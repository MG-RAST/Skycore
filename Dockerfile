# Skycore factory - a compilation container

FROM golang:1.6.3-alpine


WORKDIR /go/bin

COPY . /go/src/github.com/MG-RAST/Skycore

RUN cd /go/src/github.com/MG-RAST/Skycore && \
    CGO_ENABLED=0 go install -a -installsuffix cgo -v  ...


### build container
# docker build --tag skycore/factory:latest .

### use container to compile and get binary:
# mkdir -p ~/skycore_bin
# docker run -t -i --name sky_fac -v ~/skycore_bin:/gopath/bin skycore/factory:latest /compile.sh
# if you plan to compile multiple times with latest code:
# docker run -t -i --name sky_fac -v ~/skycore_bin:/gopath/bin skycore/factory:latest bash -c "cd /gopath/src/github.com/wgerlach/Skycore && git pull && /compile.sh"
# docker start sky_fac

### skycore execution within container
# container will need access to docker socket
# mkdir -p ~/skycore_bin
# docker run -t -i -v /var/run/docker.sock:/var/run/docker.sock --name sky_fac -v ~/skycore_bin:/gopath/bin skycore/factory:latest bash


