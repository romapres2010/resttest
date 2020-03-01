package bolt

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/FogCreek/mini"
	bolt "go.etcd.io/bbolt"

	"github.com/pkg/errors"
)

// Config конфигурационные настройки
type Config struct {
	ConfigFile string // имя файла-конфигурации для подключения
	DbFilePath string // расположение файла с Бд
	DbTimeout  int    // время ожидания открытия БД в секундах
	//	DbReadOnly        bool   // Open database in read-only mode
	//	DbNoFreelistSync  bool   // Do not sync freelist to disk
	//	DbInitialMmapSize int    // InitialMmapSize is the initial mmap size of the database in bytes
}

// Db represents a bbolt connection
type Db struct {
	db *bolt.DB
}

// SaveStr - Структура для записи в Bolt
type SaveStr struct {
	Key int    // ключ для БД
	Buf []byte // JSON
}

// SaveSlice - Слайс для пакетной записи
type SaveSlice []*SaveStr

// InitDb - инициализирует подключение к БД
func InitDb(cfg Config) (*Db, error) {
	cnt := "bolt InitDb" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START, param: '%+v'", cnt, cfg)

	var err error
	var boltDb Db

	// считаем конфигурацию Bolt DB из внешнего файла
	{
		if cfg.ConfigFile == "" {
			return nil, errors.New(fmt.Sprintf("[ERROR] %v - Bolt DB config file name is null", cnt))
		}
		log.Printf("[INFO] %v - Loading Bolt DB config from file '%s'", cnt, cfg.ConfigFile)
		// Считать информацию о файле или каталоге
		_, err := os.Stat(cfg.ConfigFile)
		// Если файл не существует
		if os.IsNotExist(err) {
			errM := fmt.Sprintf("[ERROR] %v - Bolt DB config file '%s' does not exist", cnt, cfg.ConfigFile)
			log.Printf(errM)
			return nil, errors.New(errM)
		}

		// Считать конфигурацию из файла
		pgcfg, err := mini.LoadConfiguration(cfg.ConfigFile)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - mini.LoadConfiguration(cfg.ConfigFile) \n %+v", cnt, err)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}

		// Считать  параметры Бд
		cfg.DbFilePath = pgcfg.String("DbFilePath", "")
		if cfg.DbFilePath == "" {
			errM := fmt.Sprintf("[ERROR] %v - parameter DbFilePath IS NULL", cnt)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
		log.Printf("[INFO] %v - parameter DbFilePath '%+v'", cnt, cfg.DbFilePath)

		val := pgcfg.String("DbTimeout", "10")
		cfg.DbTimeout, err = strconv.Atoi(val)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - incorrect number - parameter DbTimeout = '%v'", cnt, val)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
		log.Printf("[INFO] %v - parameter DbTimeout '%+v'", cnt, cfg.DbTimeout)

	}

	// Формируем опции Botl DB
	opt := &bolt.Options{
		Timeout: time.Duration(cfg.DbTimeout * int(time.Second)),
	}

	// Открываем Botl DB
	log.Printf("[INFO] %v - Opening Bolt DB database '%s'", cnt, cfg.DbFilePath)
	boltDb.db, err = bolt.Open(cfg.DbFilePath, 0666, opt)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - bolt.Open(cfg.dbFilePath, 0666, opt)", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	{ // Создаем Buckets "DeptBucket"
		log.Printf("[INFO] %v - create buckets 'DeptBucket'", cnt)
		err := boltDb.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("DeptBucket"))
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.CreateBucketIfNotExists([]byte('DeptBucket'))", cnt)
				log.Printf(errM)
				return errors.Wrap(err, errM)
			}
			return nil
		})
		if err != nil {
			defer boltDb.db.Close()
			errM := fmt.Sprintf("[ERROR] %v - ERROR - boltDb.db.Update(func(tx *bolt.Tx)", cnt)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
		log.Printf("[INFO] %v - buckets 'DeptBucket' was created", cnt)
	}
	{ // Создаем Buckets "EmpBucket"
		log.Printf("[INFO] %v - create buckets 'EmpBucket'", cnt)
		err := boltDb.db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("EmpBucket"))
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.CreateBucketIfNotExists([]byte('EmpBucket'))", cnt)
				log.Printf(errM)
				return errors.Wrap(err, errM)
			}
			return nil
		})
		if err != nil {
			defer boltDb.db.Close()
			errM := fmt.Sprintf("[ERROR] %v - ERROR - boltDb.db.Update(func(tx *bolt.Tx)", cnt)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
		log.Printf("[INFO] %v - buckets 'EmpBucket' was created", cnt)
	}

	log.Printf("[INFO] %v - Bolt DB database is opened", cnt)
	log.Printf("[INFO] %v - SUCCESS", cnt)
	return &boltDb, nil
}

// Close the Bolt DB. Safe to nil pointer
func (d *Db) Close() error {
	cnt := "bolt (d *BoltDb) Close()" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - d.db.Close()", cnt)
			log.Printf(errM)
			return errors.Wrap(err, errM)
		}
	} else {
		log.Printf("[INFO] %v - d.db IS NULL", cnt)
	}

	log.Printf("[INFO] %v - SUCCESS", cnt)
	return nil
}

// PutByteSlice return byte buf to cache
// =====================================================================
func (d *Db) PutByteSlice(p []byte) {
	if p != nil {
		putByteSlice(p)
	}
}

// GetByteSlice allocates a new []byte
// =====================================================================
func (d *Db) GetByteSlice() []byte {
	return getByteSlice()
}

// PrintfByteSliceStat return statistics of using byte pool
// =====================================================================
func (d *Db) PrintfByteSliceStat() string {
	return fmt.Sprintf("PrintfByteSliceStat() - countNewByteSlice = '%v', countGetByteSlice = '%v', countPutByteSlice = '%v'", stat.countNewByteSlice, stat.countGetByteSlice, stat.countPutByteSlice)
}

//Save save all json from slice into DB
// =====================================================================
func (d *Db) Save(bucket string, bufs SaveSlice) error {
	cnt := "bolt  (d *Db) Save" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, len(bufs) = '%v'", cnt, len(bufs))
	log.Printf("[INFO] %v - START, len(bufs) = '%v'", cnt, len(bufs))
	//log.Printf("[INFO] %v", d.PrintfByteSliceStat())

	if bucket == "" {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - Bolt DB Bucket Name IS NULL", cnt)
		log.Printf(errM)
		return errors.New(errM)
	}

	if bufs == nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - buf IS NULL", cnt)
		log.Printf(errM)
		return errors.New(errM)
	}

	err := d.db.Update(func(tx *bolt.Tx) error {
		// возьмем Bucket
		b := tx.Bucket([]byte(bucket))

		for _, v := range bufs {
			// Запишем в БД
			err := b.Put(itob(v.Key), v.Buf)
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - b.Put(itob(v.Key), v.Buf), args = '%+v'", cnt, v.Key)
				log.Printf(errM)
				return errors.Wrap(err, errM)
			}
		}

		//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
		return nil
	})

	return err
}

//Get get json from DB
// =====================================================================
func (d *Db) Get(bucket string, key int) ([]byte, error) {
	cnt := "bolt  (d *Db) Get" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, args = '%+v', '%v'", cnt, bucket, key)

	if bucket == "" {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - Bolt DB Bucket Name IS NULL", cnt)
		log.Printf(errM)
		return nil, errors.New(errM)
	}

	var retbuf []byte // переменная для возврата

	// считываем данные
	err := d.db.View(func(tx *bolt.Tx) error {

		// возьмем Bucket
		b := tx.Bucket([]byte(bucket))

		// Считаем из БД
		buf := b.Get(itob(key))

		if buf != nil {
			// копируем данные во внешнюю переменную
			//retbuf = make([]byte, len(buf))
			retbuf = d.GetByteSlice() // cap получаемого буффера может быть произвольным !!!
			retbuf = retbuf[:0]       // сбросим указатель на начало
			// если размер мал, то выбросим старый буффер и создадим новый
			if cap(retbuf) < len(buf) {
				retbuf = make([]byte, 0, len(buf))
			}
			//copy(retbuf, buf) // Будет скопировано не более размера целевого среза
			retbuf = append(retbuf, buf...)
			//mylog.PrintfDebug("[DEBUG] %v - SUCCESS get buf from Bolt DB, key='%v', '%p'", cnt, key, retbuf)
		} else {
			//mylog.PrintfDebug("[DEBUG] %v - buf readen from Bolt DB is NULL, key='%v'", cnt, key)
		}

		buf = nil // проверка как работает сборщик мусора

		return nil
	})

	return retbuf, err
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
