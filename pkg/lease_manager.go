package nex

import (
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	quantum = 1 * time.Minute
)

func RunLeaseManager() {

	for ; ; time.Sleep(quantum) {

		//TODO(ry,bug):
		//      there appear to be cases in which the following can return
		//      duplicate members which causes the lease recycle transaction to
		//      fail. I believe the following may be the case.
		//      - a leak causes two addresses to reference the same mac so the
		//        FetchIp4IndexMembers list has duplicate mac entries
		//      - when the list of macs with dupes is passed to FetchMembers the
		//        fetch members produces duplicate entries for the duplicate macs
		//      - now dupe members are being passed to RecycleExpiredLeases which
		//        causes transaction errors
		//members, err := FetchIp4IndexMembers()
		members, err := FetchMacIndexMembers()

		if err != nil {
			log.WithError(err).Errorf("lease-manager: fetch ip4 index failed")
			continue
		}

		err = RecycleExpiredLeases(members)
		if err != nil {
			log.WithError(err).Errorf("lease-manager: recycle failed")
			continue
		}
	}

}

func RecycleExpiredLeases(members []*Member) error {

	var updates, trash []Object
	nets := make(map[string]*NetworkObj)
	poolUpdates := make(map[string]*PoolObj)

	for _, m := range members {

		if m.Ip4 == nil {
			continue
		}
		if m.Ip4.Expires == nil {
			continue
		}

		expires := time.Unix(m.Ip4.Expires.Seconds, int64(m.Ip4.Expires.Nanos))
		if time.Now().After(expires) {

			// remove the ipv4 index
			trash = append(trash, NewIp4Index(m))

			// update the mac index by removing the address from the member object
			u := m.Clone()
			u.Ip4 = nil
			updates = append(updates, NewMacIndex(u))

			// update the pool the expired lease came from
			var pool *PoolObj

			// first try to get the network + associated pool from local cache
			network, ok := nets[m.Net]
			if !ok {

				network = NewNetworkObj(&Network{Name: m.Net})
				pool = NewPoolObj(&Pool{Net: m.Net})
				_, err := ReadObjects([]Object{network, pool})
				if err != nil {
					return err
				}
				nets[m.Net] = network
				poolUpdates[m.Net] = pool
				updates = append(updates, pool)

			} else {
				pool = poolUpdates[m.Net]
			}

			index := network.Range4.Offset(net.ParseIP(m.Ip4.Address))
			pool.CountSet = pool.CountSet.Remove(index)

			log.WithFields(log.Fields{
				"member": m.Mac,
				"addr":   m.Ip4.Address,
				"index":  index,
			}).Info("recycling expired address")

		}

	}

	return RunObjectTx(ObjectTx{Put: updates, Delete: trash})

}
