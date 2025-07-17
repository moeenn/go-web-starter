include .env

ENV_FILE = .env
ENTRYPOINT = ./cmd/app/main.go
CSS_INPUT = ./resources/css/styles.css
CSS_OUTPUT = ./public/css/styles.css


.PHONY: setup
setup:
	bash ./scripts/tools.sh


css_dev:
	tailwindcss -i ${CSS_INPUT} -o ${CSS_OUTPUT} --watch


gen_css:
	tailwindcss -i ${CSS_INPUT} -o ${CSS_OUTPUT} --minify


.PHONY: server_dev
server_dev:
	templ generate --watch \
	--proxy="http://localhost:${SERVER_PORT}" \
	--open-browser=false \
	--cmd="godotenv -f ${ENV_FILE} go run ${ENTRYPOINT}"


gen_views:
	templ generate


.PHONY: lint
lint:
	golangci-lint run ./...


build: gen_css gen_views
	go build -o app ${ENTRYPOINT}


pkg: build
	rm -rf ./dist && \
	mkdir -p ./dist && \
	mv -v ./app ./dist && \
	cp -r ./public ./dist


.PHONY: clean
clean:
	rm -rvf ./app && \
	rm -rvf ./dist && \
	rm -vf ./views/**/*_templ.go


.PHONY: gen_token
gen_token:
	go run github.com/moeenn/go-token@latest


migration_new:
	goose -s create ${NAME} sql


.PHONY: db_migrate
db_migrate:
	goose up


.PHONY: db_rollback
db_rollback:
	goose down-to 0
