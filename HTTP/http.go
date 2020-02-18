package http

// Обертка над стандартным пакетом http используется для изоляции HTTP кода и обработчиков

import (
	"github.com/romapres2010/restdept/json"
)

// Handler - представляет обертку над обработчиком для реализации всех HTTP методов
type Handler struct {
	*json.Service
}

// New - создает новый handler
func New(s *json.Service) *Handler {
	return &Handler{s}
}
