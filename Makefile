dev-build:
	docker build -t go-esun-jpy-crawler -f Dockerfile-dev .

dev-console:
	docker run --env-file .env -v $$(pwd):/opt/app/ -it --rm go-esun-jpy-crawler sh
