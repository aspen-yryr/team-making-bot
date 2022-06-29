.PHONY: gen_secret
gen_secret:
	@echo -n "Please input your bot token to make .env.dev file.\n>"
	@bash -c 'read -s -p "" KEY && \
	cat ./env/.env.dev.example | sed -e "s/__TOKEN__HERE__/$$KEY/" ./env/.env.dev.example' > ./env/.env.dev &&  echo


.PHONY: build
build:
	@CGO_ENABLED=0 go build --ldflags="-s -w" -o app cmd/app/main.go


.PHONY: run_build
run_build:
	@./app -v 6 -logtostderr --env_file=./env/.env.dev --greet=false
