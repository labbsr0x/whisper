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

RUN mkdir -p /static

COPY --from=builder /whisper /go/bin/
ADD web/ui/static/index.html /static

ENTRYPOINT [ "/go/bin/whisper" ]

CMD [ "serve" ]
