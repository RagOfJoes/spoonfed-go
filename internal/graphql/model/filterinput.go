package model

type FilterInput struct {
	Key       string          `json:"key"`
	Condition FilterCondition `json:"condition"`
	Value     interface{}     `json:"value"`
	Values    []interface{}   `json:"values"`
}

func (f FilterInput) HasValidCondition(conditions []FilterCondition) bool {
	if len(conditions) == 0 {
		return true
	}

	condition := f.Condition
	for _, vc := range conditions {
		if vc == condition {
			return true
		}
	}
	return false
}
