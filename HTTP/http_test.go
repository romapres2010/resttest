package http

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/logutils"
	"github.com/romapres2010/resttest/bolt"
	"github.com/romapres2010/resttest/cache"
	"github.com/romapres2010/resttest/json"
	"github.com/romapres2010/resttest/postgres"
)

/*
	=====================================================================
	Открытие тестовой БД вынести в отдельный пакет (кросс зависимости)
	=====================================================================
*/

type testCacheType struct {
	mx    sync.RWMutex  // для независимой блокировки отдельных кешей
	cache *cache.Cache  // тестовый кэш
	json  *json.Service //

}

var testCache testCacheType

const ConfigFileTest = "bench_test_database.cfg"
const BoltCfgFile = "test_boltdb.cfg"

var JSONCfg json.Config = json.Config{
	DeptBucketName:  "DeptBucket",
	SaveChSize:      100000,
	SaveQueueSize:   1000,
	IsPopulateCache: true,
	IsPassCache:     true,
}

var BoltCfg bolt.Config = bolt.Config{
	ConfigFile: "test_boltdb.cfg",
	DbFilePath: "",
	DbTimeout:  1000,
}

func InitDbTest(conf string) *postgres.PgDb {
	// Фильтр логирования
	var logFilter = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "INFO", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}
	// Переопределяем глобальный логер на кастомный
	log.SetOutput(logFilter)

	cfgDb := postgres.Config{ConfigFile: "", ConnectString: ""}
	cfgDb.ConfigFile = ConfigFileTest
	// Инициализируем БД
	db, err := postgres.InitDb(cfgDb)
	if err != nil {
		log.Printf("[ERROR] BenchmarkPgDb_GetDept %v", fmt.Sprintf("%+v", err))
		return nil
	}

	return db
}

func BenchmarkHandler_DeptEmpsGetHandler(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping test in short mode.")
	}

	// Инициализируем кэш для теста
	func() {
		testCache.mx.Lock()
		defer testCache.mx.Unlock()
		if testCache.cache == nil {
			// Инициируем БД
			testDB := InitDbTest(ConfigFileTest)

			// Инициируем кэш
			testCache.cache = cache.New(testDB, testDB)

			// Инициируем Bolt DB
			var boltDB *bolt.Db
			if !JSONCfg.IsPassCache {
				boltDB, err := bolt.InitDb(BoltCfg)
				defer func() {
					if boltDB != nil {
						boltDB.Close()
					}
				}()
				if err != nil {
					errM := fmt.Sprintf("[ERROR] - bolt.InitDb(BoltCfg), args = '%+v'", BoltCfg)
					b.Errorf(errM)
				}
			} else {
				log.Printf("[INFO] DO NOT use Bolt DB")
				boltDB = nil
			}

			testCache.json = json.New(JSONCfg, testCache.cache, testCache.cache, boltDB)

			if JSONCfg.IsPopulateCache {
				testCache.json.PopulateDeptCache()
			}

			// Наполним cache
			//testCache.cache.PopulateCache()
		}
	}()

	// считаем все PK для dept для формирования случайной выборки
	deptsPK, err := testCache.cache.GetDeptsPK()
	if err != nil {
		b.Errorf("\n deptsGetHandler() - db.GetDeptsPK() error = %v", err)
		return
	}
	deptsPKlen := len(deptsPK)

	// Сбрасываем таймер
	b.ResetTimer()
	b.SetBytes(200000)

	// Создаем свой HTTP обработчик
	myHandler := New(testCache.json)

	// Создадим обработчик
	handler := func() func(http.ResponseWriter, *http.Request) {
		//return myHandler.DeptEmpsGetHandler
		return myHandler.DeptEmpsGetHandler
	}()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // The loop body is executed b.N times total across all goroutines
			b.StopTimer()

			// случайным образом выбираем PK
			par := deptsPK[rand.Intn(deptsPKlen)].Deptno

			// подготовим параметры для внедрения в запрос через Gorille Mux
			vars := make(map[string]string)
			vars["id"] = strconv.Itoa(par)

			// Подготовим запрос
			r := httptest.NewRequest("GET", "/depts/id", nil)

			// внедрим параметры в запрос через Gorille Mux
			rMux := mux.SetURLVars(r, vars)

			w := httptest.NewRecorder()

			b.StartTimer()

			//Запускаем обработчик
			handler(w, rMux)

			if err != nil {
				b.Errorf("\n PgDb.GetDept() - error GetDept(d.Deptno), param: %v/n %v", par, fmt.Sprintf("%+v", err))
				return
			}
		}
	})

}
