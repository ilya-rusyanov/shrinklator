package storage

type FilePersistence struct {
}

func NewFilePersistence(filePath string) *FilePersistence {
	return &FilePersistence{}
}

func (p *FilePersistence) Append(string, string) {
}
