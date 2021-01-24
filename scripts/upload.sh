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
echo "uploading ${SERVICE}"

# get s3 bucket name from cloud formation
S3_BUCKET=$(aws cloudformation describe-stacks --stack-name artifacts --query "Stacks[0].Outputs[?OutputKey=='BucketName'].OutputValue" --output text)

# upload code to the s3 bucket
cd ./bin/${SERVICE}
chmod +x main
zip ${SERVICE}.zip main
aws s3 cp ${SERVICE}.zip s3://${S3_BUCKET}/

# get the version ID
aws s3api list-object-versions --bucket ${S3_BUCKET} --query "reverse(sort_by(Versions, &LastModified))[0].VersionId" --output text
