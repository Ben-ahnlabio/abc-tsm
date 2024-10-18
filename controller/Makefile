docker:
	docker build --no-cache --rm --network=host \
		--build-arg BUILD_TYPE=dev \
		--build-arg APP_VERSION=$(shell cat version.info) \
		--build-arg APP_NAME=tsm-controller \
		-t tsm-controller:$(shell cat version.info).dev .
	docker tag tsm-controller:$(shell cat version.info).dev tsm-controller:latest

docker-run: docker-stop
	docker run -d --rm --name tsm-controller1 --env-file=.env.node1 -p 4001:3000 tsm-controller:$(shell cat version.info).dev
	docker run -d --rm --name tsm-controller2 --env-file=.env.node2 -p 4002:3000 tsm-controller:$(shell cat version.info).dev

docker-stop:
	docker stop tsm-controller1
	docker stop tsm-controller2
	docker rm tsm-controller1
	docker rm tsm-controller2

docker-push: docker
	aws ecr get-login-password --region ap-northeast-2 | docker login --username AWS --password-stdin 915486611144.dkr.ecr.ap-northeast-2.amazonaws.com
	docker tag tsm-controller:latest 915486611144.dkr.ecr.ap-northeast-2.amazonaws.com/tsm-controller:$(shell cat version.info).dev
	docker push 915486611144.dkr.ecr.ap-northeast-2.amazonaws.com/tsm-controller:$(shell cat version.info).dev
