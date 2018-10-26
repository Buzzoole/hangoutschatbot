get-deps:
	dep ensure

build:
	go build -o hangoutschatbot .

docker-build:
	docker build -t appsterdam/hangoutschatbot .

docker-push:
	docker login -u appsterdam -p $(DOCKERPASSWORD)
	docker push appsterdam/hangoutschatbot

docker-run:
	docker run --rm -p 8080:8080 -v $$(pwd)/credentials.json:/credentials.json appsterdam/hangoutschatbot:latest