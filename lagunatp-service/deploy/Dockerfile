FROM alpine:3.6
MAINTAINER evan@viroonga.com

RUN apk --update add ca-certificates

ADD lagunatpcom /usr/bin/lagunatpcom


EXPOSE 3001
ENTRYPOINT ["/usr/bin/lagunatpcom"]