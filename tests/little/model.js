let name = "nex-"+Math.random().toString().substr(-6)
topo = {
  name: name,
  nodes: [
    deb('server'), deb('db'), 
    deb('c0'), deb('c1'), deb('c2'), deb('c3')
  ],
  switches: [cumulus('cx')],
  links: [
    Link('server', 1, 'cx', 1),
    Link('db', 1, 'cx', 2),
    Link('c0', 1, 'cx', 3, {
      mac: {
        c0: '00:00:11:11:00:01'
      }
    }),
    Link('c1', 1, 'cx', 4, {
      mac: {
        c1: '00:00:22:22:00:01'
      }
    }),
    Link('c2', 1, 'cx', 5, {
      mac: {
        c2: '00:00:33:33:00:01'
      }
    }),
    Link('c3', 1, 'cx', 6, {
      mac: {
        c2: '00:11:aa:00:00:01'
      }
    }),
  ]
}

function deb(name) {
  return {
    name: name,
    image: 'debian-bullseye',
    cpu: { cores: 4 },
    memory: { capacity: GB(4) },
    mounts: [
      { source: env.PWD+'/../../../nex', point: '/tmp/nex' }
    ]
  }
}

function cumulus(name) {
  return {
    name: name,
    image: 'cumulusvx-3.5-mvrf',
    cpu: { cores: 4 },
    memory: { capacity: GB(4) },
    mounts: [
      { source: env.PWD+'/../../../nex', point: '/tmp/nex' }
    ]
  }
}

