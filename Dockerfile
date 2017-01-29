FROM alpine:latest

MAINTAINER Edward Muller <edward@heroku.com>

WORKDIR "/opt"

ADD .docker_build/kindle-live /opt/bin/kindle-live
ADD ./static /opt/static

CMD ["/opt/bin/kindle-live"]
