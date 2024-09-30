// package main

// import (
// 	"fmt"
// 	"log"
// 	"os/exec"
// 	"runtime"

// 	"github.com/vishvananda/netns"
// )

// func main() {
// 	runtime.LockOSThread()
// 	defer runtime.UnlockOSThread()

// 	// Save the current network namespace
// 	// origns, _ := netns.Get()
// 	// defer origns.Close()
// 	newns, err := netns.GetFromName("nginx")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	netns.Set(newns)
// 	res := exec.Command("ping","8.8.8.8")

//		if err := res.Run(); err != nil {
//			fmt.Println("here 1")
//			log.Fatal(err)
//		}
//		stdout, err := res.Output()
//		if err != nil {
//			fmt.Println("here")
//			log.Fatal(err)
//		}
//		fmt.Println(string(stdout))
//		// var result  []byte
//		// _, err = res.Stderr.Write(result)
//		// if err!= nil{
//		// 	log.Fatal(err)
//		// }
//		newns.Close()
//	}
package main

import (
	"fmt"
	"strings"
)

type name struct {
	name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func main() {
	// out, errout, err := Shellout("ping 8.8.8.8")
	// if err != nil {
	// 	log.Printf("error: %v\n", err)
	// }
	// fmt.Println("--- stdout ---")
	// fmt.Println(out)
	// fmt.Println("--- stderr ---")
	// fmt.Println(errout)
	//   var data = `{"name": "amine"}`
	// net := &nex.Network{}
	// // if err := json.Unmarshal( []byte(data),net); err != nil {
	// //     log.Fatal(err)
	// // }
	// if net == nil {
	// 	fmt.Println("hello")
	// }
	aa := []string{
		"amine@test", "amine",
	}
	for _, a := range aa {
		v := strings.Split(a, "@")
		if len(v) > 1 {
			if strings.Split(a, "@")[0] == "amine" {
				fmt.Println(a)
			}
		}
	}
}
