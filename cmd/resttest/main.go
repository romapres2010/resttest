package main

import (
	"log"
	"os"
	"time"

	"github.com/hashicorp/logutils"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"github.com/romapres2010/resttest/daemon"
	mylog "github.com/romapres2010/resttest/log"
)

// входные флаги программы
var flags = []cli.Flag{
	cli.StringFlag{
		Name:        "databaseconf, dbconf",
		Usage:       "Databsee config file name",
		Destination: &dbConfigFileFlag,
	},
	cli.StringFlag{
		Name:        "listenstring, l",
		Usage:       "Listen string in format <host>:<port>. Default localhost:3000",
		Destination: &listenStringFlag,
		Value:       "localhost:3000",
	},
	cli.StringFlag{
		Name:        "usecache, cache",
		Usage:       "Use in-memory cache yes/no/populate?. Default - no",
		Destination: &useCacheFlag,
		Value:       "no",
	},
	cli.StringFlag{
		Name:        "boltconfig, boltconf",
		Usage:       "Bolt DB config file name.",
		Destination: &boltConfigFileFlag,
		Value:       "",
	},
	cli.StringFlag{
		Name:        "JSONCache, JSON",
		Usage:       "JSON cache type memory/bolt/pass?. Default - pass",
		Destination: &JSONCacheTypeFlag,
		Value:       "",
	},
	cli.StringFlag{
		Name:        "populateJSONCache, populateJSON",
		Usage:       "Populate JSON cache yes/no?. Default - no",
		Destination: &populateJSONCacheFlag,
		Value:       "",
	},
	cli.StringFlag{
		Name:        "debug, d",
		Usage:       "debug mode: DEBUG, WARN, INFO, ERROR. Default ERROR",
		Destination: &debugFlag,
		Value:       "ERROR",
		Hidden:      true,
	},
}

// Глобальные переменные в которые парсятся входные флаги
var dbConfigFileFlag string
var listenStringFlag string
var debugFlag string
var useCacheFlag string
var boltConfigFileFlag string
var JSONCacheTypeFlag string
var populateJSONCacheFlag string

// checkFlags проверяет входные флаги и входные аргументы программы
func setConfig(cfg *daemon.Config) error {
	//	defer PrintElapsed("[checkFlags takes ")()

	cfg.LogFilter = mylog.LogFilter

	// Проверяем, что конфигурационный файл для соединения с PostgreSQL указан
	if dbConfigFileFlag == "" {
		log.Fatalf("[ERROR] Database config file name is null")
	} else {
		cfg.PqDbCfg.ConfigFile = dbConfigFileFlag
	}

	// адрес и порт HTTP сервера
	if listenStringFlag == "" {
		log.Fatalf("[ERROR] HTTP listener string is null")
	} else {
		cfg.ListenSpec = listenStringFlag
	}

	// ===============================================================================
	// Тестовые значения для отладки Notify PostgreSQL
	cfg.PqDbCfg.ChanelPrefix = "events"
	cfg.PqDbCfg.ListenerMinReconnectInterval = 10 * time.Second
	cfg.PqDbCfg.ListenerMaxReconnectInterval = time.Minute
	cfg.PqDbCfg.ListenerWaitForNotificationInterval = 60 * time.Second
	// ===============================================================================

	// Если запуск с опцией отладки
	if debugFlag != "" {
		switch debugFlag {
		case "DEBUG":
			mylog.LogFilter.SetMinLevel(logutils.LogLevel("DEBUG"))
		case "WARN":
			mylog.LogFilter.SetMinLevel(logutils.LogLevel("WARN"))
		case "ERROR":
			mylog.LogFilter.SetMinLevel(logutils.LogLevel("ERROR"))
		case "INFO":
			mylog.LogFilter.SetMinLevel(logutils.LogLevel("INFO"))
		default:
			log.Fatalf("[ERROR] debug flag have to be: DEBUG, WARN, INFO, ERROR")
		}
	}

	// Флаг использования кэша
	if useCacheFlag != "" {
		switch useCacheFlag {
		case "yes":
			cfg.UseCache = "yes"
			cfg.PqDbCfg.UseNotification = true
		case "no":
			cfg.UseCache = "no"
			cfg.PqDbCfg.UseNotification = false
		case "populate":
			cfg.UseCache = "populate"
			cfg.PqDbCfg.UseNotification = true
		default:
			log.Fatalf("[ERROR] useCache flag have to be: 'yes', 'no', 'populate'")
		}
	} else {
		cfg.PqDbCfg.UseNotification = false
	}

	// Временные настройки JSON Cache
	cfg.JSONCfg.DeptBucketName = "DeptBucket"
	cfg.JSONCfg.SaveChSize = 5000000
	cfg.JSONCfg.SaveQueueSize = 10000

	// Флаг использования JSON кэша
	if JSONCacheTypeFlag != "" {
		switch JSONCacheTypeFlag {
		case "memory":
			cfg.JSONCfg.IsPassCache = false
		case "bolt":
			cfg.JSONCfg.IsPassCache = false
		case "pass":
			cfg.JSONCfg.IsPassCache = true
		default:
			log.Fatalf("[ERROR] JSONCacheTypeFlag flag have to be: 'memory', 'bolt', 'pass'")
		}
	} else {
		cfg.JSONCfg.IsPassCache = true
	}

	// Проверяем указан ли конфигурационный файл для соединения с Bolt BD
	if JSONCacheTypeFlag == "bolt" {
		if boltConfigFileFlag != "" {
			cfg.BoltDbCfg.ConfigFile = boltConfigFileFlag
			cfg.UseBoltDb = true
		} else {
			log.Fatalf("[ERROR] Bolt DB config file name must to be NOT NULL")
		}
	} else {
		cfg.BoltDbCfg.ConfigFile = ""
		cfg.UseBoltDb = false
	}

	// Флаг наполнять JSON кэш
	if JSONCacheTypeFlag == "bolt" || JSONCacheTypeFlag == "memory" {
		if populateJSONCacheFlag != "" {
			switch populateJSONCacheFlag {
			case "yes":
				cfg.JSONCfg.IsPopulateCache = true
			case "no":
				cfg.JSONCfg.IsPopulateCache = false
			default:
				log.Fatalf("[ERROR] populateJSON flag have to be: 'yes', 'no'")
			}
		} else {
			cfg.JSONCfg.IsPopulateCache = false
		}
	} else {
		cfg.JSONCfg.IsPopulateCache = false
	}

	return nil
}

func main() {
	// Переопределяем глобальный логер на кастомный
	log.SetOutput(mylog.LogFilter)

	app := cli.NewApp() // создаем приложение
	app.Name = "REST test application"
	app.Version = "0.1.0"
	app.Usage = "REST test application"
	app.UsageText = "restDept [-db PostgresDB name -l listenstring <host>:<port>]"

	app.Flags = flags // присваиваем ранее определенные флаги

	// Определяем единственной действие - запуск демона
	app.Action = func(c *cli.Context) error {

		// Настраиваем конфигурацию демона
		cfg := &daemon.Config{}
		err := setConfig(cfg) // Записываем конфигурацию из флагов и внешний переменных
		if err != nil {
			return errors.Wrap(err, "setConfig(cfg)")
		}

		err = daemon.Run(cfg) // Запускаем демон

		if err != nil {
			return errors.Wrap(err, "daemon.Run(cfg)")
		}
		return nil
	}

	err := app.Run(os.Args) // Запускаем приложение

	// Выводим в лог стек вызова
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
