FROM alpine:3.6

RUN adduser -D oauth2-proxy
USER oauth2-proxy

ADD tmp/_output/bin/oauth2-proxy /usr/local/bin/oauth2-proxy
