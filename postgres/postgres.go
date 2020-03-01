package postgres

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/FogCreek/mini"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/romapres2010/resttest/model"
)

// Config конфигурационные настройки
type Config struct {
	ConfigFile                          string        // имя файла-конфигурации для подключения
	ConnectString                       string        // строка подключения к БД
	UseNotification                     bool          // использовать нотификацию из БД
	ChanelPrefix                        string        // префикс каналов листенера для работы с Notify PostgreSQL
	ListenerMinReconnectInterval        time.Duration // минимальное время попытки восстановления подключения для работы с Notify PostgreSQL
	ListenerMaxReconnectInterval        time.Duration // максимальное время попытки восстановления подключения для работы с Notify PostgreSQL
	ListenerWaitForNotificationInterval time.Duration // время ожидания события при работе с Notify PostgreSQL
	PqConnMaxLifetime                   int           // время жизни подключения в милисекундах
	PqMaxOpenConns                      int           // максимальное количество открытых подключений
	PqMaxIdleConns                      int           // максимальное количество простаивающих подключений

}

// текст SQL запросов
const (
	// dept
	sqlSelectDeptExistsText string = "SELECT deptno FROM dept WHERE deptno = $1"
	sqlSelectDeptText       string = "SELECT deptno, dname, loc FROM dept WHERE deptno = $1"
	sqlSelectDeptsText      string = "SELECT deptno, dname, loc FROM dept"
	sqlSelectDeptsPKText    string = "SELECT deptno FROM dept"
	sqlInsertDeptText       string = "INSERT INTO dept (deptno, dname, loc) VALUES (:deptno, :dname, :loc)"
	sqlUpdateDeptText       string = "UPDATE dept SET dname = :dname, loc = :loc WHERE deptno = :deptno"

	// Emp
	sqlSelectEmpExistsText    string = "SELECT empno FROM emp WHERE empno = $1"
	sqlSelectEmpText          string = "SELECT empno, ename, job, mgr, hiredate, sal, comm, deptno FROM emp WHERE empno = $1"
	sqlSelectEmpsByDeptText   string = "SELECT empno, ename, job, mgr, hiredate, sal, comm, deptno FROM emp WHERE deptno = $1"
	sqlSelectEmpsPKByDeptText string = "SELECT empno FROM emp WHERE deptno = $1"
	sqlInsertEmpText          string = "INSERT INTO emp (empno, ename, job, mgr, hiredate, sal, comm, deptno) VALUES (:empno, :ename, :job, :mgr, :hiredate, :sal, :comm, :deptno)"
	sqlUpdateEmpText          string = "UPDATE emp SET empno = :empno, ename = :ename, job = :job, mgr = :mgr, hiredate = :hiredate, sal = :sal, comm = :comm, deptno = :deptno WHERE empno = :empno"
)

// PgDb represents a PostgreSQL connection
type PgDb struct {
	Db                                    *sqlx.DB
	pqListener                            *pq.Listener  // листенер для работы с Notify PostgreSQL
	pqListenerWaitForNotificationInterval time.Duration // время ожидания события при работе с Notify PostgreSQL
	pqChanelDeptName                      string
	pqChanelEmpName                       string
	modelListener                         *model.Listener // обработчик уведомлений

	// Dept
	sqlSelectDeptExists *sqlx.Stmt
	sqlSelectDept       *sqlx.Stmt
	sqlSelectDepts      *sqlx.Stmt
	sqlSelectDeptsPK    *sqlx.Stmt
	//sqlInsertDept       *sqlx.NamedStmt

	// Emp
	sqlSelectEmpExists    *sqlx.Stmt
	sqlSelectEmp          *sqlx.Stmt
	sqlSelectEmpsByDept   *sqlx.Stmt
	sqlSelectEmpsPKByDept *sqlx.Stmt
}

// InitDb - инициализирует подключение к БД
func InitDb(cfg Config) (*PgDb, error) {
	cnt := "postgers InitDb" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START, param: '%+v'", cnt, cfg)

	var err error
	var PgDb PgDb

	// считаем конфигурацию PostgreSQL из внешнего файла
	{
		if cfg.ConfigFile == "" {
			return nil, errors.New(fmt.Sprintf("[ERROR] %v - PostgreSQL config file name is null", cnt))
		}
		log.Printf("[INFO] %v - Loading PostgreSQL config from file '%s'", cnt, cfg.ConfigFile)
		// Считать информацию о файле или каталоге
		_, err := os.Stat(cfg.ConfigFile)
		// Если файл не существует
		if os.IsNotExist(err) {
			errM := fmt.Sprintf("[ERROR] %v - PostgreSQL config file '%s' does not exist", cnt, cfg.ConfigFile)
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
		// Сформировать строку с параметрами
		cfg.ConnectString = fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s user=%s password=%s ",
			pgcfg.String("host", "127.0.0.1"),
			pgcfg.String("port", "5432"),
			pgcfg.String("dbname", ""),
			pgcfg.String("sslmode", ""),
			pgcfg.String("user", ""),
			pgcfg.String("pass", ""),
		)
		// Считать дополнительные параметры Бд
		val := pgcfg.String("PqMaxOpenConns", "5")
		cfg.PqMaxOpenConns, err = strconv.Atoi(val)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - incorrect number - parameter PqMaxOpenConns = '%v'", cnt, val)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
		log.Printf("[INFO] %v - cfg.PqMaxOpenConns = '%v'", cnt, cfg.PqMaxOpenConns)

		val = pgcfg.String("PqMaxIdleConns", "2")
		cfg.PqMaxIdleConns, err = strconv.Atoi(val)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - incorrect number - parameter PqMaxIdleConns = '%v'", cnt, val)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
		log.Printf("[INFO] %v - cfg.PqMaxIdleConns = '%v'", cnt, cfg.PqMaxIdleConns)

		val = pgcfg.String("PqConnMaxLifetime", "5")
		cfg.PqConnMaxLifetime, err = strconv.Atoi(val)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - incorrect number - parameter PqConnMaxLifetime = '%v'", cnt, val)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
		log.Printf("[INFO] %v - cfg.PqConnMaxLifetime = '%v' millisecond", cnt, cfg.PqConnMaxLifetime)

	}

	// открываем соединение с БД
	log.Printf("[INFO] %v - Opening PostgreSQL database %s", cnt, cfg.ConnectString)
	PgDb.Db, err = sqlx.Connect("postgres", cfg.ConnectString)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - sqlx.Connect('postgres', '%v')", cnt, cfg.ConnectString)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}
	log.Printf("[INFO] %v - PostgreSQL database is opened", cnt)

	// Устанавливаем параметры подключений БД
	{
		PgDb.Db.SetMaxOpenConns(cfg.PqMaxOpenConns)
		PgDb.Db.SetMaxIdleConns(cfg.PqMaxIdleConns)
		PgDb.Db.SetConnMaxLifetime(time.Duration(cfg.PqConnMaxLifetime * int(time.Millisecond)))
	}

	// подготовим SQL команды (DML команды не готовим)
	log.Printf("[INFO] %v - Prepare SQL statements", cnt)
	err = PgDb.prepareSQLStatements()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - PgDb.prepareSQLStatements() \n %+v", cnt, err)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// Создаем листенер для ожидания уведомлений и запускаем бесконечный цикл
	if cfg.UseNotification {
		log.Printf("[INFO] %v - Create listener for notification", cnt)
		err = PgDb.createListenerForNotification(cfg)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - PgDb.createListenerForNotification(cfg) \n %+v", cnt, err)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
	} else {
		log.Printf("[INFO] %v - DO NOT use notification", cnt)
	}

	log.Printf("[INFO] %v - SUCCESS", cnt)

	return &PgDb, nil
}

// prepareSqlStatements - парсит SQL запросы после открытия БД
func (p *PgDb) prepareSQLStatements() (err error) {
	cnt := "postgers prepareSQLStatements" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	{ // Dept
		p.sqlSelectDeptExists, err = p.Db.Preparex(sqlSelectDeptExistsText)
		if err != nil {
			return errors.Wrap(err, "p.Db.Preparex(sqlSelectDeptExistsText)")
		}
		p.sqlSelectDept, err = p.Db.Preparex(sqlSelectDeptText)
		if err != nil {
			return errors.Wrap(err, "p.Db.Preparex(sqlSelectDeptText)")
		}
		p.sqlSelectDeptsPK, err = p.Db.Preparex(sqlSelectDeptsPKText)
		if err != nil {
			return errors.Wrap(err, "p.Db.Preparex(sqlSelectDeptsPKText)")
		}
		p.sqlSelectDepts, err = p.Db.Preparex(sqlSelectDeptsText)
		if err != nil {
			return errors.Wrap(err, "p.Db.Preparex(sqlSelectDeptsText)")
		}
	}
	{ // Emp
		p.sqlSelectEmpExists, err = p.Db.Preparex(sqlSelectEmpExistsText)
		if err != nil {
			return errors.Wrap(err, "p.Db.Preparex(sqlSelectEmpExistsText)")
		}
		p.sqlSelectEmp, err = p.Db.Preparex(sqlSelectEmpText)
		if err != nil {
			return errors.Wrap(err, "p.Db.Preparex(sqlSelectEmpText)")
		}
		p.sqlSelectEmpsByDept, err = p.Db.Preparex(sqlSelectEmpsByDeptText)
		if err != nil {
			return errors.Wrap(err, "p.Db.Preparex(sqlSelectEmpsByDeptText)")
		}
		p.sqlSelectEmpsPKByDept, err = p.Db.Preparex(sqlSelectEmpsPKByDeptText)
		if err != nil {
			return errors.Wrap(err, "p.Db.Preparex(sqlSelectEmpsPKByDeptText)")
		}
	}

	log.Printf("[INFO] %v - SUCCESS", cnt)

	return nil
}
