profile =personal

.PHONY: build clean deploy gomodgen

build:
	env GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/tracking-lambda/bootstrap ./cmd/empty-pixel-handler
	zip -j bin/tracking-lambda/tracking-lambda.zip bin/tracking-lambda/bootstrap

	env GOARCH=arm64 GOOS=linux go build -tags lambda.norpc -ldflags="-s -w" -o bin/app/bootstrap ./cmd/app
	zip -j bin/app/app.zip bin/app/bootstrap
clean:
	rm -rf ./bin 

deploy: build
	sls deploy --verbose --aws-profile $(profile)

