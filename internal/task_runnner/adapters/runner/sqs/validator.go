package sqs

import (
	"fmt"
	"net/url"
	"strings"

	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/runner"
)

type configValidator struct{}

func NewConfigValidator() runner.ConfigValidator {
	return &configValidator{}
}

func (v *configValidator) Validate(config entity.TaskRunnerConfig) error {
	var missing []string

	queueURL, _ := config["queue_url"].(string)
	region, _ := config["region"].(string)
	accessKey, _ := config["access_key"].(string)
	secretKey, _ := config["secret_key"].(string)

	if strings.TrimSpace(queueURL) == "" {
		missing = append(missing, "queue_url")
	}
	if strings.TrimSpace(region) == "" {
		missing = append(missing, "region")
	}
	if strings.TrimSpace(accessKey) == "" {
		missing = append(missing, "access_key")
	}
	if strings.TrimSpace(secretKey) == "" {
		missing = append(missing, "secret_key")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	if _, err := url.ParseRequestURI(queueURL); err != nil {
		return fmt.Errorf("invalid queue_url: %s", err)
	}

	if endpoint, ok := config["endpoint"].(string); ok && strings.TrimSpace(endpoint) != "" {
		if _, err := url.ParseRequestURI(endpoint); err != nil {
			return fmt.Errorf("invalid endpoint: %s", err)
		}
	}

	return nil
}
