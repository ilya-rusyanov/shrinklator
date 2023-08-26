package storage

type FilePersistence struct {
}

func NewFilePersistence(filePath string) *FilePersistence {
	return &FilePersistence{}
}

func (p *FilePersistence) Append(string, string) {
}

func (p *FilePersistence) ReadAll() (values map[string]string, err error) {
	values = make(map[string]string)
	return
}
