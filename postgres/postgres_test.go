package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/logutils"

	_ "github.com/lib/pq"
	"github.com/romapres2010/resttest/model"
)

/*
	=====================================================================
	Открытие тестовой БД вынести в отдельный пакет (кросс зависимости)
	=====================================================================
*/

var testDB *PgDb // тестовая БД
const ConfigFileTest = "bench_test_database.cfg"

func InitDbTest(conf string) *PgDb {
	// Фильтр логирования
	var logFilter = &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "INFO", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stderr,
	}
	// Переопределяем глобальный логер на кастомный
	log.SetOutput(logFilter)

	cfgDb := Config{ConfigFile: "", ConnectString: ""}
	cfgDb.ConfigFile = ConfigFileTest
	// Инициализируем БД
	db, err := InitDb(cfgDb)
	if err != nil {
		log.Printf("[ERROR] BenchmarkPgDb_GetDept %v", fmt.Sprintf("%+v", err))
		return nil
	}

	return db
}

func BenchmarkPgDb_GetDept(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping test in short mode.")
	}

	// Инициализируем БД для теста
	if testDB == nil {
		testDB = InitDbTest(ConfigFileTest)
	}

	// считаем все PK для dept
	deptsPK, err := testDB.GetDeptsPK()
	if err != nil {
		b.Errorf("\n PgDb.GetDept() - db.GetDeptsPK() error = %v", err)
		return
	}
	deptsPKlen := len(deptsPK)

	// готовим структуру для тестирования
	tests := struct {
		name string
		p    *PgDb
		args []*model.DeptPK
		//		want    *model.Dept
		//		wantErr bool
	}{
		name: "test",
		p:    testDB,
		args: deptsPK,
	}

	b.ResetTimer()
	b.SetBytes(200000)
	//b.SetParallelism(2)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // The loop body is executed b.N times total across all goroutines
			b.StopTimer()

			// случайным образом выбираем PK
			par := tests.args[rand.Intn(deptsPKlen)].Deptno

			b.StartTimer()

			_, err := tests.p.GetDept(par)
			if err != nil {
				b.Errorf("\n PgDb.GetDept() - error GetDept(d.Deptno), param: %v/n %v", par, fmt.Sprintf("%+v", err))
				return
			}
		}
	})
}

func BenchmarkPgDb_GetDeptEmps(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping test in short mode.")
	}

	// Инициализируем БД для теста
	if testDB == nil {
		testDB = InitDbTest(ConfigFileTest)
	}

	// считаем все PK для dept
	deptsPK, err := testDB.GetDeptsPK()
	if err != nil {
		b.Errorf("\n PgDb.GetDept() - db.GetDeptsPK() error = %v", err)
		return
	}
	deptsPKlen := len(deptsPK)

	// готовим структуру для тестирования
	tests := struct {
		name string
		p    *PgDb
		args []*model.DeptPK
		//		want    *model.Dept
		//		wantErr bool
	}{
		name: "test",
		p:    testDB,
		args: deptsPK,
	}

	b.ResetTimer()
	b.SetBytes(200000)
	//b.SetParallelism(2)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // The loop body is executed b.N times total across all goroutines
			b.StopTimer()

			// случайным образом выбираем PK
			par := tests.args[rand.Intn(deptsPKlen)].Deptno

			b.StartTimer()

			_, err := tests.p.GetDeptEmps(par)
			if err != nil {
				b.Errorf("\n PgDb.GetDept() - error GetDept(d.Deptno), param: %v/n %v", par, fmt.Sprintf("%+v", err))
				return
			}
		}
	})
}

/*
func BenchmarkPgDb_GetDepts(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping test in short mode.")
	}

	// Инициализируем БД для теста
	if testDB == nil {
		testDB = InitDbTest(ConfigFileTest)
	}

	b.ResetTimer()
	b.SetBytes(200000)
	//b.SetParallelism(2)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // The loop body is executed b.N times total across all goroutines
			_, err := testDB.GetDeptsPK()
			if err != nil {
				b.Errorf("\n PgDb.GetDepts() - error GetDepts() : /n %v", fmt.Sprintf("%+v", err))
				return
			}
		}
	})
}
*/
func BenchmarkPgDb_UpdateDept(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping test in short mode.")
	}

	// Инициализируем БД для теста
	if testDB == nil {
		testDB = InitDbTest(ConfigFileTest)
	}

	// считаем все PK для dept
	deptsPK, err := testDB.GetDeptsPK()
	if err != nil {
		b.Errorf("\n PgDb.GetDept() - db.GetDeptsPK() error = %v", err)
		return
	}
	deptsPKlen := len(deptsPK)

	// готовим структуру для тестирования
	tests := struct {
		name string
		p    *PgDb
		args []*model.DeptPK
		//		want    *model.Dept
		//		wantErr bool
	}{
		name: "test",
		p:    testDB,
		args: deptsPK,
	}

	b.ResetTimer()
	b.SetBytes(200000)
	//b.SetParallelism(2)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // The loop body is executed b.N times total across all goroutines
			b.StopTimer()

			// случайным образом выбираем PK
			par := tests.args[rand.Intn(deptsPKlen)].Deptno

			// считываем случайный dept
			dept, err := tests.p.GetDept(par)
			if err != nil {
				b.Errorf("\n PgDb.GetDept() - error GetDept(d.Deptno), param: %v/n %v", par, fmt.Sprintf("%+v", err))
				return
			}

			b.StartTimer()

			_, err = tests.p.UpdateDept(dept)
			if err != nil {
				b.Errorf("\n PgDb.UpdateDept() - error UpdateDept(par), param: %v/n %v", dept, fmt.Sprintf("%+v", err))
				return
			}
		}
	})
}

func BenchmarkPgDb_CreateDept(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping test in short mode.")
	}

	// Инициализируем БД для теста
	if testDB == nil {
		testDB = InitDbTest(ConfigFileTest)
	}

	//готовим фиктивные данные два вставки
	var emp = model.Emp{
		Empno: 0,
		Ename: sql.NullString{
			String: "Ename - insertBenchmark",
			Valid:  true,
		},
		Job: sql.NullString{
			String: "Job - insertBenchmark",
			Valid:  true,
		},
		Mgr: sql.NullInt64{
			Int64: 0,
			Valid: false,
		},
		Hiredate: sql.NullString{
			String: "1981-06-09T00:00:00Z",
			Valid:  true,
		},
		Sal: sql.NullInt64{
			Int64: 99,
			Valid: true,
		},
		Comm: sql.NullInt64{
			Int64: 99,
			Valid: true,
		},
		Deptno: sql.NullInt64{
			Int64: 0,
			Valid: true,
		},
	}

	var dept = model.Dept{
		Deptno: 0,
		Dname:  "Dname - insertBenchmark",
		Loc: sql.NullString{
			String: "Loc - insertBenchmark",
			Valid:  true,
		},
		Emps: nil,
	}

	// готовим структуру для тестирования
	tests := struct {
		name string
		p    *PgDb
		args []*model.Dept
		//		want    *model.Dept
		//		wantErr bool
	}{
		name: "test",
		p:    testDB,
		args: nil,
	}

	b.ResetTimer()
	b.SetBytes(200000)
	//b.SetParallelism(2)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() { // The loop body is executed b.N times total across all goroutines
			b.StopTimer()

			var depnoPK model.DeptPK // локальная копия для хранения последовательности
			var empnoPK model.EmpPK  // локальная копия для хранения последовательности
			var deptNew = dept       // локальная копия структуры
			deptNew.Emps = make([]*model.Emp, 10)

			// считываем последовательность для dept PK
			err := tests.p.Db.Get(&depnoPK, "select nextval('dept_deptno_seq') deptno")
			if err != nil {
				b.Errorf("\n PgDb.CreateDept() - error tests.p.Db.Get(&depnoPK, select nextval('dept_deptno_seq')) %v", fmt.Sprintf("%+v", err))
				return
			}

			// Подменяем PK в тестовых данных dept
			deptNew.Deptno = depnoPK.Deptno

			// Подменяем PK в тестовых данных emp
			for i, _ := range deptNew.Emps {
				// создаем новую переменную
				var e = emp
				// записываем номер dept
				e.Deptno.Int64 = int64(depnoPK.Deptno)
				e.Deptno.Valid = true
				// считываем последовательность для emp PK
				err := tests.p.Db.Get(&empnoPK, "select nextval('dept_deptno_seq') empno")
				if err != nil {
					b.Errorf("\n PgDb.CreateDept() - error tests.p.Db.Get(&depnoPK, select nextval('dept_deptno_seq')) %v", fmt.Sprintf("%+v", err))
					return
				}
				// записываем номер emp
				e.Empno = empnoPK.Empno
				deptNew.Emps[i] = &e
			}

			b.StartTimer()

			//b.Logf("&deptNew %+v/n", deptNew)
			_, err = tests.p.CreateDept(&deptNew)

			if err != nil {
				b.Errorf("\n PgDb.CreateDept() - error tests.p.CreateDept(&deptNew), param: %v/n %v", deptNew, fmt.Sprintf("%+v", err))
				return
			}
		}
	})

	/*
		=====================================================================
		Добавить удаление тестовых данных
		=====================================================================
	*/
}
