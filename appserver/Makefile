docker:
	docker build --no-cache --rm \
		--build-arg BUILD_TYPE=dev \
		--build-arg APP_VERSION=$(shell cat version.info) \
		--build-arg APP_NAME=tsm-appserver \
		-t tsm-appserver:$(shell cat version.info).dev .
	docker tag tsm-appserver:$(shell cat version.info).dev tsm-appserver:latest

docker-run:
	docker run -d --rm --name tsm-appserver -p 3000:3000 tsm-appserver:$(shell cat version.info).dev
	docker logs -f tsm-appserver

docker-stop:
	docker stop tsm-appserver

docker-push:
	aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 915486611144.dkr.ecr.ap-northeast-2.amazonaws.com
	docker tag tsm-appserver:latest 915486611144.dkr.ecr.ap-northeast-2.amazonaws.com/tsm-appserver:$(shell cat version.info).dev
	docker push 915486611144.dkr.ecr.ap-northeast-2.amazonaws.com/tsm-appserver:$(shell cat version.info).dev
