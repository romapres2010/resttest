package json

import (
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/romapres2010/restdept/bolt"
)

// Listener represents a chanel for save to Bolt DB
// =====================================================================
type Listener struct {
	db        *bolt.Db
	ch        chan *cacheStr // канал для формирования пакетной записи в BOLT
	slice     bolt.SaveSlice // очередь для записи в BOLT
	orig      []*cacheStr    // очередь оригинальных элементов
	queueSize int            // размер очереди записи в Bolt
	bucket    string         // имя bucket для записи
}

// newSaveListener create new listener for save to Bolt DB
// =====================================================================
func newSaveListener(db *bolt.Db, bucket string, chSize int, queueSize int) *Listener {
	cnt := "json newSaveListener" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	if db == nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - Bolt DB pointer IS NULL", cnt)
		log.Printf(errM)
		return nil
	}
	if bucket == "" {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - Bolt DB Bucket Name IS NULL", cnt)
		log.Printf(errM)
		return nil
	}

	l := &Listener{
		db:        db,
		ch:        make(chan *cacheStr, chSize),
		slice:     make([]*bolt.SaveStr, queueSize),
		orig:      make([]*cacheStr, queueSize),
		queueSize: queueSize,
		bucket:    bucket,
	}

	l.slice = l.slice[:0] // сбрасываем размер очереди
	l.orig = l.orig[:0]   // сбрасываем размер очереди

	// Запускем бесконечный цикл с recovery
	go l.startListener()

	log.Printf("[INFO] %v - SUCCESS", cnt)
	return l
}

// startListener - запускаем бесконечный цикл для ожидания уведомлений
// =====================================================================
func (l *Listener) startListener() {
	cnt := "json (l *Listener) startListener()" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	for {
		l.waitForSaveMessage()
	}
}

// waitForSaveMessage - ожидание уведомления
// =====================================================================
func (l *Listener) waitForSaveMessage() {
	cnt := "json (l *Listener) waitForSaveMessage" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// Восстановление после паники, иначе падает весь сервер
	var err error
	defer func() {
		r := recover()
		if r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = errors.New("Unknown error")
			}
			log.Printf("[ERROR] %v - recovery from panic \n %+v", cnt, err.Error())
		}
	}()

	// Обрабока очереди
	select {
	case m := <-l.ch:
		//mylog.PrintfDebug("[DEBUG] %v - recieved message - len(l.slice) = '%v'", cnt, len(l.slice))

		if m == nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - m := <-l.ch - m IS NULL", cnt)
			log.Printf(errM)
		}

		l.orig = append(l.orig, m) // добавляем сообщение в очередь оригинальных элементов

		// структура для записи в БД
		v := &bolt.SaveStr{
			Key: m.key,
			Buf: m.buf,
		}
		l.slice = append(l.slice, v) // добавляем сообщение в очередь для записи
		m.flag = JSON_IN_SAVE_QUEUE  // устанавливаем флаг, что поместили в очередь

		// Если очередь для записи наполнилась - отправляем на запись в синхронном режиме
		if len(l.slice) >= l.queueSize {
			//mylog.PrintfDebug("[DEBUG] %v - len(l.slice) >= l.queueSize '%v' >= '%v'", cnt, len(l.slice), l.queueSize)
			l.SaveMessageQueue()
		}

	case <-time.After(time.Minute):
		// Долго не было запросов - сбросим очередь
		log.Printf("[INFO] %v - Received no events for '%v'. Run l.SaveMessageQueue()", cnt, time.Minute)
		l.SaveMessageQueue()
		//Напечатаем статистику использования pool
		if l.db != nil {
			log.Printf("[INFO] %v", l.db.PrintfByteSliceStat())
		}
	}
}

// SaveMessageQueue - записывает очередь в БД
// =====================================================================
func (l *Listener) SaveMessageQueue() {
	cnt := "json (l *Listener) SaveMessageQueue" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	if len(l.orig) == 0 {
		return
	}

	err := l.db.Save(l.bucket, l.slice)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - l.db.SaveBatch(l.slice)", cnt)
		log.Printf(errM)
	}

	// если успешная запись, то проставим статус для кэша
	for i := range l.orig {
		l.orig[i].flag = JSON_IN_BLOB_DB // устанавливаем флаг, что записали в БД
		l.orig[i].buf = nil
	}

	// если успешная запись, то обнулим указатели
	for i := range l.slice {
		l.slice[i].Buf = nil
		l.slice[i] = nil
	}

	l.slice = l.slice[:0] // сбрасываем размер
	l.orig = l.orig[:0]   // сбрасываем размер
}

// sendSaveMessage - отправка сообщения в канал для записи
// =====================================================================
func (l *Listener) sendSaveMessage(m *cacheStr) {
	cnt := "json (l *Listener) sendSaveMessage" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	if m == nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - m *cacheStr IS NULL", cnt)
		log.Printf(errM)
	}

	// Если в канале есть место, то отправляем, иначе выбрасываем структуру
	select {
	case l.ch <- m:
		m.flag = JSON_IN_SAVE_CHANEL // устанавливаем флаг, что отдали в канал
	default:
		m.flag = JSON_NOT_ENOUGH_SPACE_IN_CHANEL // устанавливаем флаг, что не хватило место в канале
		m.buf = nil                              // выбрасываем буфер
		//log.Printf("[WARN] %v - not enough space in chanel len(l.ch) = '%v'", cnt, len(l.ch))
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return
}
