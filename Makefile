#! /usr/bin/make
#
# Makefile for goa chat
#
# Targets:
# - clean     delete all generated files
# - generate  (re)generate all goagen-generated files.
# - build     compile executable
#
# Meta targets:
# - all is the default target, it runs all the targets in the order above.
#
DBNAME = goa_chat
.PHONY: all start depend bootstrap generate models client build run

all: depend clean generate build generate-client

depend:
	@glide install

clean:
	@rm -rf app
	@rm -rf client
	@rm -rf tool
	@rm -rf public/swagger
	@rm -rf public/schema
	@rm -rf public/js
	@rm -f todo

bootstrap:
	@goagen main    -d github.com/m0a-mystudy/goa-chat/design -o controllers

generate:
	@goagen app     -d github.com/m0a-mystudy/goa-chat/design
	@goagen swagger -d github.com/m0a-mystudy/goa-chat/design -o public
	@goagen schema  -d github.com/m0a-mystudy/goa-chat/design -o public
	@goagen client  -d github.com/m0a-mystudy/goa-chat/design
	@goagen js      -d github.com/m0a-mystudy/goa-chat/design -o public

generate-client:
	@swagger-codegen generate -l typescript-fetch -t ./codegen-template -i ./public/swagger/swagger.json -o ./chat-client-api
	@jq -s '.[0] * .[1]' chat-client-api/package.json chat-client-api/package_replace.json > replaced_package.json
	@rm chat-client-api/package.json
	@mv replaced_package.json chat-client-api/package.json

models:
	@xo mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@localhost/$(DBNAME)  -o models -t ./xo_templates

build:
	@go build -o chat

build-client:
	@cd chat-client-api && yarn
	@cd chat-client-api && npm link

dev:
	@on | xargs -n1 -I{} make build
run:
	@chat

# ae-build:
# 	@if [ ! -d $(HOME)/cellar ]; then \
# 		mkdir $(HOME)/cellar; \
# 		ln -s $(CURRENT_DIR)/appengine.go $(HOME)/cellar/appengine.go; \
# 		ln -s $(CURRENT_DIR)/app.yaml     $(HOME)/cellar/app.yaml; \
# 	fi

# ae-deploy: ae-build
# 	cd $(HOME)/cellar
# 	python2 appcfg.py update .
