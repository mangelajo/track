VERSION=`git describe --tags 2> /dev/null || echo "dev"`
BUILD=`date +%FT%T%z`

LDFLAGS="-X cmd.Version=${VERSION} -X cmd.Build=${BUILD}"

track:
	go build -ldflags ${LDFLAGS} -o track main.go

test:
	go test github.com/mangelajo/track/..

clean:
	rm track

all: track

