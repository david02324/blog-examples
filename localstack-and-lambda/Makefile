init: stop zip-executable build start create-table deploy-lambda create-event-source-mapping

start:
	docker-compose -f docker-compose.yml up -d
	sleep 2 # grace time

stop:
	docker-compose -f docker-compose.yml down --remove-orphans

zip-executable:
	zip memo.zip memo

build:
	GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o memo main.go

create-table:
	aws dynamodb create-table --cli-input-json file://memos-table.json --endpoint-url http://localhost:4566

redeploy: remove-lambda zip-executable build deploy-lambda

deploy-lambda:
	aws lambda create-function \
		--function-name memo-lambda \
		--runtime go1.x \
		--handler memo \
		--timeout 30 \
		--memory-size 512 \
		--zip-file fileb://memo.zip \
		--role arn:aws:lambda:ap-northeast-2:000000000000:function:memo-lambda \
		--endpoint-url http://localhost:4566

remove-lambda:
	aws lambda delete-function \
	--function-name memo-lambda \
	--endpoint-url http://localhost:4566

create-event-source-mapping:
	aws lambda create-event-source-mapping --function-name memo-lambda \
	--batch-size 1 \
	--starting-position LATEST \
	--event-source-arn $(shell aws dynamodb describe-table --table-name memos \
	 --endpoint-url "http://localhost:4566" --output json --query 'Table'.'LatestStreamArn') \
	--maximum-batching-window-in-seconds 1 \
	--endpoint-url http://localhost:4566
