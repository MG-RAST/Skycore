
FROM nginx 
# baseimage is offical nginx release

RUN apt-get update && apt-get -y install git curl net-tools

# config files and run_confd.sh
RUN cd / && \
  git clone https://github.com/wgerlach/Skycore.git && \
  mkdir -p /Skycore/services/nginx/sites-enabled


# install confd
RUN curl -L -o /usr/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.8.0/confd-0.8.0-linux-amd64 && chmod +x /usr/bin/confd

# start confd and nginx
CMD ["/root/Skycore/services/nginx/confd/run_confd.sh"]

# Expose ports.
EXPOSE 80
EXPOSE 443
