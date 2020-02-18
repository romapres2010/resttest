package model

import (
	"errors"
	"log"
	"time"
)

// Notification represents a single notification
type Notification struct {
	ObjectName    string // наименование таблицы
	OperationName string // наименование операции
	Id            int    // идентификатор
}

// NotificationService represents service for process notification
type NotificationService interface {
	ProcessNotification(n *Notification) error
}

// Listener represents a chanel for notifications
type Listener struct {
	Notify chan *Notification

	Service NotificationService
	/*
		name                 string
		minReconnectInterval time.Duration
		maxReconnectInterval time.Duration
		dialer               Dialer
		eventCallback        EventCallbackType

		lock                 sync.Mutex
		isClosed             bool
		reconnectCond        *sync.Cond
		cn                   *ListenerConn
		connNotificationChan <-chan *Notification
		channels             map[string]struct{}
	*/
}

// NewListener create new listener
func NewListener(s NotificationService) *Listener {
	l := &Listener{
		Notify:  make(chan *Notification, 32),
		Service: s,
	}

	go l.startListener()

	return l
}

// SetNotificationService set Notification Service
func (l *Listener) SetNotificationService(s NotificationService) {
	l.Service = s
}

// startListener - запускаем бесконечный цикл для ожидания уведомлений
func (l *Listener) startListener() {
	cnt := "model (l *Listener) startListener()" // Имя текущего метода для логирования
	log.Printf("[INFO] %v - START", cnt)

	for {
		l.waitForNotification()
	}
}

// waitForNotification - ожидание уведомления
func (l *Listener) waitForNotification() {
	cnt := "model (l *Listener) waitForNotification" // Имя текущего метода для логирования
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

	select {
	case n := <-l.Notify:
		if l.Service != nil {
			l.Service.ProcessNotification(n)
		} else {
			log.Printf("[ERROR] %v - l.Service is NULL", cnt)
		}
		// Не получили уведомления - проверям соединение
	case <-time.After(time.Minute):
		log.Printf("[INFO] %v - Received no events for '%v'", cnt, time.Minute)
	}
}
