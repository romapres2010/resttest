package bolt

import (
	"sync"
)

type Counter struct {
	// на время тестирования эффективности
	countGetByteSlice uint // количество раз запроса кэша
	countPutByteSlice uint // количество раз возврата кэша
	countNewByteSlice uint // количество раз запроса кэша
	mx                sync.RWMutex
}

var stat Counter

// byteSlicePool represent pooling of []byte
// =====================================================================
var byteSlicePool = sync.Pool{
	New: func() interface{} {
		stat.mx.Lock()
		stat.countNewByteSlice++
		stat.mx.Unlock()

		return make([]byte, 512)
	},
}

// getByteSlice allocates a new []byte
// =====================================================================
func getByteSlice() []byte {
	stat.mx.Lock()
	stat.countGetByteSlice++
	stat.mx.Unlock()

	return byteSlicePool.Get().([]byte)
}

// putCacheStr return byte buf to cache
// =====================================================================
func putByteSlice(p []byte) {
	stat.mx.Lock()
	stat.countPutByteSlice++
	stat.mx.Unlock()

	p = p[:0]
	byteSlicePool.Put(p)
}
