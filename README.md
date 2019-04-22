# Whisper
[![Docker Pulls](https://img.shields.io/docker/pulls/labbsr0x/whisper.svg)](https://hub.docker.com/r/labbsr0x/whisper)

A Login and Consent provider implementation in Go for Ory Hydra



# Developer

From the project root folder, fire the following commands to execute this project in development mode:

```
go build
```

and

```
./whisper serve --port 7070 --base-ui-path ./web/ui/www --hydra-admin-url http://localhost:4445 --hydra-public-url http://localhost:4444 --client-id whisper --client-secret whisper --log-level debug --scopes-file-path ./scopes.json
```

This will serve Whisper at the 7070 port; endpoints `/login` and `/consent` will display our incredibly shy `index.html` example.
