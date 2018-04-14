FROM ubuntu:trusty

RUN apt-get update
RUN apt-get install -y ca-certificates

ADD quote /srv/quote

COPY .missy.yml /.missy.yml

ENTRYPOINT /srv/quote

EXPOSE 8080
