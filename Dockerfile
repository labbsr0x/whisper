# BUILD
FROM abilioesteves/gowebbuilder:v0.2.3 as builder

ENV p $GOPATH/src/github.com/abilioesteves/whisper

ADD ./ ${p}
WORKDIR ${p}
RUN go get -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /whisper main.go

## PKG
FROM alpine

ENV WHISPER_BASE_UI_PATH "/www/"
ENV WHISPER_SCOPES_FILE_PATH "/scopes.json"
ENV WHISPER_PORT ""
ENV WHISPER_DATABASE_URL ""
ENV WHISPER_HYDRA_ADMIN_ENDPOINT ""
ENV WHISPER_HYDRA_PUBLIC_ENDPOINT ""
ENV WHISPER_HYDRA_CLIENT_ID ""
ENV WHISPER_HYDRA_CLIENT_SECRET ""

RUN mkdir -p ${WHISPER_BASE_UI_PATH}

COPY --from=builder /whisper /
COPY web/ui/www/ /www/
COPY scopes.json /scopes.json

CMD [ "/whisper", "serve" ]
