# Email subscription platform

![Register service](https://github.com/vranystepan/email/workflows/Register%20service/badge.svg)
[![Codeac](https://static.codeac.io/badges/2-332261928.svg "Codeac")](https://app.codeac.io/github/vranystepan/email)
[![Go Report Card](https://goreportcard.com/badge/github.com/vranystepan/email)](https://goreportcard.com/report/github.com/vranystepan/email)

I'd like to own my audience as per [Always Own Your platform](https://www.alwaysownyourplatform.com). Hence I'm preparing this simple platform for
emailing that's gonna help me with registration and stuff like such.

And as I'm AWS partner, It's gonna be built from AWS components. Also
it's gonna be a bit over engineered because I'd like to play with some
services.

## Stack

- **AWS Lambda** for all backend operations
- **AWS SQS** for messaging
- **AWS Simple Email Service** for the delivering of e-mail messages
- **AWS DynamoDB** for the persistence
- **AWS API Gateway**
- **AWS S3**

## Infrastructure

All the infrastructure components are provisioned with CloudFormation.
You can find the concrete instructions below.

### Artifacts

```bash
aws cloudformation deploy --stack-name artifacts --template-file ./infrastructure/cloudformation/artifacts.yaml --no-fail-on-empty-changeset
```

### Create register-email resources

```bash
S3_BUCKET=$(aws cloudformation describe-stacks \
  --stack-name artifacts \
  --query "Stacks[0].Outputs[?OutputKey=='BucketName'].OutputValue" \
  --output text)

aws cloudformation deploy \
    --stack-name email-register \
    --template-file ./infrastructure/cloudformation/email-register.yaml \
    --capabilities CAPABILITY_IAM \
    --no-fail-on-empty-changeset \
    --parameter-overrides \
        BucketName="${S3_BUCKET}"
```

### Create issue resources

```bash
LAMBDA_REGISTER_EMAIL_QUEUE_ARN=$(aws cloudformation describe-stacks \
  --stack-name email-register \
  --query "Stacks[0].Outputs[?OutputKey=='QueueARN'].OutputValue" \
  --output text)

aws cloudformation deploy \
    --stack-name email-issue \
    --template-file ./infrastructure/cloudformation/email-issue.yaml \
    --capabilities CAPABILITY_IAM \
    --no-fail-on-empty-changeset \
    --parameter-overrides \
        BucketName="${S3_BUCKET}" \
        QueueARN="${LAMBDA_REGISTER_EMAIL_QUEUE_ARN}"
```

### Create verify resources

```bash
LAMBDA_ISSUE_VERIFICATION_QUEUE_ARN=$(aws cloudformation describe-stacks \
  --stack-name email-issue \
  --query "Stacks[0].Outputs[?OutputKey=='QueueARN'].OutputValue" \
  --output text)

aws cloudformation deploy \
    --stack-name email-verify \
    --template-file ./infrastructure/cloudformation/email-verify.yaml \
    --capabilities CAPABILITY_IAM \
    --no-fail-on-empty-changeset \
    --parameter-overrides \
        BucketName="${S3_BUCKET}" \
        EmailTemplateBucketName="${S3_BUCKET}" \
        QueueARN="${LAMBDA_ISSUE_VERIFICATION_QUEUE_ARN}" \
        SenderEmailAddress="stepan@vrany.dev"
```

### Create gateway

```bash
LAMBDA_REGISTER_EMAIL_ARN=$(aws cloudformation describe-stacks \
  --stack-name email-register \
  --query "Stacks[0].Outputs[?OutputKey=='FunctionARN'].OutputValue" \
  --output text)

LAMBDA_REGISTER_EMAIL_NAME=$(aws cloudformation describe-stacks \
  --stack-name email-register \
  --query "Stacks[0].Outputs[?OutputKey=='FunctionName'].OutputValue" \
  --output text)

aws cloudformation deploy \
    --stack-name gateway \
    --template-file ./infrastructure/cloudformation/gateway.yaml \
    --no-fail-on-empty-changeset \
    --parameter-overrides \
        RegisterFunctionARN="${LAMBDA_REGISTER_EMAIL_ARN}"
```

## Services

```bash
export SERVICE=register
./scripts/build.sh -s "${SERVICE}"
./scripts/upload.sh -s "${SERVICE}"
./scripts/promote.sh -s "${SERVICE}"
```
