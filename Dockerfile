FROM golang:alpine

ADD ./*.go /go/src/github.com/vangroan/shorter/

# Add dependencies
ADD ./vendor/**/* /go/src/

# 0.    Set some shell flags like `-e` to abort the 
#       execution in case of any failure (useful if we 
#       have many ';' commands) and also `-x` to print to 
#       stderr each command already expanded.
# 1.    Get into the directory with the golang source code
# 2.    Perform the go build with some flags to make our
#       build produce a static binary (CGO_ENABLED=0 and 
#       the `netgo` tag).
# 3.    copy the final binary to a suitable location that
#       is easy to reference in the next stage
RUN set -ex && \
    cd /go/src/github.com/vangroan/shorter && \       
    CGO_ENABLED=0 go build \
    -tags netgo \
    -v -a \
    -ldflags '-extldflags "-static"' && \
    mv ./shorter /usr/bin/shorter

# Set the binary as the entrypoint of the container
ENTRYPOINT [ "shorter" ]