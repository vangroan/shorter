FROM golang:alpine as build

# installs `gcc`, required to compile sqlite3
RUN apk add build-base

COPY ./*.go /go/src/github.com/vangroan/shorter/

# Add dependencies
COPY ./vendor/ /go/src/github.com/vangroan/shorter/vendor/

# 0.    Set some shell flags like `-e` to abort the 
#       execution in case of any failure (useful if we 
#       have many ';' commands) and also `-x` to print to 
#       stderr each command already expanded.
# 1.    Get into the directory with the golang source code
# 2.    Perform the go build with some flags to make our
#       build produce a static binary (CGO_ENABLED=1 and 
#       the `netgo` tag). `cgo` is required to compile
#       `sqlite3`.
# 3.    copy the final binary to a suitable location that
#       is easy to reference in the next stage
RUN set -ex && \
    cd /go/src/github.com/vangroan/shorter && \       
    CGO_ENABLED=1 go build \
    -tags netgo \
    -v -a \
    -ldflags '-extldflags "-static"' && \
    mv ./shorter /usr/bin/shorter

# New image since build is bloated by `gcc`
FROM golang:alpine as runtime

# Override target file for SQLite3 database
ENV SHORTER_SQLITE /var/db/shorter.sqlite3

COPY --from=build /usr/bin/shorter /usr/bin/shorter

# Static files
RUN mkdir /assets/
COPY ./assets/* /assets/
WORKDIR /assets/

EXPOSE 8000

VOLUME ["/var/db"]

# Set the binary as the entrypoint of the container
ENTRYPOINT [ "shorter" ]