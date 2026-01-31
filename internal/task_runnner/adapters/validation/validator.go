package validation

import (
	"fmt"
	"net/url"
	"strings"

	"nimbus/internal/task_runnner/domain/entity"
	"nimbus/internal/task_runnner/domain/runner"
)

type fieldRuleValidator struct {
	rules []runner.FieldRule
}

// NewConfigValidator builds a ConfigValidator from declarative field rules.
func NewConfigValidator(rules []runner.FieldRule) runner.ConfigValidator {
	return &fieldRuleValidator{rules: rules}
}

func (v *fieldRuleValidator) Validate(config entity.TaskRunnerConfig) error {
	var missing []string

	for _, r := range v.rules {
		val, _ := config[r.Name].(string)
		if r.Required && strings.TrimSpace(val) == "" {
			missing = append(missing, r.Name)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required config fields: %s", strings.Join(missing, ", "))
	}

	for _, r := range v.rules {
		if !r.URL {
			continue
		}
		val, _ := config[r.Name].(string)
		if strings.TrimSpace(val) == "" {
			continue
		}
		if _, err := url.ParseRequestURI(val); err != nil {
			return fmt.Errorf("invalid %s: %s", r.Name, err)
		}
	}

	return nil
}
