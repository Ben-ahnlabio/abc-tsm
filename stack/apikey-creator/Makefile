init: clean
	poetry export --without-hashes --format=requirements.txt > ./key_creator/requirements.txt
	sam build --use-container
	cd .aws-sam/build/ApiKeyCreatorFunction && zip -r ApiKeyCreatorFunction.zip ./
	mkdir -p ./dist
	mv .aws-sam/build/ApiKeyCreatorFunction/ApiKeyCreatorFunction.zip ./dist

clean:
	rm -rf .aws-sam/build
	rm ./dist/*

deploy: init
	aws s3 cp ./dist/ApiKeyCreatorFunction.zip s3://abc-tsm/function/