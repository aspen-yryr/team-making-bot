.PHONY: gen_secret
gen_secret:
	@echo -n "Please input your bot token to make .env.dev file.\n>"
	@bash -c 'read -s -p "" KEY && \
	cat ./env/.env.dev.example | sed -e "s/__TOKEN__HERE__/$$KEY/" ./env/.env.dev.example' > ./env/.env.dev &&  echo


