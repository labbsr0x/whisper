# Whisper
[![Docker Pulls](https://img.shields.io/docker/pulls/labbsr0x/whisper.svg)](https://hub.docker.com/r/labbsr0x/whisper)

Whisper is an Identity and OAuth provider implemented in Go.

By definition, it securely stores user credentials and authorizes clients in the context of OAuth 2.0.

The OAuth flows adopted by Whisper are the `authorization_code` and `authorization_code + pkce` flows. 

To implement such OAuth flows, we have integrated Whisper with [Hydra](https://github.com/ory/hydra), through its **Login and Consent flow** spec.

To easily integrate with Whisper, client applications can make use of a special library and cli utility called [whisper-client](https://github.com/labbsr0x/whisper-client), available on github and also implemented in Go.

# Passwords

Passwords are cryptographically stored with a salt and secret-key with the help of the following equation:

```
HMAC(SHA512(password+salt), secret-key)
```

Salts are generated randomically each time a password is stored.

The secret-key is unique and should not be changed after the app goes up, otherwise Whisper will be unable to verify the validity of old passwords.

# Client registration

To register your application as a client, you need to be able to talk privately with Whisper. 

The admin endpoint is managed by Hydra and is by default not publicly available. Use the provided whisper-client library to register your app as a valid Whisper client.

If the application you are developing is a command-line interface or a mobile device, you need to use the `authorization_code` grant type with `pkce`.

The `pkce` (pronouced pixy) flow is specified by the [RFC7636](https://tools.ietf.org/html/rfc7636).

# Login and Consent

After Whisper is up and running, to use it as a login and consent provider one needs to generate an authorization url and redirect the browser to it.

Then, the user will be prompted to insert its login info and consent to any scopes the client application is requesting from Whisper (the requested scopes need to be previously known to Whisper).

Whisper will then redirect the browser to the client registered `redirect_uri` with a Code. This code should then be exchanged for a token to finish login.

All this operations can be more easily accomplished using the whisper-client library.

# Try it yourself

From the project root folder, fire the following commands to execute this project in development mode:

```
docker-compose up -d local
```

```
go build
```

```
./whisper serve --port 7070 \
    --base-ui-path ./web/ui/www \
    --hydra-admin-url http://localhost:4445 \
    --hydra-public-url http://localhost:4444 \
    --secret-key uhSunsodnsuBsdjsbds \
    --log-level debug \
    --scopes-file-path ./scopes.json \
    --database-url "mysql://root:secret@tcp(localhost:3306)/whisper?charset=utf8mb4&parseTime=True&loc=Local"
```

This will serve the auxiliary services (databases and Hydra) and will run Whisper at the `7070` port; endpoints `/login` and `/consent` will display our incredibly simple user interface.

The only way to access this endpoints is through a valid authorization url.

**With some small Dockerfile trickery, it is possible to override the provided UI files to use your custom page layout and icons.**

# Credential Update

A special interface is provided for the update of user credentials. 

The interface is found in the `/secure/update` endpoint and needs a valid token to access it.

The token will be searched in the request's bearer authorization header or via a `token` query param, e.g.:

```
https://<your-whisper-domain>/secure/update?token=nZgyaH1JthU0GIsp2ndRDrYNFE_6ivOqjrQhikIQ5rk.8u5lgf7OtGDbN4Y2GXcTudf1u8lLX3kvsYkFH3uPxrY&redirect_to=http://<your-app-domain>/<where-you-were>
```

After successfuly updating credentials, the UI is redirected back to where it came from, via the provided `redirect_to` query param.