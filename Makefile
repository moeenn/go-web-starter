setup:
	sh ./scripts/tools.sh

dev:
	templ generate --watch --proxy="http://localhost:3000" --open-browser=false --cmd="go run ."



