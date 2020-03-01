package daemon

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/logutils"

	"github.com/gorilla/mux"

	"github.com/pkg/errors"
	"github.com/romapres2010/resttest/bolt"
	"github.com/romapres2010/resttest/cache"
	myhttp "github.com/romapres2010/resttest/http"
	"github.com/romapres2010/resttest/json"
	//mylog "github.com/romapres2010/resttest/log"
	"github.com/romapres2010/resttest/model"
	"github.com/romapres2010/resttest/postgres"
)

// Config - представляет конфигурационные настройки
type Config struct {
	LogFilter  *logutils.LevelFilter // ссылка на кастомные loger - пока не используется
	ListenSpec string                // строка HTTP листенера
	UseCache   string                // использовать кэш?
	PqDbCfg    postgres.Config       // конфигурация PostgreSQL DB
	UseBoltDb  bool                  // признак что используем Bolt DB
	BoltDbCfg  bolt.Config           // конфигурация Bolt DB
	JSONCfg    json.Config           // конфигурация JSON Cahce
}

// Run - cтартует демон и ожидает сигнала выхода
func Run(cfg *Config) error {
	cnt := "daemon Run" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START, param: '%+v'", cnt, cfg)

	// Основной сервис
	var jsonService *json.Service

	// Инициируем PostgreSQL БД
	pqDb, err := postgres.InitDb(cfg.PqDbCfg)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - postgres.InitDb(cfg.PqDbCfg), args = '%+v'", cnt, cfg.PqDbCfg)
		log.Printf(errM)
		return errors.WithMessage(err, errM)
	}

	// Инициируем Bolt DB
	var boltDB *bolt.Db
	if cfg.UseBoltDb {
		boltDB, err = bolt.InitDb(cfg.BoltDbCfg)
		defer func() {
			if boltDB != nil {
				boltDB.Close()
			}
		}()
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - bolt.InitDb(cfg.BoltDbCfg), args = '%+v'", cnt, cfg.BoltDbCfg)
			log.Printf(errM)
			return errors.WithMessage(err, errM)
		}
	} else {
		log.Printf("[INFO] %v - DO NOT use Bolt DB", cnt)
		boltDB = nil
	}

	// Инициируем и наполняем кэш
	if cfg.UseCache == "yes" || cfg.UseCache == "populate" {
		log.Printf("[INFO] %v - create cache", cnt)

		// Инициируем кэш
		cacheService := cache.New(pqDb, pqDb)

		// Создаем листенер и устанавливаем его для БД и кэша
		listener := model.NewListener(cacheService) // устанавливаем кэш в качестве обработчика
		pqDb.SetModelListener(listener)             // БД использует канал листенера на передачу
		cacheService.SetModelListener(listener)     // Кэш использует канал листенера на прием

		// преднаполним кэш - для целей тестирования
		if cfg.UseCache == "populate" {
			log.Printf("[INFO] %v - populate cache", cnt)
			go cacheService.PopulateCache()
		} else {
			log.Printf("[INFO] %v - DO NOT populate cache", cnt)
		}

		jsonService = json.New(cfg.JSONCfg, cacheService, cacheService, boltDB)
		if cfg.JSONCfg.IsPopulateCache {
			go jsonService.PopulateDeptCache()
		}
	} else {
		jsonService = json.New(cfg.JSONCfg, pqDb, pqDb, boltDB)
		if cfg.JSONCfg.IsPopulateCache {
			go jsonService.PopulateDeptCache()
		}
		log.Printf("[INFO] %v - DO NOT use cache", cnt)
	}

	// Создаем свой HTTP обработчик
	myHandler := myhttp.New(jsonService)

	// стартуем сервер, ожидаем прерывания
	err = StartHTTPHandlerWait(cfg, myHandler, cfg.ListenSpec)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - StartHTTPHandlerWait", cnt)
		log.Printf(errM)
		return errors.WithMessage(err, errM)
	}

	log.Printf("[INFO] %v - SUCCESS", cnt)

	return nil
}

// StartHTTPHandlerWait - устанавливает обработчиков и стартует в фоне сервер
func StartHTTPHandlerWait(cfg *Config, myHandler *myhttp.Handler, listenSpec string) error {
	cnt := "daemon StartHTTPHandlerWait" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START, param: '%+v'", cnt, cfg)

	// Определяем  листенер
	listener, err := net.Listen("tcp", listenSpec)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - net.Listen('tcp', listenSpec) param: '%v'", cnt, listenSpec)
		log.Printf(errM)
		return errors.Wrap(err, errM)
	}
	log.Printf("[INFO] %v - Starting, HTTP on: %s", cnt, listenSpec)

	// Настраиваем роутер
	router := mux.NewRouter()
	router.HandleFunc("/depts/{id:[0-9]+}", RecoverWrap(http.HandlerFunc(setDeptGetHandler(myHandler)))).Methods("GET")
	router.HandleFunc("/deptemps/{id:[0-9]+}", RecoverWrap(http.HandlerFunc(setDeptEmpsGetHandler(myHandler)))).Methods("GET")
	router.HandleFunc("/depts", RecoverWrap(http.HandlerFunc(setDeptAddHandler(myHandler)))).Methods("POST")
	router.HandleFunc("/depts", RecoverWrap(http.HandlerFunc(setDeptUpdateHandler(myHandler)))).Methods("PUT")
	router.HandleFunc("/depts/emps/{id:[0-9]+}", RecoverWrap(http.HandlerFunc(setEmpGetHandler(myHandler)))).Methods("GET")
	router.HandleFunc("/depts/emps", RecoverWrap(http.HandlerFunc(setEmpAddHandler(myHandler)))).Methods("POST")
	// Для целей нагрузочного тестирования
	{
		router.HandleFunc("/depts/random", RecoverWrap(http.HandlerFunc(setRandomDeptGetHandler(myHandler)))).Methods("GET")
		router.HandleFunc("/depts/random", RecoverWrap(http.HandlerFunc(setRandomDeptUpdateHandler(myHandler)))).Methods("PUT")
	}

	// Регистрация pprof-обработчиков для gorilla/mux
	{
		pprofrouter := router.PathPrefix("/debug/pprof").Subrouter()
		pprofrouter.HandleFunc("/", pprof.Index)
		pprofrouter.HandleFunc("/cmdline", pprof.Cmdline)
		pprofrouter.HandleFunc("/symbol", pprof.Symbol)
		pprofrouter.HandleFunc("/trace", pprof.Trace)

		profile := pprofrouter.PathPrefix("/profile").Subrouter()
		profile.HandleFunc("", pprof.Profile)
		profile.Handle("/goroutine", pprof.Handler("goroutine"))
		profile.Handle("/threadcreate", pprof.Handler("threadcreate"))
		profile.Handle("/heap", pprof.Handler("heap"))
		profile.Handle("/block", pprof.Handler("block"))
		profile.Handle("/mutex", pprof.Handler("mutex"))
	}

	// Конфигурация сервера
	server := &http.Server{
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	// устанавливаем роутер
	http.Handle("/", router)

	// Запускаем сервер
	go server.Serve(listener)

	log.Printf("[INFO] %v - Daemon is running. For exit <CTRL-c>", cnt)

	// Ожидание сигнала прерывания
	waitForSignal()

	// Попытка корректного завершения после 5 секунд
	log.Printf("[INFO] %v - Waiting for shutdown 5 sec", cnt)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Printf("[INFO] %v - Shutdown completed", cnt)

	return nil
}

func waitForSignal() {
	cnt := "daemon waitForSignal" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	s := <-ch
	log.Printf("[INFO] %v - Got signal: %v, exiting", cnt, s)
}

// RecoverWrap cover handler functions with panic recoverer
func RecoverWrap(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cnt := "daemon RecoverWrap - HandlerFunc" // Имя текущего метода для логирования
		//mylog.PrintfDebug("[DEBUG] ==========================================================================")
		//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

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
				//err = errors.Wrap(err, "[ERROR] %v - ERROR", cnt)
				//sendMeMail(err)
				log.Printf("[ERROR] %v - ERROR - \n %+v", cnt, err.Error())

				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)

		//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
		//mylog.PrintfDebug("[DEBUG] ==========================================================================")
	})
}

func setDeptGetHandler(h *myhttp.Handler) func(http.ResponseWriter, *http.Request) {
	return h.DeptGetHandler
}

func setDeptEmpsGetHandler(h *myhttp.Handler) func(http.ResponseWriter, *http.Request) {
	return h.DeptEmpsGetHandler
}

func setDeptAddHandler(h *myhttp.Handler) func(http.ResponseWriter, *http.Request) {
	return h.DeptAddHandler
}

func setDeptUpdateHandler(h *myhttp.Handler) func(http.ResponseWriter, *http.Request) {
	return h.DeptUpdateHandler
}

func setEmpGetHandler(h *myhttp.Handler) func(http.ResponseWriter, *http.Request) {
	return h.EmpGetHandler
}

func setEmpAddHandler(h *myhttp.Handler) func(http.ResponseWriter, *http.Request) {
	return h.EmpAddHandler
}

func setRandomDeptGetHandler(h *myhttp.Handler) func(http.ResponseWriter, *http.Request) {
	return h.RandomDeptGetHandler
}

func setRandomDeptUpdateHandler(h *myhttp.Handler) func(http.ResponseWriter, *http.Request) {
	return h.RandomDeptUpdateHandler
}
