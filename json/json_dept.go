package json

import (
	"fmt"
	"log"

	"github.com/francoispqt/gojay"
	"github.com/pkg/errors"

	//mylog "github.com/romapres2010/restdept/log"
	"github.com/romapres2010/restdept/model"
)

// GetDept return a JSON for a given id
// =====================================================================
func (s *Service) GetDept(id int) ([]byte, error) {
	cnt := "json (s *Service)  GetDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// считаем данные из сервиса
	v, err := s.deptService.GetDept(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.GetDept(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// Если не используется кэш, то возврашаем в пул
	if !s.deptService.IsUseCache() {
		defer model.PutDept(v, true)
	}

	// сформируем json
	if v != nil {
		// сформируем json
		buf, err := gojay.Marshal(v)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - gojay.Marshal(v), args = '%+v'", cnt, v)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
		//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
		return buf, nil
	}
	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return nil, nil
}

// GetDeptEmps return a JSON for a given id and all Emps
// =====================================================================
func (s *Service) GetDeptEmps(id int) ([]byte, error) {
	cnt := "json (s *Service)  GetDeptEmps" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	var val *cacheStr
	var inCache bool

	if !s.isPassCache {
		// Ищем в кэше
		val, inCache = s.deptCache.get(id)

		/*
			{ // отладка
				if inCache && val != nil {
					//mylog.PrintfDebug("[DEBUG] %v - s.deptCache.get(id), args = '%+v', IN CACHE - val.flag ='%+v'", cnt, id, val.flag)
				} else if inCache && val == nil {
					//mylog.PrintfDebug("[DEBUG] %v - s.deptCache.get(id), args = '%+v', IN CACHE - val IS NULL", cnt, id)
				} else if !inCache {
					//mylog.PrintfDebug("[DEBUG] %v - s.deptCache.get(id), args = '%+v', NOT IN CACHE", cnt, id)
				}
			} // отладка
		*/

		// Если данные есть в кэше и буфер памяти заполнен
		if val != nil && inCache && val.flag == JSON_IN_MEMORY {
			if s.db != nil {
				// во время обработки буфер может быть сброшел при записи в Bolt DB
				// добавить копирование в дополнительный буфер и mutex
			}
			// Некорректная ситуация - на время отладки
			if val.buf == nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptCache.get(id), args = '%+v' - flag = JSON_IN_MEMORY but val.buf IS NULL", cnt, id)
				log.Printf(errM)
				return nil, errors.New(errM)
			}
			return val.buf, nil
		}

		// Если буфер в кэше Blob DB, либо не нашли в памяти - попробуем поискать в Bolt DB
		if !inCache || (inCache && val.flag == JSON_IN_BLOB_DB) {
			if s.db != nil {
				// Считаем Bolt DB
				retbuf, err := s.db.Get(s.deptBucket, id)
				if err != nil {
					errM := fmt.Sprintf("[ERROR] %v - ERROR - s.db.Get(s.deptBucket, id), args = '%+v', '%+v'", cnt, s.deptBucket, id)
					log.Printf(errM)
					return nil, errors.WithMessage(err, errM)
				}
				// Если данные найдены
				if retbuf != nil {
					//mylog.PrintfDebug("[DEBUG] %v - s.db.Get(s.deptBucket, id), args = '%+v', '%+v', '%p' - JSON has been found in DB", cnt, s.deptBucket, id)
					// после первого включения приложения, записи может не быть в памяти - добавляем
					if val == nil {
						//mylog.PrintfDebug("[DEBUG] %v - s.deptCache.set(id, nil, JSON_IN_BLOB_DB), args = '%+v' - add to Cache", cnt, id)

						// запишем c пустым кэшем и флагом JSON_IN_BLOB_DB
						s.deptCache.set(id, nil, JSON_IN_BLOB_DB)
						//s.deptCache.set(id, retbuf, JSON_IN_MEMORY) // для тестирования памяти
					}
					//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
					return retbuf, nil
				}
				// Нормальная ситуация - продолжим обработку
				//mylog.PrintfDebug("[DEBUG] %v - s.db.Get(s.deptBucket, id), args = '%+v', '%+v' - JSON DID NOT find in DB", cnt, s.deptBucket, id)
			}
		}
	}

	/*
		{ // отладка
			if inCache {
				//mylog.PrintfDebug("[INFO] %v - Strange !!!! - inCahce=True but not in Bolt, args = '%v'", cnt, id)
			}
		} // отладка
	*/

	// Запрашиваем базовый сервис

	// считаем родительскую запись из сервиса
	dept, err := s.deptService.GetDept(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.GetDept(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// Если не используется кэш сервиса, то возврашаем структуру в пул
	if !s.deptService.IsUseCache() {
		defer model.PutDept(dept, false)
	}

	// Если данные не найдены
	if dept == nil {
		//mylog.PrintfDebug("[DEBUG] %v - NO_DATA_FOUND", cnt)
		return nil, nil
	}

	// Считываем подчиненные записи
	dept.Emps, err = s.empService.GetEmpsByDept(dept.Deptno)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.GetEmpsByDept(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// сформируем json
	retbuf, err := gojay.Marshal(dept)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - gojay.Marshal(), args = '%+v'", cnt, dept)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// Если не прокидываем кэш и ранее не было записи в кэше
	if !inCache && !s.isPassCache {
		//mylog.PrintfDebug("[DEBUG] %v - s.deptCache.set(dept.Deptno, buf, JSON_IN_MEMORY), args = '%+v' - add to Cache", cnt, dept.Deptno)
		valNew := s.deptCache.set(dept.Deptno, retbuf, JSON_IN_MEMORY) // запишем структуру в кэш

		// Если используется запись в Bolt DB, то отправляем сообщение в канал для подготовки к записи
		if s.db != nil {
			s.saveListener.sendSaveMessage(valNew)
		}
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return retbuf, nil
}

// CreateDept create dept and return a JSON for a given id
// =====================================================================
func (s *Service) CreateDept(buf []byte) ([]byte, error) {
	cnt := "json (s *Service)  CreateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// Извлечем из pool новую структуру
	v := model.GetDept()
	defer model.PutDept(v, true) // возвращаем в pool струкуру

	// Парсим JSON в структуру
	err := gojay.Unmarshal(buf, v)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - gojay.Unmarshal(v), args = '%+v'", cnt, v)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// считаем данные из сервиса
	newV, err := s.deptService.CreateDept(v)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.CreateDept(v), args = '%v'", cnt, v)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// Если не используется кэш, то возврашаем в пул
	if !s.deptService.IsUseCache() {
		defer model.PutDept(newV, true)
	}

	// сформируем json
	if newV != nil {
		bufNew, err := gojay.Marshal(newV)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - gojay.Marshal(newV), args = '%+v'", cnt, newV)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
		//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
		return bufNew, nil
	}

	// Данные не менялись - вернуть исходную информацию
	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return buf, nil
}

// UpdateDept update dept and return a JSON for a given id
// =====================================================================
func (s *Service) UpdateDept(buf []byte) ([]byte, error) {
	cnt := "json (s *Service)  UpdateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// Извлечем из pool новую структуру
	v := model.GetDept()
	defer model.PutDept(v, true) // возвращаем в pool струкуру

	// Парсим JSON в структуру
	err := gojay.Unmarshal(buf, v)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - gojay.Unmarshal(v), args = '%+v'", cnt, v)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// считаем данные из сервиса
	newV, err := s.deptService.UpdateDept(v)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.UpdateDept(v), args = '%v'", cnt, v)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// Если не используется кэш, то возврашаем в пул
	if !s.deptService.IsUseCache() {
		defer model.PutDept(newV, true)
	}

	// сформируем json
	if newV != nil {
		bufNew, err := gojay.Marshal(newV)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - gojay.Marshal(newV), args = '%+v'", cnt, newV)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
		//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
		return bufNew, nil
	}

	// Данные не найдены
	//mylog.PrintfDebug("[DEBUG] %v - s.deptService.UpdateDept(v) NO_DATA_FOUND", cnt)
	return nil, nil
}
