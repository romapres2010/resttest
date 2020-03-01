package cache

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/romapres2010/resttest/model"
)

// deptStr - структура для хранения одного объекта
type deptStr struct {
	val        *model.Dept
	isOutdate  bool // признак, что данные устрели и требуют перегрузки
	isToRemove bool // признак, что данные к удалению
	getsCnt    uint // счетчик количества обращений
}

// empStr - структура для хранения одного объекта
type empStr struct {
	val        *model.Emp
	isOutdate  bool // признак, что данные устрели и требуют перегрузки
	isToRemove bool // признак, что данные к удалению
	getsCnt    uint // счетчик количества обращений
}

type deptCache struct {
	ch map[int]*deptStr
	mx sync.RWMutex // для независимой блокировки отдельных кешей
}

func (c *deptCache) get(id int) (*model.Dept, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	vStr, ok := c.ch[id]
	if ok {
		val := vStr.val
		if vStr != nil && !vStr.isOutdate && !vStr.isToRemove {
			if val != nil && !val.IsPutToPool {
				//mylog.PrintfDebug("[DEBUG] cache (c *deptCache) get()- data in the cache is valid, p = '%p'", val)
				return val, true // данные есть и они корректные
			}
			//mylog.PrintfDebug("[DEBUG] cache (c *deptCache) get() - data in the cache is INVALID - vStr.isOutdate='%v', vStr.isToRemove='%v', vStr.val.IsPutToPool='%v'", vStr.isOutdate, vStr.isToRemove, vStr.val.IsPutToPool)
		}
		delete(c.ch, id)          // Пустой указатель или инвалидный - удалить из карты
		model.PutDept(val, false) // вернем объект в пул
		return nil, false
	}
	return nil, false
}

func (c *deptCache) set(id int, v *model.Dept) {
	if v != nil {
		v.IsPutToPool = false // признак, что структура валидная - можно использовать в кэш
		c.mx.Lock()
		defer c.mx.Unlock()
		vStr, ok := c.ch[id]
		if ok {
			vStr.val = v
			vStr.isOutdate = false
			vStr.isToRemove = false
		} else {
			c.ch[id] = &deptStr{
				val:        v,
				isOutdate:  false,
				isToRemove: false,
			}
		}
	}
	return
}

func (c *deptCache) outdate(id int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	vStr, ok := c.ch[id]
	if ok {
		delete(c.ch, id)               // удалим из карты
		model.PutDept(vStr.val, false) // вернем объект в пул, без подчиненных записей
		vStr.val = nil
	} else {
		//log.Printf("[ERROR] cache (c *deptCache) outdate - id = '%v' not in the cache", id)
	}
	return
}

type empCache struct {
	ch map[int]*empStr
	mx sync.RWMutex // для независимой блокировки отдельных кешей
}

func (c *empCache) get(id int) (*model.Emp, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	vStr, ok := c.ch[id]
	if ok {
		val := vStr.val
		if vStr != nil && !vStr.isOutdate && !vStr.isToRemove {
			if val != nil && !val.IsPutToPool {
				//mylog.PrintfDebug("[DEBUG] cache (c *empCache) get()- data in the cache is valid, p = '%p'", val)
				return val, true // данные есть и они корректные
			}
			//mylog.PrintfDebug("[DEBUG] cache (c *empCache) get() - data in the cache is INVALID - vStr.isOutdate='%v', vStr.isToRemove='%v', vStr.val.IsPutToPool='%v'", vStr.isOutdate, vStr.isToRemove, vStr.val.IsPutToPool)
		}
		delete(c.ch, id)  // Пустой указатель или инвалидный  - удалить из карты
		model.PutEmp(val) // вернем объект в пул
		return nil, false
	}
	return nil, false
}

func (c *empCache) set(id int, v *model.Emp) {
	if v != nil {
		v.IsPutToPool = false // признак, что структура валидная - можно использовать в кэш
		c.mx.Lock()
		defer c.mx.Unlock()
		vStr, ok := c.ch[id]
		if ok {
			vStr.val = v
			vStr.val.IsPutToPool = false // признак, что структура валидная - можно использовать в кэш
			vStr.isOutdate = false
			vStr.isToRemove = false
		} else {
			c.ch[id] = &empStr{
				val:        v,
				isOutdate:  false,
				isToRemove: false,
			}
		}
	}
	return
}

func (c *empCache) outdate(id int) {
	c.mx.Lock()
	defer c.mx.Unlock()
	vStr, ok := c.ch[id]
	if ok {
		/*
			if vStr.val != nil && vStr.val.Dept != nil {
				vStr.val.Dept.outdate()
			}
		*/
		delete(c.ch, id)       // удалим из карты
		model.PutEmp(vStr.val) // вернем объект в пул
		vStr.val = nil
	} else {
		//		log.Printf("[ERROR] cache (c *empCache) outdate - id = '%v' not in the cache", id)
	}
	return
}

// Cache wraps a model.Services to provide an in-memory cache
type Cache struct {
	depts         deptCache
	emps          empCache
	empService    model.EmpService
	deptService   model.DeptService
	modelListener *model.Listener // обработчик уведомлений
}

// New returns a new read-through cache for service.
// =====================================================================
func New(empService model.EmpService, deptService model.DeptService) *Cache {
	return &Cache{
		depts: deptCache{
			ch: make(map[int]*deptStr),
		},
		emps: empCache{
			ch: make(map[int]*empStr),
		},
		empService:  empService,
		deptService: deptService,
	}
}

// SetModelListener set model Listener for notifications
// =====================================================================
func (c *Cache) SetModelListener(l *model.Listener) {
	cnt := "cache (c *Cache) SetModelListener" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	c.modelListener = l
}

// PrintElapsed при отложенном запуске в начале процедуры измеряет время ее выполнение (с закрытием)
func PrintElapsed(what string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s%v", what, time.Since(start))
	}
}

// PopulateCache populate all caches
// =====================================================================
func (c *Cache) PopulateCache() {
	cnt := "cache (c *Cache) PopulateCache" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	// считаем все PK для dept
	deptsPK, err := c.GetDeptsPK()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.GetDeptsPK()", cnt)
		log.Printf(errM)
		return
	}

	//	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	log.Printf("[INFO] %v - runtime.NumCPU()= '%v'", cnt, numCPU)

	lenPK := len(deptsPK)
	log.Printf("[INFO] %v - lenPK= '%v'", cnt, lenPK)

	var waitgroup sync.WaitGroup // для ожидания всех запущенных горутин

	for i := 0; i < numCPU; i++ {
		log.Printf("[INFO] %v - start gorutine '%v'", cnt, i)
		waitgroup.Add(1) // добавляем в список ожидания горутин

		go func(index int) {
			defer PrintElapsed(fmt.Sprintf("[INFO] %v - gorutine '%v' - takes", cnt, index))()
			for j := 0; j < lenPK/numCPU; j++ {
				pk := deptsPK[j+(index*(lenPK/numCPU))].Deptno
				_, _ = c.GetDeptEmps(pk)
			}
			waitgroup.Done() // сигнал завершения горутины
		}(i)

	}
	waitgroup.Wait() // ожидание завершения всех горутин

	log.Printf("[INFO] %v - SUCCESS", cnt)
	return
}

//ProcessNotification - process notification and invalidate cache
// =====================================================================
func (c *Cache) ProcessNotification(n *model.Notification) error {
	cnt := "cache (c *Cache) processNotification" // Имя текущего метода для логирования

	if n.ObjectName == "dept" {
		if n.OperationName == "outDate" {
			//mylog.PrintfDebug("[DEBUG] %v - set outDate for deptno = '%v'", cnt, n.Id)
			c.depts.outdate(n.Id)
			return nil
		}
	}
	if n.ObjectName == "emp" {
		if n.OperationName == "outDate" {
			//mylog.PrintfDebug("[DEBUG] %v - set outDate for empno = '%v'", cnt, n.Id)
			c.emps.outdate(n.Id)
			return nil
		}
	}
	errM := fmt.Sprintf("[ERROR] %v - UNKNOWUN notification'%+v'", cnt, n)
	log.Printf(errM)
	return errors.New(errM)
}
