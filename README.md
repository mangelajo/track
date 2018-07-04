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
$ ./track bz
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

```bash
$ ./track bz list --help
This command will list and retrieve details for bugzillas
based on configuration and query.

Usage:
  track bz list [flags]

Flags:
  -a, --assignee string           Filter by assignee (you can use 'me'
      --changed                   Show bugs changed since last run
      --class string              Class on bugzilla (default "Red Hat")
  -c, --component string          Component
  -d, --dfg string                Openstack DFG
  -f, --flags-on string           List bugs with flags on somebody (you can use 'me')
  -h, --help                      help for list
  -m, --me                        List only bugs assigned to me
  -p, --product string            Product
      --squad string              Openstack DFG Squad
  -s, --status string             Status list separated by commas (default "NEW,ASSIGNED,POST,MODIFIED,ON_DEV,ON_QA")
  -t, --target-milestone string   Target milestone
  -r, --target-release string     Target release

Global Flags:
      --bzemail string           Bugzilla login email
  -k, --bzpass string            Bugzilla login password
  -b, --bzurl string             Bugzilla URL (default "https://bugzilla.redhat.com")
      --config string            config file (default is $HOME/.track.yaml)
  -x, --html                     Pre-cache html for bz show command
      --htmlOpenCommand string   Command to open an html file (default "xdg-open")
  -i, --ignorecerts              Ignore SSL certificates
  -l, --limit int                Max entries to list (default 50)
  -o, --offset int               Offset on the bug listing
      --shell                    Start an interactive shell once the command is done
  -u, --summary                  Show a summary of the bugs we retrieve
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
htmlOpenCommand: xdg-open  # notes: for OSX use open
queries:
    ovn-new: https://bugzilla.redhat.com/buglist.cgi?bug_status=NEW&classification=Red%20Hat&component=python-networking-ovn&list_id=8959835&product=Red%20Hat%20OpenStack&query_format=advanced
    ovn-rfes: https://bugzilla.redhat.com/buglist.cgi?bug_status=NEW&bug_status=ASSIGNED&bug_status=MODIFIED&bug_status=ON_DEV&bug_status=POST&bug_status=ON_QA&classification=Red%20Hat&component=python-networking-ovn&f1=keywords&f2=short_desc&j_top=OR&list_id=8959855&o1=substring&o2=substring&product=Red%20Hat%20OpenStack&query_format=advanced&v1=RFE&v2=RFE

```

If you want to list bugs on you, regardless of DFG (and you have a DFG in config)

```bash
$ ./track bz list --me -d "" -x
1546996 (     NEW)	majopela@redhat.com	https://bugzilla.redhat.com/1546996	python-networking-ovn	[RFE] [Neutron] [OVN] QoS support
1570843 (ASSIGNED)	majopela@redhat.com	https://bugzilla.redhat.com/1570843	python-networking-ovn	East/West traffic goes through controller node on DVR-VLAN deployment
1565563 (ASSIGNED)	majopela@redhat.com	https://bugzilla.redhat.com/1565563	python-networking-ovn	OVN L3HA - Not rescheduling gateways upon chassis addition can lead to routers not being in HA
1581332 (    POST)	majopela@redhat.com	https://bugzilla.redhat.com/1581332	python-networking-ovn	Internal DNS resolution does not work for fqdn
1546994 (    POST)	majopela@redhat.com	https://bugzilla.redhat.com/1546994	python-networking-ovn	[RFE] [Neutron] [OVN] Productize a migration tool from ML2/OVS to OVN

Grabbing bug details: done.
Pre caching HTML: bz#1546996 done.

5 bugs found.
```

You can get an extended summary of each BZ with -u (or --summary):
```bash
$ ./track bz list --me -d "" -x -u
1546996 (     NEW)	majopela@redhat.com	https://bugzilla.redhat.com/1546996	python-networking-ovn	[RFE] [Neutron] [OVN] QoS support
1570843 (ASSIGNED)	majopela@redhat.com	https://bugzilla.redhat.com/1570843	python-networking-ovn	East/West traffic goes through controller node on DVR-VLAN deployment
1565563 (ASSIGNED)	majopela@redhat.com	https://bugzilla.redhat.com/1565563	python-networking-ovn	OVN L3HA - Not rescheduling gateways upon chassis addition can lead to routers not being in HA
1581332 (    POST)	majopela@redhat.com	https://bugzilla.redhat.com/1581332	python-networking-ovn	Internal DNS resolution does not work for fqdn
1546994 (    POST)	majopela@redhat.com	https://bugzilla.redhat.com/1546994	python-networking-ovn	[RFE] [Neutron] [OVN] Productize a migration tool from ML2/OVS to OVN

Grabbing bug details: done.

BZ 1581332 (    POST) Internal DNS resolution does not work for fqdn
  Product: Red Hat OpenStack ver: 13.0 (Queens) target: 13.0 (Queens) (z1)
  Keywords: Triaged, ZStream
  Assigned to: majopela@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1581332
  * OpenStack gerrit : https://review.openstack.org/#/c/556828/
  * Red Hat Engineering Gerrit : https://code.engineering.redhat.com/gerrit/#/c/140039
  * Launchpad : https://bugs.launchpad.net/bugs/1757074

BZ 1565563 (ASSIGNED) OVN L3HA - Not rescheduling gateways upon chassis addition can lead to routers not being in HA
  Product: Red Hat OpenStack ver: 13.0 (Queens) target: 13.0 (Queens) (zstream)
  Keywords: Triaged, ZStream
  Assigned to: majopela@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1565563
  * Launchpad : https://bugs.launchpad.net/bugs/1762691

BZ 1546996 (     NEW) [RFE] [Neutron] [OVN] QoS support
  Product: Red Hat OpenStack ver: 14.0 (Rocky) target: --- (---)
  Keywords: FutureFeature, RFE
  Assigned to: majopela@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1546996
  * OpenStack gerrit : https://review.openstack.org/#/c/265798/
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/02058676

BZ 1570843 (ASSIGNED) East/West traffic goes through controller node on DVR-VLAN deployment
  Product: Red Hat OpenStack ver: 13.0 (Queens) target: 13.0 (Queens) (zstream)
  Keywords: Reopened, Triaged, ZStream
  Assigned to: majopela@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1570843

BZ 1546994 (    POST) [RFE] [Neutron] [OVN] Productize a migration tool from ML2/OVS to OVN
  Product: Red Hat OpenStack ver: 14.0 (Rocky) target: 14.0 (Rocky) (Upstream M2)
  Keywords: FutureFeature, RFE, Triaged
  Assigned to: majopela@redhat.com
  * bugzilla: http://bugzilla.redhat.com/1546994
  * OpenStack gerrit : https://review.openstack.org/#/c/510460/
  * Red Hat Customer Portal : https://access.redhat.com/support/cases/02058676

 done.
Pre caching HTML: done.

5 bugs found.

```


This will let you open a pre-cached bugzilla in your browser.
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
