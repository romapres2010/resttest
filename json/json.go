package json

import (
	"github.com/romapres2010/resttest/bolt"
	"github.com/romapres2010/resttest/model"
)

// Config - представляет конфигурационные настройки
type Config struct {
	DeptBucketName  string // имя bucket для хранения dept JSON
	SaveChSize      int    // размер канала для формирования очереди записи в Bolt
	SaveQueueSize   int    // размер очереди записи в Bolt
	IsPopulateCache bool   // признак наполнять кэш при старте
	IsPassCache     bool   // пробрассывать кэш
}

// Service wraps a model.Services to provide an Json cache
// =====================================================================
type Service struct {
	empService   model.EmpService
	deptService  model.DeptService
	db           *bolt.Db
	deptCache    Cache
	saveListener *Listener // ссылка на обработчик канала для формирования очереди на запись в Bolt
	deptBucket   string    // имя bucket для хранения dept JSON
	isPassCache  bool      // пробрассывать кэш
}

// New returns a new Service
// =====================================================================
func New(cfg Config, eService model.EmpService, dService model.DeptService, db *bolt.Db) *Service {
	s := &Service{
		empService:  eService,
		deptService: dService,
		deptCache: Cache{
			ch: make(map[int]*cacheStr),
		},
		isPassCache: cfg.IsPassCache,
	}

	// если работаем в режиме с Bolt BD
	if db != nil {
		s.deptBucket = cfg.DeptBucketName
		s.db = db

		// создаем и запускаем листенер для обработки канала формирования очереди на запись в Bolt DB
		s.saveListener = newSaveListener(db, cfg.DeptBucketName, cfg.SaveChSize, cfg.SaveQueueSize)
	}

	return s
}

// PutByteSlice return byte buf to cache in Bolt DB pool
// =====================================================================
func (s *Service) PutByteSlice(p []byte) {
	if s.db != nil && p != nil {
		s.db.PutByteSlice(p)
	}
}

// GetByteSlice allocates a new []byte from  Bolt DB pool
// =====================================================================
func (s *Service) GetByteSlice() []byte {
	if s.db != nil {
		return s.db.GetByteSlice()
	}
	return nil
}
