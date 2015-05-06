# image to run nginx and confd containers (one container for each)

FROM nginx 
# baseimage is offical nginx release

RUN apt-get update && apt-get -y install git curl net-tools

# config files and run_confd.sh
RUN cd / && \
  git clone https://github.com/wgerlach/Skycore.git && \
  mkdir -p /etc/nginx/sites-enabled/

# required to update nginx in another container
ADD https://get.docker.com/builds/Linux/x86_64/docker-1.6.0 /usr/bin/docker

# install confd
RUN curl -L -o /usr/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.9.0/confd-0.9.0-linux-amd64 && \
  chmod +x /usr/bin/confd && chmod +x /usr/bin/docker


# Expose nginx ports.
EXPOSE 80
EXPOSE 443
