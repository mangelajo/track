# track

Track is a little tool to help you out deal with bugzilla and trello
in a quick way from the command line.

It features smart caching, and preloading of HTML to it's internal cache,
so, once listed if you asked it to cache HTML (-x option) you will
be able to quickly do *track bz show ID* and you will have the bugzilla
in your browser.

## installing

```bash
export GOPATH=${GOPATH:-$HOME/go}
export PATH=$PATH:$GOPATH/bin

go get github.com/mangelajo/track
```

## upgrading

```bash
export GOPATH=${GOPATH:-$HOME/go}
export PATH=$PATH:$GOPATH/bin

go get -u github.com/mangelajo/track
```

## basic help

```bash
$$ ./track bz
 Bugzilla related commands

 Usage:
   track bz [command]

 Available Commands:
   list        List bugzillas based on parameters and configuration
   query       Grab query parameters from your config
   rh-query    Grab query parameters from https://url.corp.redhat.com/< name >
   show        Open cached HTML for bugzilla

 Flags:
   -u, --bzemail string   Bugzilla login email
   -k, --bzpass string    Bugzilla login password
   -b, --bzurl string     Bugzilla URL (default "https://bugzilla.redhat.com")
   -h, --help             help for bz
   -x, --html             Pre-cache html for bz show command
       --shell            Start an interactive shell once the command is done

 Global Flags:
       --config string            config file (default is $HOME/.track.yaml)
       --htmlOpenCommand string   Command to open an html file (default "xdg-open")
   -i, --ignorecerts              Ignore SSL certificates
   -l, --limit int                Max entries to list (default 50)
   -o, --offset int               Offset on the bug listing
   -w, --workers int              Workers for http retrieval (default 4)
```

## usage examples

If you don't have proper config, track will explain you how to create a config file
```bash
$ track bz list
Could not read config file: Config File ".track" Not Found in "[/Users/ajo]"
No email address provided either in parameters or ~/.track.yaml file

An example ~/.track.yaml:

bzurl: https://bugzilla.redhat.com
bzemail: xxxxx@xxxx
bzpass: xxxxxxxx # you can omit this field and track will ask you when needed
dfg: Networking
htmlOpenCommand: xdg-open
queries:
    ovn-new: https://bugzilla.redhat.com/buglist.cgi?bug_status=NEW&classification=Red%20Hat&component=python-networking-ovn&list_id=8959835&product=Red%20Hat%20OpenStack&query_format=advanced
    ovn-rfes: https://bugzilla.redhat.com/buglist.cgi?bug_status=NEW&bug_status=ASSIGNED&bug_status=MODIFIED&bug_status=ON_DEV&bug_status=POST&bug_status=ON_QA&classification=Red%20Hat&component=python-networking-ovn&f1=keywords&f2=short_desc&j_top=OR&list_id=8959855&o1=substring&o2=substring&product=Red%20Hat%20OpenStack&query_format=advanced&v1=RFE&v2=RFE

# notes: for OSX use htmlOpenCommand: open

```

If you want to list bugs on you, regardless of DFG:

```bash
$ track bz list --me -d "" -x

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
$ track bz show 1546994
Wrote /tmp/bz1546994.html
```

You can also open predefined queries
```bash
$ track bz rh-query network-dfg-untriaged -x
...
...
...

```

## The interactive shell

Just append --shell to bz list or bz rh-query , and there it is.
```bash
$ track bz list -x --shell
Track interactive shell
BZ 1578502 (     NEW) [RFE] Networker Node replacement documentation
  Product: Red Hat OpenStack ver: 10.0 (Newton) target: 10.0 (Newton) (---)
  Keywords:
  Assigned to: rhos-docs@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1578502
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/02101007

>>> help

Commands:
  clear      clear the screen
  exit       exit the program
  go         open bugzilla from server url
  help       display help
  links      open links from bugzilla
  next       next bugzilla
  open       open a bugzilla from cache
  prev       previous bugzilla
  show       show a bugzilla
```
