FROM alpine:3.14.2
WORKDIR /srv/
VOLUME /srv/runtime
COPY sufficient /srv/sufficient
COPY migrations /srv/migrations
COPY conf.toml /srv/conf.toml
RUN chmod +x /srv/sufficient
EXPOSE 9080
