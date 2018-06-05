VERSION=`git describe --tags 2> /dev/null || echo "dev"`
BUILD=`date +%FT%T%z`

LDFLAGS="-X cmd.Version=${VERSION} -X cmd.Build=${BUILD}"

pkg/bugzilla/xmlbugzilla.go:
	tools/gen_xml_code.sh > $@

track: pkg/bugzilla/*.go cmd/*.go main.go
	go build -ldflags ${LDFLAGS} -o track main.go

test:
	go test github.com/mangelajo/track/..

clean:
	rm track

clean-xml:
	rm pkg/bugzilla/xmlbugzilla.go

all: track

