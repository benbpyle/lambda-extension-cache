build:
	cd src/ext;GOOS=linux GOARCH=amd64 go build -o bin/extensions/lambda-cache-layer main.go
package: build
	cd src/ext/bin;zip -r extension.zip extensions/ 
deploy: build package
	cd src/ext/bin;aws lambda publish-layer-version  --layer-name 'lambda-cache-layer' --region us-west-2 --zip-file 'fileb://extension.zip' --profile=dev
