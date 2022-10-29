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