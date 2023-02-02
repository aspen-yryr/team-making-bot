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

.PHONY: gen_protobuf
gen_protobuf:
	@protoc --go_out=. --go_opt=paths=source_relative  \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative  \
	--proto_path=. \
	./proto/match/*.proto
