package postgres

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/romapres2010/restdept/model"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// SetModelListener set model Listener for notifications
// =====================================================================
func (p *PgDb) SetModelListener(l *model.Listener) {
	cnt := "postgers (p *PgDb) SetModelListener" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	p.modelListener = l
}

// createListenerForNotification - создаем листенер для ожидания уведомлений и запускаем бесконеный цикл
func (p *PgDb) createListenerForNotification(cfg Config) error {
	cnt := "postgers (p *PgDb) createListenerForNotification" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	var err error

	// создаем листенер для отслеживания Notify БД Postgres
	if cfg.ChanelPrefix == "" {
		log.Printf("[INFO]  %v - PostgreSQL Notify ChanelPrefix is NULL. Listener will not created", cnt)
		return nil
	} else {
		// Сформировать ошибочную строку с параметрами для теста
		// cfg.ConnectString = "host=127.0.0.1 port=54322 dbname=test_database sslmode=disable user=test_user password=qwerty"

		log.Printf("[INFO] %v - Creating PostgreSQL Notify listener", cnt)
		p.pqListener = pq.NewListener(cfg.ConnectString, cfg.ListenerMinReconnectInterval, cfg.ListenerMaxReconnectInterval, listenerEventCallback)
		p.pqListenerWaitForNotificationInterval = cfg.ListenerWaitForNotificationInterval

		// =======================================================
		// !!!!!! переделать
		// =======================================================
		p.pqChanelDeptName = cfg.ChanelPrefix + "_dept"
		p.pqChanelEmpName = cfg.ChanelPrefix + "_emp"
		// =======================================================

		// Запускаем прослущивание каналов
		err = p.listenChanel(p.pqChanelDeptName)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - PgDb.Listener.Listen('%v') \n %+v", cnt, p.pqChanelDeptName, err)
			log.Printf(errM)
			return errors.WithMessage(err, errM)
		}
		err = p.listenChanel(p.pqChanelEmpName)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - PgDb.Listener.Listen('%v') \n %+v", cnt, p.pqChanelEmpName, err)
			log.Printf(errM)
			return errors.WithMessage(err, errM)
		}

		// Запускаем бесконечный цикл для отслеживания уведомлений
		log.Printf("[INFO] %v - Starting loop wait for PostgreSQL notification in separate GO", cnt)
		go p.startWaitForNotification()
	}

	// Создаем внутренний листенер на уровне model
	//	p.modelListener = model.NewListener()

	log.Printf("[INFO] %v - SUCCESS", cnt)

	return nil
}

// createListenerForNotification - создаем листенер для ожидания уведомлений и запускаем бесконеный цикл
func (p *PgDb) listenChanel(chanel string) error {
	cnt := "postgers (p *PgDb) listenChanel" // Имя текущего метода для логирования
	log.Printf("[INFO] %v(%v) - START", cnt, chanel)

	err := p.pqListener.Listen(chanel) // !!!!при отсутсвии ответа от сервера бесконечно блокирует и вызывает listenerEventCallback
	if err != nil {
		// Если уже открыт, то ни чего не делаем
		if err == pq.ErrChannelAlreadyOpen {
			log.Printf("[INFO] %v - p.Listener.Listen('%v') - ErrChannelAlreadyOpen", cnt, chanel)
		} else { // другие ошибки
			errM := fmt.Sprintf("[ERROR] %v - ERROR - p.Listener.Listen('%v') \n %+v", cnt, chanel, err)
			log.Printf(errM)
			return errors.Wrap(err, errM)
		}
	}
	log.Printf("[INFO] %v(%v) - SUCCESS", cnt, chanel)

	return nil
}

// listenerEventCallback - функция для обработки событий и ошибок от листенера Notify
func listenerEventCallback(ev pq.ListenerEventType, err error) {
	cnt := "postgers listenerEventCallback" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	// Определяем тип события
	switch ev {
	case pq.ListenerEventConnected: // нормальная ситуация - все в порядке
		log.Printf("[INFO] %v - ListenerEventConnected", cnt)
	case pq.ListenerEventReconnected: // нормальная ситуация - все в порядке
		log.Printf("[INFO] %v - ListenerEventReconnected", cnt)
	case pq.ListenerEventDisconnected:
		log.Printf("[ERROR] %v - ERROR - ListenerEventDisconnected \n %+v", cnt, err)
	case pq.ListenerEventConnectionAttemptFailed:
		log.Printf("[ERROR] %v - ERROR - ListenerEventConnectionAttemptFailed \n %+v", cnt, err)
	}
}

// startWaitForNotification - запускаем бесконечный цикл для ожидания уведомлений
func (p *PgDb) startWaitForNotification() {
	cnt := "postgers (p *PgDb) startWaitForNotification" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	if p.pqListener != nil {
		for {
			p.waitForNotification()
		}
	}
}

// waitForNotification - ожидание уведомления
func (p *PgDb) waitForNotification() {
	cnt := "postgers (p *PgDb) waitForNotification" // Имя текущего метода для логирования

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

	select {
	case n := <-p.pqListener.Notify:
		// При переподключении к БД может выдавать разовый nil
		if n != nil {
			//mylog.PrintfDebug("[DEBUG] %v - Received data from channel '%+v' n.Extra = '%v'", cnt, n.Channel, n.Extra)

			// Обрабатываем в зависимости от типа канала
			if n.Channel == p.pqChanelDeptName {
				// Пустые сообщения игнорируем
				if n.Extra == "" {
					log.Printf("[ERROR] %v - p.pqChanelDeptName - n.Extra is NULL", cnt)
					return
				}

				// преобразуем сообщение в число
				id, err := strconv.Atoi(n.Extra)
				if err != nil {
					log.Printf("[ERROR] %v - p.pqChanelDeptName - strconv.Atoi(n.Extra) - incorrect number id='%v'", cnt, id)
					return
				}

				// Нулевые идентификаторы не обрабатываем
				if id == 0 {
					log.Printf("[ERROR] %v - p.pqChanelDeptName - strconv.Atoi(n.Extra) - incorrect id='0'", cnt)
					return
				}

				// Сформируем текст сообщения
				nt := &model.Notification{
					ObjectName:    "dept",
					OperationName: "outDate",
					Id:            id,
				}
				//mylog.PrintfDebug("[DEBUG] %v - p.pqChanelDeptName - sent into cahnel 'p.modelListener.Notify' - '%+v'", cnt, nt)

				// Передаем сообщение в канал
				if p.modelListener != nil {
					p.modelListener.Notify <- nt
				} else { // мало вероятная ситуация, но на всякий случай
					log.Printf("[ERROR] %v - p.modelListener is NULL", cnt)
				}
			} else if n.Channel == p.pqChanelEmpName {
				// Пустые сообщения игнорируем
				if n.Extra == "" {
					log.Printf("[ERROR] %v - p.pqChanelEmpName - n.Extra is NULL", cnt)
					return
				}

				// преобразуем сообщение в число
				id, err := strconv.Atoi(n.Extra)
				if err != nil {
					log.Printf("[ERROR] %v - p.pqChanelEmpName - strconv.Atoi(n.Extra) - incorrect number id='%v'", cnt, id)
					return
				}

				// Нулевые идентификаторы не обрабатываем
				if id == 0 {
					log.Printf("[ERROR] %v - p.pqChanelEmpName - strconv.Atoi(n.Extra) - incorrect id='0'", cnt)
					return
				}

				// Сформируем текст сообщения
				nt := &model.Notification{
					ObjectName:    "emp",
					OperationName: "outDate",
					Id:            id,
				}
				//mylog.PrintfDebug("[DEBUG] %v - p.pqChanelEmpName - sent into cahnel 'p.modelListener.Notify' - '%+v'", cnt, nt)

				// Передаем сообщение в канал
				if p.modelListener != nil {
					p.modelListener.Notify <- nt
				} else { // мало вероятная ситуация, но на всякий случай
					log.Printf("[ERROR] %v - p.modelListener is NULL", cnt)
				}
			}
			/*
				// Разберем ответ в JSON
				var prettyJSON bytes.Buffer
				err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
				if err != nil {
					log.Printf("[ERROR] %v - ERROR - json.Indent(&prettyJSON, []byte(n.Extra)) '%v'", cnt, err)
					return
				}
				//mylog.PrintfDebug("[DEBUG] %v - Json value = '%v'", cnt, string(prettyJSON.Bytes()))
			*/
		}

	// Не получили уведомления - проверям соединение
	case <-time.After(p.pqListenerWaitForNotificationInterval):
		log.Printf("[INFO] %v - Received no events for '%v', checking connection", cnt, p.pqListenerWaitForNotificationInterval)
		// запускаем Ping для проверки подключения
		go func() {
			err := p.pqListener.Ping()
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - p.pqListener.Ping() /n %+v", cnt, err)
				log.Printf(errM)
			}
		}()
	}
}
