package storage

type Persistence interface {
	Append(short, long string)
}

type Hybrid struct {
	memory      *inMemory
	persistence Persistence
}

func NewHybrid(inMemory *inMemory, persistence Persistence) *Hybrid {
	return &Hybrid{
		memory:      inMemory,
		persistence: persistence,
	}
}

func (s *Hybrid) Put(id, value string) {
	s.memory.Put(id, value)
	s.persistence.Append(id, value)
}

func (s *Hybrid) ByID(id string) (string, error) {
	return s.memory.ByID(id)
}
