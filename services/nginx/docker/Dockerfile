
FROM dockerfile/nginx
# https://github.com/dockerfile/nginx/blob/master/Dockerfile

RUN git clone https://github.com/wgerlach/Skycore.git
RUN rm -f /etc/nginx/nginx.conf /etc/nginx/sites-enabled/* ; \
	ln -s /root/Skycore/services/nginx/nginx.conf /etc/nginx/nginx.conf ; \
	ln -s /root/Skycore/services/nginx/default_server /etc/nginx/default_server

CMD ["nginx"]

# Expose ports.
EXPOSE 80
EXPOSE 443
