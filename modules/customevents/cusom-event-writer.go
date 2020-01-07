package customevents

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

//NewCustomEventWriter Creates a new instance of the CustomEventWriter
func NewCustomEventWriter(buffer *bytes.Buffer) CustomEventWriter {
	return &customEventWriterImpl{
		buffer: buffer,
	}
}

//CustomEventWriter Interface definition for the CustomEventWriter
type CustomEventWriter interface {
	//Write writes an array/slice of CustomEventSpecification
	Write(events []*CustomEventSpecification) error
}

type customEventWriterImpl struct {
	buffer *bytes.Buffer
}

func (w *customEventWriterImpl) Write(events []*CustomEventSpecification) error {
	for _, event := range events {
		err := w.writeEventSpecification(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *customEventWriterImpl) writeEventSpecification(event *CustomEventSpecification) error {
	rule := event.Rules[0]
	resourceType, err := w.getTerraformResourceType(event, &rule)
	if err != nil {
		return err
	}

	resourceName, err := w.getResourceName(event)
	if err != nil {
		return err
	}

	w.appendString(fmt.Sprintf("resource \"%s\" \"%s\" {\n", resourceType, resourceName))
	w.appendString(fmt.Sprintf("  name = \"%s\"\n", event.Name))
	if rule.RuleType == "threshold" {
		w.appendString(fmt.Sprintf("  entity_type = \"%s\"\n", event.EntityType))
	}
	if event.Query != nil {
		w.appendString(fmt.Sprintf("  query = \"%s\"\n", strings.ReplaceAll(*event.Query, "\"", "\\\"")))
	}
	w.appendString(fmt.Sprintf("  enabled = %t\n", event.Enabled))
	w.appendString(fmt.Sprintf("  triggering = %t\n", event.Triggering))
	if event.Description != nil {
		w.appendString(fmt.Sprintf("  description = \"%s\"\n", strings.ReplaceAll(strings.ReplaceAll(*event.Description, "\"", "\\\""), "\n", "\\n")))
	}
	w.appendString(fmt.Sprintf("  expiration_time = %d\n", *event.ExpirationTime))
	err = w.writeRule(event, &rule)
	if err != nil {
		return err
	}

	if event.Downstream != nil {
		w.writeDownstreamIntegration(event.Downstream)
	}
	w.appendString(fmt.Sprintf("}\n\n"))
	return nil
}

func (w *customEventWriterImpl) writeRule(event *CustomEventSpecification, rule *RuleSpecification) error {
	severity, err := w.formatRuleSeverity(event, rule)
	if err != nil {
		return err
	}
	w.appendString(fmt.Sprintf("  rule_severity = \"%s\"\n", severity))
	ruleType := rule.RuleType
	if ruleType == "system" {
		w.writeSystemRule(rule)
	} else if ruleType == "threshold" {
		w.writeThresholdRule(rule)
	} else if ruleType == "entity_verification" {
		w.writeEntityVerificationRule(rule)
	} else {
		return fmt.Errorf("Error writing event specification '%s' (ID %s): Unsupported rule type %s", event.Name, event.Name, ruleType)
	}
	return nil
}

func (w *customEventWriterImpl) writeSystemRule(rule *RuleSpecification) {
	if rule.SystemRuleID != nil {
		w.appendString(fmt.Sprintf("  rule_system_rule_id = \"%s\"\n", *rule.SystemRuleID))
	}
}

func (w *customEventWriterImpl) writeThresholdRule(rule *RuleSpecification) {
	if rule.MetricName != nil {
		w.appendString(fmt.Sprintf("  rule_metric_name = \"%s\"\n", *rule.MetricName))
	}
	if rule.Window != nil {
		w.appendString(fmt.Sprintf("  rule_window = %d\n", *rule.Window))
	}
	if rule.Rollup != nil {
		w.appendString(fmt.Sprintf("  rule_rollup = %d\n", *rule.Rollup))
	}
	if rule.Aggregation != nil {
		w.appendString(fmt.Sprintf("  rule_aggregation = \"%s\"\n", *rule.Aggregation))
	}
	if rule.ConditionOperator != nil {
		w.appendString(fmt.Sprintf("  rule_condition_operator = \"%s\"\n", *rule.ConditionOperator))
	}
	if rule.ConditionValue != nil {
		w.appendString(fmt.Sprintf("  rule_condition_value = %f\n", *rule.ConditionValue))
	}
}

func (w *customEventWriterImpl) writeEntityVerificationRule(rule *RuleSpecification) {
	if rule.MatchingEntityType != nil {
		w.appendString(fmt.Sprintf("  rule_matching_entity_type = \"%s\"\n", *rule.MatchingEntityType))
	}
	if rule.ConditionOperator != nil {
		w.appendString(fmt.Sprintf("  rule_matching_operator = \"%s\"\n", *rule.MatchingOperator))
	}
	if rule.MatchingEntityLabel != nil {
		w.appendString(fmt.Sprintf("  rule_matching_entity_label = \"%s\"\n", *rule.MatchingEntityLabel))
	}
	if rule.OfflineDuration != nil {
		w.appendString(fmt.Sprintf("  rule_offline_duration = %d\n", *rule.OfflineDuration))
	}
}

func (w *customEventWriterImpl) formatRuleSeverity(event *CustomEventSpecification, rule *RuleSpecification) (string, error) {
	severity := rule.Severity
	if severity == 5 {
		return "warning", nil
	} else if severity == 10 {
		return "critical", nil
	} else {
		return "", fmt.Errorf("Error writing event specification '%s' (ID %s): Unsupported severity %d", event.Name, event.ID, severity)
	}
}

func (w *customEventWriterImpl) writeDownstreamIntegration(ds *EventSpecificationDownstream) {
	w.appendString(fmt.Sprintf("  downstream_integration_ids = %+q\n", ds.IntegrationIds))
	w.appendString(fmt.Sprintf("  downstream_broadcast_to_all_alerting_configs = %t\n", ds.BroadcastToAllAlertingConfigs))
}

func (w *customEventWriterImpl) getTerraformResourceType(event *CustomEventSpecification, rule *RuleSpecification) (string, error) {
	ruleType := rule.RuleType
	if ruleType == "system" {
		return "instana_custom_event_spec_system_rule", nil
	} else if ruleType == "threshold" {
		return "instana_custom_event_spec_threshold_rule", nil
	} else if ruleType == "entity_verification" {
		return "instana_custom_event_spec_entity_verification_rule", nil
	} else {
		return "", fmt.Errorf("Error writing event specification '%s' (ID %s): Unsupported rule type %s", event.Name, event.ID, ruleType)
	}
}

func (w *customEventWriterImpl) getResourceName(spec *CustomEventSpecification) (string, error) {
	alphaNumericOnlyRegexp, err := regexp.Compile("[^a-zA-Z0-9]")
	if err != nil {
		return "", err
	}
	multiUnderscoreRegexp, err := regexp.Compile("_+")
	if err != nil {
		return "", err
	}
	return multiUnderscoreRegexp.ReplaceAllString(alphaNumericOnlyRegexp.ReplaceAllString(spec.Name, "_"), "_"), nil
}

func (w *customEventWriterImpl) appendString(s string) {
	w.buffer.WriteString(s)
}
