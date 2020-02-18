package json

import (
	"fmt"
	"log"

	"github.com/francoispqt/gojay"
	"github.com/pkg/errors"
	"github.com/romapres2010/restdept/model"
)

// GetEmp return a JSON for a given id
// =====================================================================
func (s *Service) GetEmp(id int) ([]byte, error) {
	cnt := "json (s *Service)  GetEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// считаем данные из сервиса
	v, err := s.empService.GetEmp(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.GetEmp(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// Если не используется кэш, то возврашаем в пул
	if !s.empService.IsUseCache() {
		defer model.PutEmp(v)
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

// CreateEmp create emp and return a JSON for a given id
// =====================================================================
func (s *Service) CreateEmp(buf []byte) ([]byte, error) {
	cnt := "json (s *Service)  CreateEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// Извлечем из pool новую структуру
	v := model.GetEmp()
	defer model.PutEmp(v) // возвращаем в pool струкуру

	// Парсим JSON в структуру
	err := gojay.Unmarshal(buf, v)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - gojay.Unmarshal(v), args = '%+v'", cnt, v)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// считаем данные из сервиса
	newV, err := s.empService.CreateEmp(v)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - s.empService.CreateEmp(v), args = '%v'", cnt, v)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// Если не используется кэш, то возврашаем в пул
	if !s.empService.IsUseCache() {
		defer model.PutEmp(newV)
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
