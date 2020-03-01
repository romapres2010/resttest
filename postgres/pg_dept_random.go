package postgres

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/pkg/errors"

	"github.com/romapres2010/resttest/model"
)

var testDeptPK []*model.DeptPK

// RandomGetDept return a random row
// =====================================================================
func (p *PgDb) RandomGetDept() (*model.Dept, error) {
	cnt := "postgres (p *PgDb) RandomGetDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)
	var err error

	// Считаем все PK
	if testDeptPK == nil {
		// считаем все PK для dept
		testDeptPK, err = p.GetDeptsPK()
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptsPK()", cnt)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
	}
	deptsPKlen := len(testDeptPK)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// случайным образом выбираем PK
	par := testDeptPK[rand.Intn(deptsPKlen)].Deptno
	// Вызов НЕ в транзакционном режиме
	v, err := p.GetDeptEmpsTx(par, nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptEmpsTx(par, nil), args = '%v'", cnt, par)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// RandomUpdateDept update a random row
// =====================================================================
func (p *PgDb) RandomUpdateDept() (*model.Dept, error) {
	cnt := "postgres (p *PgDb) RandomUpdateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)
	var err error

	// Считаем все PK
	if testDeptPK == nil {
		// считаем все PK для dept
		testDeptPK, err = p.GetDeptsPK()
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptsPK()", cnt)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
	}
	deptsPKlen := len(testDeptPK)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// случайным образом выбираем PK
	par := testDeptPK[rand.Intn(deptsPKlen)].Deptno

	// Начинаем транзакцию
	tx, err := p.Db.Beginx()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.Db.Beginx()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// Считаем dept и все emp в транзакционном режиме
	v, err := p.GetDeptEmpsTx(par, tx)
	defer model.PutDept(v, true) // возвращаем в кэш структуру, в которую считывали данные
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptEmpsTx(par, nil), args = '%v'", cnt, par)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// Обновляем dept и все emp Вызов в транзакционном режиме
	vNew, err := p.UpdateDeptTx(v, tx)
	if err != nil {
		defer tx.Rollback()
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.UpdateDeptTx(v, tx), args = '%+v'", cnt, v)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// завершаем транзакцию
	if err := tx.Commit(); err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.Commit()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return vNew, nil
}
