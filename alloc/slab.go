package alloc

import (
	"errors"
	"unsafe"
)

type slab struct {
	memory []byte
	chunks []chunk
}

func NewSlab(alloc *allocator, slabIdx int, slabSize, chunkSize int) (slab, error) {
	mem := alloc.malloc(slabSize)
	if mem == nil {
		return slab{}, errors.New("malloc fail")
	}
	chunks := make([]chunk, slabSize/chunkSize)
	for i := range chunks {
		chunks[i].loc.slabId = slabIdx
		chunks[i].loc.chunkId = i
	}
	return slab{
		memory: mem,
		chunks: chunks,
	}, nil
}

func (s *slab) Addr() uintptr {
	return uintptr(unsafe.Pointer(&s.memory[0]))
}

func (s *slab) Chunk(i int) *chunk {
	if s == nil || i >= len(s.chunks) {
		return nil
	}
	return &s.chunks[i]
}

type chunk struct {
	loc
	next loc
	refs int32
}

type loc struct {
	slabId  int
	chunkId int
}

var nilLoc = loc{-1, -1}

// IsNil returns true if the loc came from NilLoc().
func (l loc) IsNil() bool {
	return l.slabId < 0 && l.chunkId < 0
}
