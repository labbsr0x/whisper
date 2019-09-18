docker-compose up -d local

go build

./whisper serve --port 7070 --base-ui-path ./web/ui/www --hydra-admin-url http://localhost:4445 --hydra-public-url http://localhost:4444 --secret-key uhSunsodnsuBsdjsbds --log-level debug --scopes-file-path ./scopes.json --database-url "mysql://root:secret@tcp(localhost:3306)/whisper?charset=utf8mb4&parseTime=True&loc=Local"

docker-compose exec hydra  hydra clients create  --endpoint http://localhost:4445 --id auth-code-client --secret secret --grant-types authorization_code,refresh_token --response-types code,id_token --scope openid,offline --callbacks http://127.0.0.1:5555/callback

docker-compose exec hydra hydra token user --client-id auth-code-client --client-secret secret --endpoint http://localhost:4444/ --port 5555 --scope openid,offline