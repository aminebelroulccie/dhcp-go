let name = "nex_basic-"+Math.random().toString().substr(-6)

topo = {
  name: name,
  nodes: [ deb('server'), deb('tango'), deb('foxtrot') ],
  switches: [ cumulus('cx') ],
  links: [ 
    Link('server',  1, 'cx', 1),
    Link('tango',   1, 'cx', 2, { mac: { tango:   '00:00:11:11:00:01' } }),
    Link('foxtrot', 1, 'cx', 2, { mac: { foxtrot: '00:00:22:22:00:02' } }) 
  ]
}

function deb(name) {
  return {
    name: name,
    image: 'debian-buster',
    cpu: { cores: 2 },
    memory: { capacity: GB(2) },
    mounts: [
      { source: env.PWD+'/../..', point: '/tmp/nex' }
    ]
  }
}

function cumulus(name) {
  return {
    name: name,
    image: 'cumulusvx-3.7',
    cpu: { cores: 2 },
    memory: { capacity: GB(2) },
  }
}
