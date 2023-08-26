package storage

type Hybrid struct {
	memory *inMemory
}

func NewHybrid(inMemory *inMemory) *Hybrid {
	return &Hybrid{
		memory: inMemory,
	}
}

func (s *Hybrid) Put(id, value string) {
	s.memory.Put(id, value)
}

func (s *Hybrid) ByID(id string) (string, error) {
	return s.memory.ByID(id)
}
