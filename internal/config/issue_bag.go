package config

// issueBag collects validation defects so ValidateAll can join them in one
// error. Each field check remembers its message; list() is what Join sees.
type issueBag struct {
	items []error
}

func (b *issueBag) remember(err error) {
	if err == nil {
		return
	}
	b.items = append(b.items, err)
}

func (b *issueBag) list() []error {
	if len(b.items) == 0 {
		return nil
	}
	return b.items
}
