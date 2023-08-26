package storage

type NullPersistence struct {
}

func NewNullPersistence() *NullPersistence {
	return &NullPersistence{}
}

func (p *NullPersistence) Append(string, string) {
}
