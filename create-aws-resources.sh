#!/bin/bash

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
 --endpoint-url=http://localhost:4566 \
 --region us-east-1 \
 iam create-role \
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
echo "Create Order Transmission SQS Queue"
echo "Subscribe Order Transmission Topic to Order Transmission Queue"
echo "Create Order Tranmsission Lambda"
aws \
  lambda create-function \
  --endpoint-url=http://localhost:4566 \
  --function-name orderTransmission \
  --role arn:aws:iam::000000000000:role/admin-role \
  --code S3Bucket=lambda-functions,S3Key=lambdas.zip \
  --handler transmission.handler \
  --runtime nodejs20.x \
  --description "Order Transmission Lambda Function" \
  --timeout 60 \
  --memory-size 128
echo "Map Order Transmission Lambda to Order Transmission Queue"

echo "Create Order Timeout SNS Topic"
echo "Create Order Timeout SQS Queue"
echo "Subscribe Order Timeout Topic to Order Timeout Queue"
echo "Create Order Timeout Lambda"
aws \
  lambda create-function \
  --endpoint-url=http://localhost:4566 \
  --function-name orderTimeout \
  --role arn:aws:iam::000000000000:role/admin-role \
  --code S3Bucket=lambda-functions,S3Key=lambdas.zip \
  --handler timeout.handler \
  --runtime nodejs20.x \
  --description "Order Timeout Lambda Function" \
  --timeout 60 \
  --memory-size 128
echo "Map Order Timeout Lambda to Order Timeout Queue"

echo "Create Order Timeout Poller Lambda"
aws \
  lambda create-function \
  --endpoint-url=http://localhost:4566 \
  --function-name orderTimeoutPoller \
  --role arn:aws:iam::000000000000:role/admin-role \
  --code S3Bucket=lambda-functions,S3Key=lambdas.zip \
  --handler timeoutPoller.handler \
  --runtime nodejs20.x \
  --description "Order Timeout Poller Lambda Function" \
  --timeout 60 \
  --memory-size 128
echo "Create CloudWatch event with a cron schedule"
echo "Add Lambda function as a target of the CloudWatch event"
