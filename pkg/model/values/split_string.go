package values

type SplitString struct {
	Kind   ValueKind
	Target Value
	Sep    Value
}

func (s SplitString) GetKind() ValueKind {
	return s.Kind
}

func (s SplitString) GetName() string {
	return ""
}
