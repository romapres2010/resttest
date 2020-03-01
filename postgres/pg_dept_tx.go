package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/romapres2010/resttest/model"
)

// DeptExistsTx returns TRUE if row for a given id exists. Transaction mode
// =====================================================================
func (p *PgDb) DeptExistsTx(id int, tx *sqlx.Tx) (bool, error) {
	cnt := "postgres (p *PgDb) DeptExistsTx"
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	var stm *sqlx.Stmt
	if tx == nil {
		// Работаем без транзакции
		//mylog.PrintfDebug("[DEBUG] %v - NON TX mode", cnt)
		stm = p.sqlSelectDeptExists
	} else {
		// Помещаем запрос в транзакцию
		stm = tx.Stmtx(p.sqlSelectDeptExists)
	}
	// возможен случай TWO_MANY_ROWS, поэтому определяем срез
	pk := model.GetDeptPK()
	defer model.PutDeptPK(pk) // возвращаем в pool
	//	pk := make([]*model.DeptPK, 0)
	// =====================================================================

	//Выполняем запрос
	err := stm.Select(&pk, id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - stm.Select(&pk, id), args = '%v'", cnt, id)
		log.Printf(errM)
		return false, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - get %v rows", cnt, len(pk))

	// Маловероятная ситуация - только если проблема с БД
	if len(pk) > 1 {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - TWO_MANY_ROWS", cnt)
		log.Printf(errM)
		return false, errors.New(errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	// проверим количество полученных строк
	if len(pk) == 0 {
		return false, nil
	}
	return true, nil
}

// GetDeptTx return a row for a given id without Emps. Transaction mode
// =====================================================================
func (p *PgDb) GetDeptTx(id int, tx *sqlx.Tx) (*model.Dept, error) {
	cnt := "postgres (p *PgDb) GetDeptTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	var stm *sqlx.Stmt
	if tx == nil {
		// Работаем без транзакции
		//mylog.PrintfDebug("[DEBUG] %v - NON TX mode", cnt)
		stm = p.sqlSelectDept
	} else {
		// Помещаем запрос в транзакцию
		stm = tx.Stmtx(p.sqlSelectDept)
	}

	// Извлечем из кэша новую структуру
	v := model.GetDept()
	// =====================================================================

	//Выполняем запрос - ожидаем одну строку
	err := stm.Get(v, id)
	if err != nil {
		// проверим NO_DATA_FOUND
		if err == sql.ErrNoRows {
			errM := fmt.Sprintf("[DEBUG] %v - stm.Get(&v, id) - NO_DATA_FOUND, param: '%v', SQL err: '%+v'", cnt, id, err)
			log.Printf(errM)
			return nil, nil
		}
		// Другие ошибки
		errM := fmt.Sprintf("[ERROR] %v - ERROR - stm.Get(&v, id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Считываем подчиненные записи только PK
	v.EmpsPK, err = p.GetEmpsPKByDeptTx(v.Deptno, tx)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetEmpsPKByDept(v.Deptno, tx), args = '%v'", cnt, v.Deptno)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return v, nil
}

// GetDeptEmpsTx return a row for a given id and all Emps. Transaction mode
// =====================================================================
func (p *PgDb) GetDeptEmpsTx(id int, tx *sqlx.Tx) (*model.Dept, error) {
	cnt := "postgres (p *PgDb) GetDeptEmpsTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Запрашиваем родительскую запись
	v, err := p.GetDeptTx(id, tx)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptTx(id, tx), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Считываем подчиненные записи все поля
	v.Emps, err = p.GetEmpsByDeptTx(v.Deptno, tx)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetEmpsByDeptTx(v.Deptno, tx), args = '%v'", cnt, v.Deptno)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return v, nil
}

// GetDeptsTx returns all depts. Transaction mode
// =====================================================================
func (p *PgDb) GetDeptsTx(tx *sqlx.Tx) ([]*model.Dept, error) {
	cnt := "postgres (p *PgDb) GetDeptsTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	var stm *sqlx.Stmt
	if tx == nil {
		// Работаем без транзакции
		//mylog.PrintfDebug("[DEBUG] %v - NON TX mode", cnt)
		stm = p.sqlSelectDepts
	} else {
		// Помещаем запрос в транзакцию
		stm = tx.Stmtx(p.sqlSelectDepts)
	}
	// Определеляем переменную для результата
	v := make([]*model.Dept, 0)

	//Выполняем запрос
	err := stm.Select(&v)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - stm.Select(&v)", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return v, nil
}

// GetDeptsPKTx returns Primary Key (PK) for all depts. Transaction mode
// =====================================================================
func (p *PgDb) GetDeptsPKTx(tx *sqlx.Tx) ([]*model.DeptPK, error) {
	cnt := "postgres (p *PgDb) GetDeptsPKTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	var stm *sqlx.Stmt
	if tx == nil {
		// Работаем без транзакции
		//mylog.PrintfDebug("[DEBUG] %v - NON TX mode", cnt)
		stm = p.sqlSelectDeptsPK
	} else {
		// Помещаем запрос в транзакцию
		stm = tx.Stmtx(p.sqlSelectDeptsPK)
	}
	// Определеляем переменную для результата
	v := model.GetDeptPK()
	//v := make([]*model.DeptPK, 0)

	//Выполняем запрос
	err := stm.Select(&v)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - stm.Select(&v)", cnt)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return v, nil
}

// CreateDeptTx create new dept. Transaction mode
// =====================================================================
func (p *PgDb) CreateDeptTx(r *model.Dept, tx *sqlx.Tx) (*model.Dept, error) {
	cnt := "postgres (p *PgDb) CreateDeptTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// Проверяем определен ли контекст транзакции
	if tx == nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx *sqlx.Tx is NULL", cnt)
		log.Printf(errM)
		return nil, errors.New(errM)
	}

	//=====================================================================
	// Добавить валидацию входной структуры
	// Добавить валидацию, что подчиненные ресурсы имеют правильные ссылки на родительский
	//=====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	{ //Если строка существует, то ни чего не делаем
		deptExists, err := p.DeptExistsTx(r.Deptno, tx)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - p.DeptExistsTx(r.deptno, tx), args = '%v'", cnt, r.Deptno)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
		// Если запись существует, то ни чего не делаем, возвращем, что пришло на вход
		if deptExists {
			errM := fmt.Sprintf("[DEBUG] %v - dept '%v' already exist - nothing to do", cnt, r.Deptno)
			log.Printf(errM)
			return nil, nil
		}
		//mylog.PrintfDebug("[DEBUG] %v - dept '%v' does not exist - can be created", cnt, r.Deptno)
	}
	// =====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	stmText := sqlInsertDeptText
	// =====================================================================

	//Выполняем команду
	res, err := tx.NamedExec(stmText, r)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.NamedExec(stmText, r), args = '%+v'", cnt, r)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	{ // Необязательная часть - можно удалить в последствии
		// Проверим количество обработанных строк
		rowCount, err := res.RowsAffected()
		_ = rowCount
		//mylog.PrintfDebug("[DEBUG] %v -- process %v rows", cnt, rowCount)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - res.RowsAffected()", cnt)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Обработаем emps
	if r.Emps != nil {
		for _, e := range r.Emps {
			// вызываем в режиме без валидации и без возварта результата
			_, err = p.CreateEmpTx(e, tx, false)
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - p.CreateEmpTx(e, tx), args = '%v'", cnt, e)
				log.Printf(errM)
				return nil, errors.WithMessage(err, errM)
			}
		}
	}
	// =====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// считаем данные обновленные данные - в БД могли быть тригера, которые поменяли данные
	v, err := p.GetDeptEmpsTx(r.Deptno, tx)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptTx(r.Deptno, tx), args = '%v'", cnt, r.Deptno)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// =====================================================================

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return v, nil
}

// UpdateDeptTx update dept. Transaction mode
// =====================================================================
func (p *PgDb) UpdateDeptTx(r *model.Dept, tx *sqlx.Tx) (*model.Dept, error) {
	cnt := "postgres (p *PgDb) UpdateDeptTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// Проверяем определен ли контекст транзакции
	if tx == nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx *sqlx.Tx is NULL", cnt)
		log.Printf(errM)
		return nil, errors.New(errM)
	}

	//=====================================================================
	// Добавить валидацию входной структуры
	// Добавить валидацию, что подчиненные ресурсы имеют правильные ссылки на родительский
	//=====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	{ //Если строка не существует, то ни чего не делаем
		deptExists, err := p.DeptExistsTx(r.Deptno, tx)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - p.DeptExistsTx(r.deptno, tx), args = '%v'", cnt, r.Deptno)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
		if !deptExists {
			errM := fmt.Sprintf("[DEBUG] %v - dept '%v' does not exist - nothing to do", cnt, r.Deptno)
			log.Printf(errM)
			return nil, nil
		}
		//mylog.PrintfDebug("[DEBUG] %v - dept '%v' exist - can be updated", cnt, r.Deptno)
	}
	// =====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	stmText := sqlUpdateDeptText
	// =====================================================================

	//Выполняем команду
	res, err := tx.NamedExec(stmText, r)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx.NamedExec(stmText, r), args = '%v'", cnt, r)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	{ // Необязательная часть - можно удалить в последствии
		// Проверим количество обработанных строк
		rowCount, err := res.RowsAffected()
		_ = rowCount
		//mylog.PrintfDebug("[DEBUG] %v -- process %v rows", cnt, rowCount)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - res.RowsAffected()", cnt)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Обработаем emps
	if r.Emps != nil {
		for _, e := range r.Emps {
			// вызываем в режиме без валидации и без возварта результата
			_, err = p.UpdateEmpTx(e, tx, false)
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - p.UpdateEmpTx(e, tx), args = '%v'", cnt, e)
				log.Printf(errM)
				return nil, errors.WithMessage(err, errM)
			}
		}
	}

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// считаем данные обновленные данные - в БД могли быть тригера, которые поменяли данные
	v, err := p.GetDeptEmpsTx(r.Deptno, tx)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptTx(r.Deptno, tx), args = '%v'", cnt, r.Deptno)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// =====================================================================

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return v, nil
}
