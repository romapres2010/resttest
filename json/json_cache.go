package json

import (
	"sync"
)

// флаги состояния записи кеша
const (
	FROM_POOL                       uint = 0
	JSON_IN_BLOB_DB                 uint = 1
	JSON_IN_SAVE_CHANEL             uint = 2
	JSON_NOT_ENOUGH_SPACE_IN_CHANEL uint = 2
	JSON_IN_SAVE_QUEUE              uint = 4
	JSON_IN_MEMORY                  uint = 5
	JSON_INALID                     uint = 6
)

// cacheStr - структура для хранения одного объекта
type cacheStr struct {
	key  int    // ключ для БД
	buf  []byte // JSON - временное хранение до записи в БД
	flag uint   // состояния записи кеша
}

// cacheStrPool represent pooling of cacheStr
// =====================================================================
var cacheStrPool = sync.Pool{
	New: func() interface{} {
		return &cacheStr{
			key:  0,
			buf:  nil,
			flag: FROM_POOL,
		}
	},
}

// getCacheStr allocates a new cacheStr
// =====================================================================
func getCacheStr() *cacheStr {
	return cacheStrPool.Get().(*cacheStr)
}

// putCacheStr return byte buf to cache
// =====================================================================
func putCacheStr(p *cacheStr) {
	p.buf = nil
	p.key = 0
	p.flag = FROM_POOL
	cacheStrPool.Put(p)
}

// Cache represent in memory cache of Json
type Cache struct {
	ch map[int]*cacheStr
	mx sync.RWMutex
}

func (c *Cache) set(key int, buf []byte, flag uint) *cacheStr {
	c.mx.Lock()
	defer c.mx.Unlock()

	// Ишем с таким же ключем
	if v, ok := c.ch[key]; ok {
		v.key = key
		v.buf = buf
		v.flag = flag
		return v
	} else {
		vStr := getCacheStr() // возьмем новую структуру из pool
		vStr.key = key
		vStr.buf = buf
		vStr.flag = flag
		c.ch[key] = vStr
		return vStr
	}
}

func (c *Cache) get(key int) (*cacheStr, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if val, ok := c.ch[key]; ok {
		if val != nil {
			if val.flag == JSON_INALID {
				////mylog.PrintfDebug("[DEBUG] json (c *Cache) get() - JSON_INALID - val.key = '%v'", val.key)
				delete(c.ch, key) // Пустой указатель или инвалидный  - удалить из карты
				putCacheStr(val)  // вернем объект в пул
				return nil, false
			} else if val.flag == JSON_NOT_ENOUGH_SPACE_IN_CHANEL {
				////mylog.PrintfDebug("[DEBUG] json (c *Cache) get() - JSON_NOT_ENOUGH_SPACE_IN_CHANEL - val.key = '%v'", val.key)
				putCacheStr(val) // вернем объект в пул
				return nil, false
			}

			////mylog.PrintfDebug("[DEBUG] json (c *Cache) get()- data in the cache is valid - val.key = '%v', val.flag = '%v'", val.key, val.flag)
			return val, true // данные есть и они корректные

		}
		delete(c.ch, key) // Пустой указатель или инвалидный  - удалить из карты
	}
	return nil, false
}
