package domain

type Domain struct {
	values map[int]bool
}

func New() (domain Domain) {
	domain.values = make(map[int]bool)
	return
}

func (domain Domain) Add(val int) {
	domain.values[val] = true
}

func (domain Domain) Has(val int) bool {
	_, ok := domain.values[val]
	return ok
}

func (domain Domain) Remove(val int) {
	delete(domain.values, val)
}

func (domain Domain) Size() (size int) {
	return len(domain.values)
}

func (domain Domain) Any() (int, bool) {
	for key, _ := range domain.values {
		return key, true
	}
	return 0, false
}

func (domain Domain) ForAll(cb func(int)) {
	for key, _ := range domain.values {
		cb(key)
	}
}
