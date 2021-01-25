package handler

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/vranystepan/email/pkg/service"
)

// VerifyParams contains params for Verify function to minify signature
type VerifyParams struct {
	SES service.SES
}

// Verify sends verification email via SES
func Verify(p VerifyParams) func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		return nil
	}
}
