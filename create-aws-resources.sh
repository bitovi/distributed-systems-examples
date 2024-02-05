#!/bin/bash

API_URL=${API_URL}
AWS_URL=${AWS_URL}
ORDER_TIMEOUT_TOPIC=${ORDER_TIMEOUT_TOPIC}
POSTGRES_HOST=${POSTGRES_HOST}
POSTGRES_USER=${POSTGRES_USER}
POSTGRES_PASS=${POSTGRES_PASS}
POSTGRES_PORT=${POSTGRES_PORT}

cat <<EOF > admin-policy.json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "*",
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF

echo "Create admin iam role"
aws \
  iam create-role \
  --endpoint-url http://localhost:4566
  --region us-east-1 \
  --role-name admin-role \
  --path / \
  --assume-role-policy-document file://admin-policy.json

echo "Make an S3 bucket for the lambdas"
aws \
  s3 mb s3://lambda-functions \
  --endpoint-url http://localhost:4566

echo "Copy the lambda function to the S3 bucket"
aws \
  s3 cp lambdas.zip s3://lambda-functions \
  --endpoint-url http://localhost:4566

echo "Create Order Transmission SNS Topic"
aws \
  sns create-topic \
  --name orderTransmission \
  --endpoint-url http://localhost:4566

echo "Create Order Transmission SQS Queue"
aws \
  sqs create-queue \
  --queue-name orderTransmission \
  --endpoint-url http://localhost:4566

echo "Subscribe Order Transmission Topic to Order Transmission Queue"
aws \
  sns subscribe \
  --endpoint-url http://localhost:4566 \
  --topic-arn arn:aws:sns:us-east-1:000000000000:orderTransmission \
  --protocol sqs \
  --notification-endpoint arn:aws:sqs:us-east-1:000000000000:orderTransmission

echo "Create Order Tranmsission Lambda"
aws \
  lambda create-function \
  --endpoint-url http://localhost:4566 \
  --function-name orderTransmission \
  --role arn:aws:iam::000000000000:role/admin-role \
  --code S3Bucket=lambda-functions,S3Key=lambdas.zip \
  --handler transmission.handler \
  --runtime nodejs20.x \
  --description "Order Transmission Lambda Function" \
  --environment Variables="{API_URL=$API_URL}" \
  --timeout 60 \
  --memory-size 128

echo "Create Order Timeout SNS Topic"
aws \
  sns create-topic \
  --endpoint-url http://localhost:4566 \
  --name orderTimeout

echo "Create Order Timeout SQS Queue"
aws \
  sqs create-queue \
  --endpoint-url http://localhost:4566 \
  --queue-name orderTimeout

echo "Subscribe Order Timeout Topic to Order Timeout Queue"
aws \
  sns subscribe \
  --endpoint-url http://localhost:4566 \
  --topic-arn arn:aws:sns:us-east-1:000000000000:orderTimeout \
  --protocol sqs \
  --notification-endpoint arn:aws:sqs:us-east-1:000000000000:orderTimeout

echo "Create Order Timeout Lambda"
aws \
  lambda create-function \
  --endpoint-url http://localhost:4566 \
  --function-name orderTimeout \
  --role arn:aws:iam::000000000000:role/admin-role \
  --code S3Bucket=lambda-functions,S3Key=lambdas.zip \
  --handler timeout.handler \
  --runtime nodejs20.x \
  --description "Order Timeout Lambda Function" \
  --environment Variables="{API_URL=$API_URL}" \
  --timeout 60 \
  --memory-size 128

echo "Create Order Timeout Poller Lambda"
aws \
  lambda create-function \
  --endpoint-url http://localhost:4566 \
  --function-name orderTimeoutPoller \
  --role arn:aws:iam::000000000000:role/admin-role \
  --code S3Bucket=lambda-functions,S3Key=lambdas.zip \
  --handler timeoutPoller.handler \
  --runtime nodejs20.x \
  --description "Order Timeout Poller Lambda Function" \
  --environment Variables="{AWS_URL=$AWS_URL,ORDER_TIMEOUT_TOPIC=$ORDER_TIMEOUT_TOPIC,POSTGRES_PORT=$POSTGRES_PORT,POSTGRES_HOST=$POSTGRES_HOST,POSTGRES_USER=$POSTGRES_USER,POSTGRES_PASS=$POSTGRES_PASS}" \
  --timeout 60 \
  --memory-size 128

echo "Create EventBridge rule to run every minute"
aws \
  events put-rule \
  --endpoint-url http://localhost:4566 \
  --name "OrderTimeoutPollerRule" \
  --schedule-expression "rate(1 minute)"

echo "Grant EventBridge service principal permission to run the Lambda"
aws \
  lambda add-permission \
  --endpoint-url http://localhost:4566 \
  --statement-id "InvokeOrderTimeoutPoller" \
  --action "lambda:InvokeFunction" \
  --principal "events.amazonaws.com" \
  --function-name orderTimeoutPoller \
  --source-arn "arn:aws:events:us-east-1:000000000000:rule/OrderTimeoutPollerRule"

echo "Add Lambda function as a target of the EventBridge rule"
aws \
  events put-targets \
  --endpoint-url http://localhost:4566 \
  --rule OrderTimeoutPollerRule \
  --targets "Id"="1","Arn"="arn:aws:lambda:us-east-1:000000000000:function:orderTimeoutPoller"

echo "Map Order Timeout Lambda to Order Timeout Queue"
aws \
  lambda create-event-source-mapping \
  --endpoint-url http://localhost:4566 \
  --function-name orderTimeout \
  --batch-size 1 \
  --event-source-arn "arn:aws:sqs:us-east-1:000000000000:orderTimeout"

echo "Map Order Transmission Lambda to Order Transmission Queue"
aws \
  lambda create-event-source-mapping \
  --endpoint-url http://localhost:4566 \
  --function-name orderTransmission \
  --batch-size 1 \
  --event-source-arn "arn:aws:sqs:us-east-1:000000000000:orderTransmission"