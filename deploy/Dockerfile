FROM golang:1.22

RUN mkdir -p /data/bin/etc/migrations

WORKDIR /data/bin

COPY ../app-go /data/bin/app-go
RUN chmod +x /data/bin/app-go

COPY ./etc/* /data/bin/etc/
COPY ./etc/migrations/* /data/bin/etc/migrations/

CMD ["/data/bin/app-go", "-f", "/data/bin/etc/config.yaml"]
