package json

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/francoispqt/gojay"
	"github.com/pkg/errors"
	"github.com/romapres2010/restdept/model"
)

var testDeptPK []*model.DeptPK

// RandomGetDept return a JSON for a given id and all Emps
// =====================================================================
func (s *Service) RandomGetDept() ([]byte, error) {
	cnt := "json (s *Service)  RandomGetDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	var err error

	// Считаем все PK
	if testDeptPK == nil {
		// считаем все PK для dept
		testDeptPK, err = s.deptService.GetDeptsPK()
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.GetDeptsPK()", cnt)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
	}
	deptsPKlen := len(testDeptPK)

	// случайным образом выбираем PK
	par := testDeptPK[rand.Intn(deptsPKlen)].Deptno

	// получаем JSON
	buf, err := s.GetDeptEmps(par)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptTx(par), args = '%v'", cnt, par)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	return buf, nil
}

// RandomUpdateDept update dept and return a JSON for a given id and all Emps
// =====================================================================
func (s *Service) RandomUpdateDept() ([]byte, error) {
	cnt := "json (s *Service)  RandomUpdateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// считаем данные из сервиса
	newV, err := s.deptService.RandomUpdateDept()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - s.deptService.RandomUpdateDept()", cnt)
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
