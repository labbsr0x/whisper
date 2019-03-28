# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial 

ENV p $GOPATH/src/github.com/abilioesteves/whisper

ADD ./ ${p}
WORKDIR ${p}
RUN go get -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /whisper main.go

## PKG
FROM alpine

RUN mkdir -p /www

COPY --from=builder /whisper /
COPY web/ui/www/ /www/

CMD [ "/whisper", "serve" ]
