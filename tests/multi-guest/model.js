
function cumulus(name) {
  return {
    name: name,
    image: 'cumulusvx-3.5-mvrf',
    cpu: { cores: 2 },
    memory: { capacity: MB(512) }
  }
}

function fedora(name, mem) {
  return {
    name: name,
    image: 'fedora-28',
    cpu: { cores: 2 },
    memory: { capacity: MB(mem) },
    mounts: [
      { source: env.PWD+'/../../../nex', point: '/tmp/nex' }
    ]
  }
}

function debian(name, mem) {
  return {
    name: name,
    image: 'debian-buster',
    cpu: { cores: 2 },
    memory: { capacity: MB(mem) },
    mounts: [
      { source: env.PWD+'/../../../nex', point: '/tmp/nex' }
    ]
  }
}

function onie(name, mem) {
  return {
    name: name,
    image: 'onie-x86',
    firmware: 'OVMF-pure-efi.fd',
    os: 'onie',
    cpu: { cores: 2 },
    memory: { capacity: MB(mem) }
  }
}

topo = {
  name: 'nex0',
  nodes: [
    // infrastructure 
    debian('s0', 1024),
    debian('s1', 1024),
    debian('db', 1024),

    // clients

    // internal
    debian('i0', 512),
    debian('i1', 512),
    // pxeboot
    onie('p0', 512),
    onie('p1', 512),
    onie('p2', 512),
    // embedded
    debian('e0', 512),
    debian('e1', 512),
    debian('e2', 512),
    // vms
    debian('v0', 512),
    debian('v1', 512),

  ],
  switches: [
    cumulus('sw')
  ],
  links: [
    Link('s0', 1, 'sw', 1),
    Link('s1', 1, 'sw', 2),
    Link('db', 1, 'sw', 3),

    Link('i0', 1, 'sw', 4),
    Link('i1', 1, 'sw', 5),


    Link('p0', 1, 'sw', 6, {
      boot: 1, 
      mac: {
        p0: '00:00:99:10:00:01'
      }
    }),
    Link('p1', 1, 'sw', 7, {
      boot: 1, 
      mac: {
        p1: '00:00:99:22:00:11'
      }
    }),
    Link('p2', 1, 'sw', 8, {
      boot: 1, 
      mac: {
        p2: '00:00:99:AB:00:CA'
      }
    }),

    Link('e0', 1, 'sw', 9),
    Link('e1', 1, 'sw', 10),
    Link('e2', 1, 'sw', 11),

    Link('v0', 1, 'sw', 12),
    Link('v1', 1, 'sw', 13)
  ]
}
