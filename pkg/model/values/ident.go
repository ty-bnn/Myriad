package values

// Ident represents the right value defined by just a name.
// e.g. a := b の時の b
type Ident struct {
	Kind ValueKind
	Name string
}

func (i Ident) GetKind() ValueKind {
	return i.Kind
}

func (i Ident) GetName() string {
	return i.Name
}
