package http

// Обертка над стандартным пакетом http используется для изоляции HTTP кода и обработчиков

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// EmpGetHandler handle get emp
// =====================================================================
func (h *Handler) EmpGetHandler(w http.ResponseWriter, r *http.Request) {
	cnt := "http (h *Handler) EmpGetHandler" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// Считаем параметры и проверим на число
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR - incorrect id = %+v", cnt, idStr)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// получаем Json
	buf, err := h.GetEmp(id)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR - h.GetEmp(id), args = '%v'\n %+v", cnt, id, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}
	// Если данные не найдены
	if err == nil && buf == nil {
		//mylog.PrintfDebug("[DEBUG] %v - h.GetEmp(id), args = '%v' - NO_DATA_FOUND", cnt, id)
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

// EmpAddHandler handle get emp
// =====================================================================
func (h *Handler) EmpAddHandler(w http.ResponseWriter, r *http.Request) {
	cnt := "http (h *Handler) EmpAddHandler" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// считаем тело запроса
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] EmpAddHandler - ioutil.ReadAll(r.Body) \n %+v", err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusBadRequest)
		return
	}

	// передаем Json на обработку и получаем в ответ изменения после создания
	bufNew, err := h.CreateEmp(buf)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR -h.CreateEmp(buf), args = '%v'\n %+v", cnt, buf, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	// Успешное создание
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	// записываем Json в ответ
	_, err = w.Write(bufNew)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR - w.Write(bufNew) \n %+v", cnt, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
}
