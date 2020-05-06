GOARCH=amd64 GOOS=linux go build -i -o ./build/github_tgbot_linux main.go
chmod 777 ./build/github_tgbot_linux
