VERSION=`git describe --tags 2> /dev/null || echo "dev"`
BUILD=`date +%FT%T%z`

LDFLAGS="-X cmd.Version=${VERSION} -X cmd.Build=${BUILD}"

all: pkg/bugzilla/xmlbugzilla.go track

pkg/bugzilla/xmlbugzilla.go: tools/gen_xml_code.sh
	tools/gen_xml_code.sh > $@

track: Makefile pkg/*/*.go pkg/storecache/*.go cmd/*.go main.go
	go build -ldflags ${LDFLAGS} -o track main.go

test:
	go test github.com/mangelajo/track/..

clean:
	rm track

clean-xml:
	rm pkg/bugzilla/xmlbugzilla.go

