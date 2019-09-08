FROM  alpine:3.10
MAINTAINER gw123  <963353840@qq.com>
COPY release/alpine/server.upload  /entry/server
COPY release/config.server.yaml  /etc/gserver/config.server.yaml
#RUN apk add redis && rm -f /var/cache/apk/*
WORKDIR /entry
EXPOSE 8080
ENTRYPOINT ["/entry/server"]
