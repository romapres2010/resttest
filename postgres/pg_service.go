package postgres

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/romapres2010/restdept/model"
)

// IsUseCache returns FALSE по Postgres DB
// =====================================================================
func (p *PgDb) IsUseCache() bool {
	return false
}

// DeptExists returns TRUE if row for a given id exists
// =====================================================================
func (p *PgDb) DeptExists(id int) (bool, error) {
	cnt := "postgers (p *PgDb) DeptExists" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.DeptExistsTx(id, nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.DeptExistsTx(id, nil), args = '%v'", cnt, id)
		log.Printf(errM)
		return false, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// EmpExists returns TRUE if row for a given id exists
// =====================================================================
func (p *PgDb) EmpExists(id int) (bool, error) {
	cnt := "postgers (p *PgDb) EmpExists" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.EmpExistsTx(id, nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.EmpExistsTx(id, nil), args = '%v'", cnt, id)
		log.Printf(errM)
		return false, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// GetDept return a row for a given id without Emps. Non TX mode
// =====================================================================
func (p *PgDb) GetDept(id int) (*model.Dept, error) {
	cnt := "postgers (p *PgDb) GetDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.GetDeptTx(id, nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptTx(id, nil), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// GetDeptEmps return a row for a given id and all Emps. Non TX mode
// =====================================================================
func (p *PgDb) GetDeptEmps(id int) (*model.Dept, error) {
	cnt := "postgers (p *PgDb) GetDeptEmps" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.GetDeptEmpsTx(id, nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptEmpsTx(id, nil), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// GetEmpsByDept returns all emps for a given dept
// =====================================================================
func (p *PgDb) GetEmpsByDept(id int) ([]*model.Emp, error) {
	cnt := "postgers (p *PgDb) GetEmpsByDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.GetEmpsByDeptTx(id, nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetEmpsByDeptTx(id, nil), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// GetEmpsPKByDept returns all emps for a given dept
// =====================================================================
func (p *PgDb) GetEmpsPKByDept(id int) ([]*model.EmpPK, error) {
	cnt := "postgers (p *PgDb) GetEmpsPKByDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.GetEmpsPKByDeptTx(id, nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetEmpsPKByDeptTx(id, nil), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// GetDeptsPK returns Primary Key (PK) for all depts
// =====================================================================
func (p *PgDb) GetDeptsPK() ([]*model.DeptPK, error) {
	cnt := "postgers (p *PgDb) GetDeptsPK" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.GetDeptsPKTx(nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptsPKTx(tx)", cnt)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// GetDepts returns all depts
// =====================================================================
func (p *PgDb) GetDepts() ([]*model.Dept, error) {
	cnt := "postgers (p *PgDb) GetDepts" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.GetDeptsTx(nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDepts(tx)", cnt)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// GetEmp return a row for a given id
// =====================================================================
func (p *PgDb) GetEmp(id int) (*model.Emp, error) {
	cnt := "postgers (p *PgDb) GetEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов НЕ в транзакционном режиме
	v, err := p.GetEmpTx(id, nil)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetEmpTx(id, tx), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// CreateEmp create new emp
// =====================================================================
func (p *PgDb) CreateEmp(r *model.Emp) (*model.Emp, error) {
	cnt := "postgers (p *PgDb) CreateEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// Начинаем транзакцию
	tx, err := p.Db.Beginx()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.Db.Beginx()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов в транзакционном режиме
	v, err := p.CreateEmpTx(r, tx, true)
	// =====================================================================
	if err != nil {
		defer tx.Rollback()
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.CreateEmpTx(id, tx)", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// завершаем транзакцию
	if err := tx.Commit(); err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.Commit()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS - tx.Commit()", cnt)
	return v, nil
}

// CreateDept create new Dept
func (p *PgDb) CreateDept(r *model.Dept) (*model.Dept, error) {
	cnt := "postgers (p *PgDb) CreateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// Начинаем транзакцию
	tx, err := p.Db.Beginx()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.Db.Beginx()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов в транзакционном режиме
	v, err := p.CreateDeptTx(r, tx)
	// =====================================================================
	if err != nil {
		defer tx.Rollback()
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.CreateDeptTx(id, tx)", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// завершаем транзакцию
	if err := tx.Commit(); err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.Commit()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS - tx.Commit()", cnt)
	return v, nil
}

// UpdateEmp Update new emp
func (p *PgDb) UpdateEmp(r *model.Emp) (*model.Emp, error) {
	cnt := "postgers (p *PgDb) UpdateEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// Начинаем транзакцию
	tx, err := p.Db.Beginx()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.Db.Beginx()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов в транзакционном режиме
	v, err := p.UpdateEmpTx(r, tx, true)
	// =====================================================================
	if err != nil {
		defer tx.Rollback()
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.UpdateEmpTx(id, tx)", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// завершаем транзакцию
	if err := tx.Commit(); err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.Commit()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS - tx.Commit()", cnt)
	return v, nil
}

// UpdateDept Update new Dept
func (p *PgDb) UpdateDept(r *model.Dept) (*model.Dept, error) {
	cnt := "postgers (p *PgDb) UpdateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// Начинаем транзакцию
	tx, err := p.Db.Beginx()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.Db.Beginx()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Вызов в транзакционном режиме
	v, err := p.UpdateDeptTx(r, tx)
	// =====================================================================
	if err != nil {
		defer tx.Rollback()
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.UpdateDeptTx(r, tx)", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// завершаем транзакцию
	if err := tx.Commit(); err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.Commit()", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS - tx.Commit()", cnt)
	return v, nil
}
