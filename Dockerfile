# BUILD
FROM abilioesteves/gowebbuilder:unstable as builder

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
