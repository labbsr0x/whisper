# Careful Whisper
A Login and Consent provider implementation in go.

# Use it yourself

# Developer

From the project root folder, fire the following commands to execute this project in development mode:

```
go build
```

and

```
./whisper serve --port 7070 --base-ui-path ./web/ui/static
```

This will serve Whisper at the 7070 port; endpoints `/login` and `/consent` will display our incredibly shy `index.html` example.