package cache

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/logutils"

	_ "github.com/lib/pq"

	"github.com/romapres2010/resttest/postgres"
)

/*
	=====================================================================
	Открытие тестовой БД вынести в отдельный пакет (кросс зависимости)
	=====================================================================
*/

type testCacheType struct {
	mx    sync.RWMutex // для независимой блокировки отдельных кешей
	cache *Cache       // тестовый кэш
}

var testCache testCacheType

const ConfigFileTest = "bench_test_database.cfg"

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

func BenchmarkCache_GetDept(b *testing.B) {
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
			testCache.cache = New(testDB, testDB)

			// Наполним cache
			//testCache.cache.PopulateCache()
		}
	}()

	// считаем все PK для dept
	deptsPK, err := testCache.cache.GetDeptsPK()
	if err != nil {
		b.Errorf("\n PgDb.GetDept() - db.GetDeptsPK() error = %v", err)
		return
	}
	deptsPKlen := len(deptsPK)

	b.ResetTimer()
	b.SetBytes(200000)
	//	b.SetParallelism(1)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // The loop body is executed b.N times total across all goroutines
			b.StopTimer()

			// случайным образом выбираем PK
			par := deptsPK[rand.Intn(deptsPKlen)].Deptno

			b.StartTimer()

			_, err := testCache.cache.GetDept(par)
			//model.PutDept(dept, true)

			if err != nil {
				b.Errorf("\n PgDb.GetDept() - error GetDept(d.Deptno), param: %v/n %v", par, fmt.Sprintf("%+v", err))
				return
			}
		}
	})
}

func BenchmarkCache_GetDeptEmps(b *testing.B) {
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
			testCache.cache = New(testDB, testDB)

			// Наполним cache
			//testCache.cache.PopulateCache()
		}
	}()

	// считаем все PK для dept
	deptsPK, err := testCache.cache.GetDeptsPK()
	if err != nil {
		b.Errorf("\n PgDb.GetDept() - db.GetDeptsPK() error = %v", err)
		return
	}
	deptsPKlen := len(deptsPK)

	b.ResetTimer()
	b.SetBytes(200000)
	//b.SetParallelism(4)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // The loop body is executed b.N times total across all goroutines
			b.StopTimer()

			// случайным образом выбираем PK
			par := deptsPK[rand.Intn(deptsPKlen)].Deptno

			b.StartTimer()

			_, err := testCache.cache.GetDeptEmps(par)
			b.StopTimer()

			if err != nil {
				b.Errorf("\n PgDb.GetDept() - error GetDept(d.Deptno), param: %v/n %v", par, fmt.Sprintf("%+v", err))
				return
			}
		}
	})
}
