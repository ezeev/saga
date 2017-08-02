FROM alpine:3.6
MAINTAINER evan@viroonga.com

RUN apk --update add ca-certificates

ADD viroongacom /usr/bin/viroongacom
ADD static static
ADD view view

EXPOSE 3001
ENTRYPOINT ["/usr/bin/viroongacom"]