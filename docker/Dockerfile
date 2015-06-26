# Skycore factory - a compilation container

FROM ubuntu:14.04

RUN apt-get update && apt-get install -y \
	git-core \
	bzr \
	make \
	gcc \
	mercurial \
	ca-certificates \
	curl


RUN curl -s https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar -v -C /usr/local -xz

ENV GOROOT /usr/local/go
ENV PATH /usr/local/go/bin:/gopath/bin:/usr/local/bin:$PATH
ENV GOPATH /gopath/

# compile.sh script for 
RUN echo '#!/bin/bash' > /compile.sh ; \
	echo 'export GOPATH=/gopath/ ; export CGO_ENABLED=0 ; go get -a -v github.com/wgerlach/Skycore/skycore' >> /compile.sh ; \
	chmod +x /compile.sh


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


