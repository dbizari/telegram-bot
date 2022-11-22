BUILDPATH := $(CURDIR)/build
PKGNAME := telegram-bot

.PHONY: build
build:
	@mkdir -p $(BUILDPATH)
	@CGO_ENABLED=0 go build -mod=vendor -ldflags -s -o $(BUILDPATH)/$(PKGNAME) ./cmd/telegram-bot-service/

deploy:
	heroku container:login
	heroku container:push web --app telegram-bot-tdl
	heroku container:release web --app telegram-bot-tdl

mockgen:
	mockgen -source=./internal/handlers/cmd/cmd.go -destination=./testing/mocks/handlers_mock/cmd/cmd_mock.go
	mockgen -source=./internal/handlers/cmd/getter/cmd_getter.go -destination=./testing/mocks/handlers_mock/cmdgetter/cmd_getter_mock.go
	mockgen -source=./internal/handlers/telegram/telegram.go -destination=./testing/mocks/handlers_mock/telegram/telegram_bot_mock.go