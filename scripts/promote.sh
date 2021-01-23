#!/bin/sh

# make sure SERVICE is not propagated from the parent
unset SERVICE

while getopts 's:' c
do
  case $c in
    s) SERVICE=$OPTARG ;;
  esac
done

[ -z "${SERVICE}" ] && { echo "SERVICE (-s) can't be empty!"; exit 1; }
echo "promoting ${SERVICE}"

S3_BUCKET=$(aws cloudformation describe-stacks --stack-name artifacts --query "Stacks[0].Outputs[?OutputKey=='BucketName'].OutputValue" --output text)
LAMBDA_VERSION=$(aws s3api list-object-versions --bucket ${S3_BUCKET} --query "reverse(sort_by(Versions, &LastModified))[?Key=='${SERVICE}.zip'] | [0].VersionId" --output text)
LAMBDA_NAME=$(aws cloudformation describe-stacks --stack-name "email-${SERVICE}" --query "Stacks[0].Outputs[?OutputKey=='FunctionName'].OutputValue" --output text)
GATEWAY_ID=$(aws cloudformation describe-stacks --stack-name gateway --query "Stacks[0].Outputs[?OutputKey=='GatewayID'].OutputValue" --output text)

echo "setting ${LAMBDA_NAME} to artiface version ${LAMBDA_VERSION}"
aws lambda update-function-code --function-name "${LAMBDA_NAME}" --s3-object-version "${LAMBDA_VERSION}" --s3-bucket="${S3_BUCKET}" --s3-key="${SERVICE}.zip"

echo "promoting API gateway"
aws apigateway create-deployment --rest-api-id "${GATEWAY_ID}" --stage-name main
