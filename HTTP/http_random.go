package http

// Обертка над стандартным пакетом http используется для изоляции HTTP кода и обработчиков

import (
	"fmt"
	"log"
	"net/http"
	//mylog "github.com/romapres2010/restdept/log"
)

// RandomDeptGetHandler handle get dept
// =====================================================================
func (h *Handler) RandomDeptGetHandler(w http.ResponseWriter, r *http.Request) {
	cnt := "http (h *Handler) RandomDeptGetHandler" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// получаем Json
	buf, err := h.RandomGetDept()
	defer h.PutByteSlice(buf) // вернем в pool для повторного использования. Safe for nil pointer

	if err != nil {
		log.Printf("[ERROR] %v - ERROR - h.RandomGetDept() \n %+v", cnt, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}
	// Если данные не найдены
	if err == nil && buf == nil {
		//mylog.PrintfDebug("[DEBUG] %v - h.RandomGetDept() - NO_DATA_FOUND", cnt)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Устанавливаем тип контента
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// записываем Json в ответ
	_, err = w.Write(buf)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR - w.Write(buf) \n %+v", cnt, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
}

// RandomDeptUpdateHandler handle get dept
// =====================================================================
func (h *Handler) RandomDeptUpdateHandler(w http.ResponseWriter, r *http.Request) {
	cnt := "http (h *Handler) RandomDeptUpdateHandler" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// получаем Json
	buf, err := h.RandomUpdateDept()
	if err != nil {
		log.Printf("[ERROR] %v - ERROR - h.RandomUpdateDept() \n %+v", cnt, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}
	// Если данные не найдены
	if err == nil && buf == nil {
		//mylog.PrintfDebug("[DEBUG] %v - h.RandomUpdateDept() - NO_DATA_FOUND", cnt)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Устанавливаем тип контента
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// записываем Json в ответ
	_, err = w.Write(buf)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR - w.Write(buf) \n %+v", cnt, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
}
