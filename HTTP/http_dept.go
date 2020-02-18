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

// DeptGetHandler handle get dept
// =====================================================================
func (h *Handler) DeptGetHandler(w http.ResponseWriter, r *http.Request) {
	cnt := "http (h *Handler) DeptGetHandler" // Имя текущего метода для логирования
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
	buf, err := h.GetDept(id)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR - h.GetDept(id), args = '%v'\n %+v", cnt, id, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}
	// Если данные не найдены
	if err == nil && buf == nil {
		//mylog.PrintfDebug("[DEBUG] %v - h.GetDept(id), args = '%v' - NO_DATA_FOUND", cnt, id)
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

	buf = nil
	return
	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
}

// DeptEmpsGetHandler handle get deptemps
// =====================================================================
func (h *Handler) DeptEmpsGetHandler(w http.ResponseWriter, r *http.Request) {
	cnt := "http (h *Handler) DeptEmpsGetHandler" // Имя текущего метода для логирования
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
	buf, err := h.GetDeptEmps(id)
	defer h.PutByteSlice(buf) // вернем в pool для повторного использования. Safe for nil pointer
	if err != nil {
		log.Printf("[ERROR] %v - ERROR - h.GetDept(id), args = '%v'\n %+v", cnt, id, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}
	// Если данные не найдены
	if err == nil && buf == nil {
		//mylog.PrintfDebug("[DEBUG] %v - h.GetDept(id), args = '%v' - NO_DATA_FOUND", cnt, id)
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

	buf = nil
	return
	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
}

// DeptAddHandler handle add dept
// =====================================================================
func (h *Handler) DeptAddHandler(w http.ResponseWriter, r *http.Request) {
	cnt := "http (h *Handler) DeptAddHandler" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// считаем тело запроса
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] DeptAddHandler - ioutil.ReadAll(r.Body) \n %+v", err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusBadRequest)
		return
	}

	// передаем Json на обработку и получаем в ответ изменения после создания
	bufNew, err := h.CreateDept(buf)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR -h.CreateDept(buf), args = '%v'\n %+v", cnt, buf, err)
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

// DeptUpdateHandler handle update dept
// =====================================================================
func (h *Handler) DeptUpdateHandler(w http.ResponseWriter, r *http.Request) {
	cnt := "http (h *Handler) DeptUpdateHandler" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] --------------------------------------------------------------------------")
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// считаем тело запроса
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] DeptUpdateHandler - ioutil.ReadAll(r.Body) \n %+v", err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusBadRequest)
		return
	}

	// передаем Json на обработку и получаем в ответ изменения после обновления
	bufNew, err := h.UpdateDept(buf)
	if err != nil {
		log.Printf("[ERROR] %v - ERROR -h.CreateDept(buf), args = '%v'\n %+v", cnt, buf, err)
		http.Error(w, fmt.Sprintf("%+v", err), http.StatusInternalServerError)
		return
	}
	// Если данные не найдены
	if err == nil && bufNew == nil {
		log.Printf("[INFO] DeptUpdateHandler - call h.UpdateDept(buf) NO_DATA_FOUND")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Успешное обновление
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

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
