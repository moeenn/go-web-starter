include .env

ENTRYPOINT = ./cmd/app/main.go
CSS_INPUT = ./resources/css/styles.css
CSS_OUTPUT = ./public/css/styles.css

setup:
	bash ./scripts/tools.sh

css_dev:
	tailwindcss -i ${CSS_INPUT} -o ${CSS_OUTPUT} --watch

css_build:
	tailwindcss -i ${CSS_INPUT} -o ${CSS_OUTPUT} --minify

server_dev:
	templ generate --watch --proxy="http://localhost:${SERVER_PORT}" --open-browser=false --cmd="go run ${ENTRYPOINT}"

views_build:
	templ generate

build:
	make css_build && \
	make views_build && \
	go build -o app ${ENTRYPOINT} && \
	rm -rf ./dist && \
	mkdir -p ./dist && \
	mv -v ./app ./dist && \
	cp -r ./public ./dist
