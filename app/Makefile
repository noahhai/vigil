.PHONY: build clean deploy

build-api:
	env CG_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/main main.go

deploy-api: build-api
	sls deploy --verbose

build-web:
	(cd webapp && npm run build)

deploy-web: build-web
	aws s3 sync ./webapp/dist s3://getvigil.io
	aws cloudfront create-invalidation --distribution-id E2VR2UI5YH62G7 --paths /index.html /main.js "/static/*"
