
FROM nginx 
# baseimage is offical nginx release

RUN apt-get update && apt-get -y install git curl net-tools
RUN cd / && git clone https://github.com/wgerlach/Skycore.git
RUN rm -f /etc/nginx/nginx.conf /etc/nginx/sites-enabled/* /etc/nginx/sites-available/default ; \
	ln -s /Skycore/services/nginx/nginx.conf /etc/nginx/nginx.conf ;
	mkdir -p /etc/nginx/sites-enabled/
	#ln -s /Skycore/services/nginx/default_server /etc/nginx/default_server


#confd
RUN curl -L -o /usr/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.8.0/confd-0.8.0-linux-amd64 && chmod +x /usr/bin/confd

RUN mkdir -p /etc/confd/conf.d/ /etc/confd/templates/ ; \
	ln -s /Skycore/services/nginx/confd/nginx.toml /etc/confd/conf.d/nginx.toml ; \
	ln -s /Skycore/services/nginx/confd/nginx.tmpl /etc/confd/templates/nginx.tmpl


CMD ["/root/Skycore/services/nginx/confd/run_confd.sh"]

# Expose ports.
EXPOSE 80
EXPOSE 443
