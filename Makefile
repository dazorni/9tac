ifndef version
		$(error version must be defined. make version=someVersion)
endif

docker:
	rm -rf build/
	GOOS=linux GOARCH=amd64 go build -o build/9tac.x64.1.5 main.go
	npm install
	npm run build:js:prod
	npm run build:css
	mkdir build/public/fonts
	cp -R node_modules/font-awesome/fonts/* build/public/fonts
	cp app/views/* build/public
	docker build -t dazorni/9tac .
	docker tag dazorni/9tac dazorni/9tac:$(version)
	docker tag dazorni/9tac dazorni/9tac:latest
