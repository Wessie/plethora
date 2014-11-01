package scheduler

import (
	"encoding/binary"
	"sync"
	"sync/atomic"
	"time"
)

func newCodec() *codec {
	return &codec{
		mapping: make(map[uint64]Task, 8),
	}
}

// codec implements the encoding and decoding functions required
// for storing our schedule in bolt key/values
type codec struct {
	// id is the ID to be used in encodeTask
	id uint64

	// mapping is a map between integer IDs and tasks
	mapping map[uint64]Task
	// mu protects the mapping map
	mu sync.Mutex
}

func (c *codec) decodeTime(b []byte) time.Time {
	t, err := time.Parse(string(b), time.RFC3339Nano)
	if err != nil {
		panic(err)
	}
	return t
}

func (c *codec) encodeTime(t time.Time) []byte {
	return []byte(t.Format(time.RFC3339Nano))
}

func (c *codec) encodeTask(t Task) []byte {
	id := atomic.AddUint64(&c.id, 1)

	c.mu.Lock()
	c.mapping[id] = t
	c.mu.Unlock()

	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, id)
	return b
}

func (c *codec) decodeTask(b []byte) Task {
	id := binary.BigEndian.Uint64(b)

	c.mu.Lock()
	tsk, ok := c.mapping[id]
	c.mu.Unlock()

	if !ok {
		panic("unknown task")
	}

	return tsk
}
