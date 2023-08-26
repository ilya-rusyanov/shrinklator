package storage

type NullPersistence struct {
}

func NewNullPersistence() *NullPersistence {
	return &NullPersistence{}
}

func (p *NullPersistence) Append(string, string) {
}

func (p *NullPersistence) ReadAll() (values map[string]string, err error) {
	values = make(map[string]string)
	return
}
