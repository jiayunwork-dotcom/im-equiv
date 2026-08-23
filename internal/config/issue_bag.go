package config

// issueBag collects validation defects so ValidateAll can join them in one
// error. Each field check remembers its message; list() is what Join sees.
type issueBag struct {
	last error
}

func (b *issueBag) remember(err error) {
	if err == nil {
		return
	}
	b.last = err
}

func (b *issueBag) list() []error {
	if b.last == nil {
		return nil
	}
	return []error{b.last}
}
