# track

Track is a little tool to help you out deal with bugzilla (in the future 
trello too) in a quick way from the command line.

It features smart caching, and preloading of HTML to it's internal cache,
so, once listed if you asked it to cache HTML (-x option) you will
be able to quickly do *track bz-show ID* and you will have the bugzilla
in your browser.

## installing

```bash
export GOPATH=${GOPATH:-$HOME/go}
export PATH=$PATH:$GOPATH/bin

go get github.com/mangelajo/track
```

## usage examples

If you don't have proper config, track will explain you how to create a config file
```bash
$ track bz-list
Could not read config file: Config File ".track" Not Found in "[/Users/ajo]"
No email address provided either in parameters or ~/.track.yaml file

An example ~/.track.yaml:

bzurl: https://bugzilla.redhat.com
bzemail: xxxxx@xxxx
bzpass: xxxxxxxx
dfg: Networking
htmlOpenCommand: xdg-open

```

If you want to list bugs on you, reardless of DFG:

```bash
$ track bz-list --me -d "" -x

1399987 (ASSIGNED)	majopela@redhat.com	https://bugzilla.redhat.com/1399987	   openstack-neutron	[RFE] allow to limit conntrack entries per tenant to avoid "nf_conntrack: table full, dropping packet"
1546996 (     NEW)	majopela@redhat.com	https://bugzilla.redhat.com/1546996	python-networking-ovn	[RFE] [Neutron] [OVN] QoS support
1546994 (ASSIGNED)	majopela@redhat.com	https://bugzilla.redhat.com/1546994	python-networking-ovn	[RFE] [Neutron] [OVN] Productize a migration tool from ML2/OVS to OVN

BZ 1546994 (ASSIGNED) [RFE] [Neutron] [OVN] Productize a migration tool from ML2/OVS to OVN
  Keywords: FutureFeature, Triaged
  Assigned to: majopela@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1546994
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/02058676

BZ 1546996 (     NEW) [RFE] [Neutron] [OVN] QoS support
  Keywords: FutureFeature
  Assigned to: majopela@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1546996
  * OpenStack gerrit : https://review.openstack.org/#/c/265798/
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/02058676

BZ 1399987 (ASSIGNED) [RFE] allow to limit conntrack entries per tenant to avoid "nf_conntrack: table full, dropping packet"
  Keywords: FutureFeature, ZStream
  Assigned to: majopela@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1399987
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/02037820
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/01973106
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/01955752
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/01747905

Pre caching HTML
 - bz#1546994 cached
 - bz#1546996 cached
 - bz#1399987 cached
```

This will let you open a bugzilla
```bash
$ track bz-show 1546994
Wrote /tmp/bz1546994.html
```

You can also open predefined queries
```bash
$ track bz-rh-query network-dfg-untriaged -x
...
...
...

```

