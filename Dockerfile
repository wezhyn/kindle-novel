FROM golang:1.13 as builder

WORKDIR /go/src/kindle
COPY . .
RUN  go get -d -v ./... \
    && CGO_ENABLED=0 GOOS=linux go build -a -v -o kindle . \
    && cp kindle /root \
    && cp config.yml /root
	&& cp ./kindle.sh /root 

FROM alpine:3.11
COPY --from=builder /root .
WORKDIR /root

RUN  chmod +x /root/kindle && chmod +x /root/kindle.sh


CMD ["sh","/root/kindle.sh"]



