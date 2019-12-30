package customevents

//RuleSpecification is the representation of a rule specification in Instana
type RuleSpecification struct {
	//Common Fields
	RuleType string `json:"ruleType"`
	Severity int    `json:"severity"`

	//System Rule fields
	SystemRuleID *string `json:"systemRuleId"`

	//Threshold Rule fields
	MetricName        *string  `json:"metricName"`
	Rollup            *int     `json:"rollup"`
	Window            *int     `json:"window"`
	Aggregation       *string  `json:"aggregation"`
	ConditionOperator *string  `json:"conditionOperator"`
	ConditionValue    *float64 `json:"conditionValue"`

	//Entity Verification Rule
	MatchingEntityType  *string `json:"matchingEntityType"`
	MatchingOperator    *string `json:"matchingOperator"`
	MatchingEntityLabel *string `json:"matchingEntityLabel"`
	OfflineDuration     *int    `json:"offlineDuration"`
}

//EventSpecificationDownstream is the representation of a downstream reporting in Instana
type EventSpecificationDownstream struct {
	IntegrationIds                []string `json:"integrationIds"`
	BroadcastToAllAlertingConfigs bool     `json:"broadcastToAllAlertingConfigs"`
}

//CustomEventSpecification is the representation of a custom event specification in Instana
type CustomEventSpecification struct {
	ID             string                        `json:"id"`
	Name           string                        `json:"name"`
	EntityType     string                        `json:"entityType"`
	Query          *string                       `json:"query"`
	Triggering     bool                          `json:"triggering"`
	Description    *string                       `json:"description"`
	ExpirationTime *int                          `json:"expirationTime"`
	Enabled        bool                          `json:"enabled"`
	Rules          []RuleSpecification           `json:"rules"`
	Downstream     *EventSpecificationDownstream `json:"downstream"`
}
