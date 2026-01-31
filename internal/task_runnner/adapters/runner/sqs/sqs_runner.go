package sqs

import (
	"context"
	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/runner"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type sqsRunner struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSRunner(client *sqs.Client, queueURL string) runner.Runner {
	return &sqsRunner{
		client:   client,
		queueURL: queueURL,
	}
}

func (r *sqsRunner) Execute(ctx context.Context, payload string) error {
	_, err := r.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(r.queueURL),
		MessageBody: aws.String(payload),
	})
	return err
}

func NewFactory() runner.Factory {
	return func(config entity.TaskRunnerConfig) runner.Runner {
		queueURL, _ := config["queue_url"].(string)
		region, _ := config["region"].(string)
		accessKey, _ := config["access_key"].(string)
		secretKey, _ := config["secret_key"].(string)
		endpoint, _ := config["endpoint"].(string)

		client := NewSQSClient(endpoint, region, accessKey, secretKey)
		return NewSQSRunner(client, queueURL)
	}
}
