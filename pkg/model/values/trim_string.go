package values

type TrimString struct {
	Kind   ValueKind
	Target Value
	Trim   Value
	From   FromKind
}

func (t TrimString) GetKind() ValueKind {
	return t.Kind
}

func (t TrimString) GetName() string {
	return ""
}

type FromKind int

const (
	LEFT FromKind = iota
	RIGHT
)
