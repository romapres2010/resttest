package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/romapres2010/resttest/model"
)

// EmpExistsTx returns TRUE if row for a given id exists. Transaction mode
// =====================================================================
func (p *PgDb) EmpExistsTx(id int, tx *sqlx.Tx) (bool, error) {
	cnt := "postgres (p *PgDb) EmpExistsTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	var stm *sqlx.Stmt
	if tx == nil {
		// Работаем без транзакции
		stm = p.sqlSelectEmpExists
	} else {
		// Помещаем запрос в транзакцию
		stm = tx.Stmtx(p.sqlSelectEmpExists)
	}
	// возможен случай TWO_MANY_ROWS, поэтому определяем срез
	pk := model.GetEmpPK()
	defer model.PutEmpPK(pk) // возвращаем в pool
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
		return false, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	// проверим количество полученных строк
	if len(pk) == 0 {
		return false, nil
	}
	return true, nil
}

// GetEmpTx return a row for a given id. Transaction mode
// =====================================================================
func (p *PgDb) GetEmpTx(id int, tx *sqlx.Tx) (*model.Emp, error) {
	cnt := "postgres (p *PgDb) GetEmpTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	var stm *sqlx.Stmt
	if tx == nil {
		// Работаем без транзакции
		//mylog.PrintfDebug("[DEBUG] %v - NON TX mode", cnt)
		stm = p.sqlSelectEmp
	} else {
		// Помещаем запрос в транзакцию
		stm = tx.Stmtx(p.sqlSelectEmp)
	}

	// Извлечем из кэша новую структуру
	v := model.GetEmp()
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

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return v, nil
}

// GetEmpsByDeptTx returns all emps for a given dept. Transaction mode
// =====================================================================
func (p *PgDb) GetEmpsByDeptTx(id int, tx *sqlx.Tx) ([]*model.Emp, error) {
	cnt := "postgres (p *PgDb) GetEmpsByDeptTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	var stm *sqlx.Stmt
	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	if tx == nil {
		// Работаем без транзакции
		//mylog.PrintfDebug("[DEBUG] %v - NON TX mode", cnt)
		stm = p.sqlSelectEmpsByDept
	} else {
		// Помещаем запрос в транзакцию
		stm = tx.Stmtx(p.sqlSelectEmpsByDept)
	}
	// =====================================================================

	//Выполняем запрос
	// err := stm.Select(&v, id)

	// Открываем курсор
	rows, err := stm.Queryx(id)
	defer rows.Close() // Закрываем в defer, на случай если парсинг будет с panic
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - stm.Select(&v, id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}
	// Определеляем переменную для результата
	//	list := make([]*model.Emp, 0)
	list := model.GetEmpSlice()
	// перебираем все записи
	for rows.Next() {
		// Извлечем из кэша новую структуру
		v := model.GetEmp()
		// разберем поля курсора в поля структуры
		err := rows.StructScan(v)
		// добавляем в коллецию
		list = append(list, v)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - rows.StructScan(r), args = '%+v'", cnt, v)
			log.Printf(errM)
			return nil, errors.Wrap(err, errM)
		}
	}
	rows.Close()

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return list, nil
}

// GetEmpsPKByDeptTx returns all emps PK for a given dept. Transaction mode
// =====================================================================
func (p *PgDb) GetEmpsPKByDeptTx(id int, tx *sqlx.Tx) ([]*model.EmpPK, error) {
	cnt := "postgres (p *PgDb) GetEmpsPKByDeptTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	var stm *sqlx.Stmt
	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	if tx == nil {
		// Работаем без транзакции
		//mylog.PrintfDebug("[DEBUG] %v - NON TX mode", cnt)
		stm = p.sqlSelectEmpsPKByDept
	} else {
		// Помещаем запрос в транзакцию
		stm = tx.Stmtx(p.sqlSelectEmpsPKByDept)
	}
	// Определеляем переменную для результата
	v := model.GetEmpPK()
	//	v := make([]*model.EmpPK, 0)
	// =====================================================================

	//Выполняем запрос
	err := stm.Select(&v, id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - stm.Select(&v, id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.Wrap(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return v, nil
}

// CreateEmpTx create new emp. Transaction mode
// =====================================================================
func (p *PgDb) CreateEmpTx(r *model.Emp, tx *sqlx.Tx, isValidate bool) (*model.Emp, error) {
	cnt := "postgres (p *PgDb) CreateEmpTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// Проверяем определен ли контекст транзакции
	if tx == nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx *sqlx.Tx is NULL", cnt)
		log.Printf(errM)
		return nil, errors.New(errM)
	}

	//=====================================================================
	// Добавить валидацию входной структуры
	//=====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Если запускаем с проверками
	if isValidate {
		{ //Если Dept NULL или НЕ существует, то ошибка
			if !r.Deptno.Valid {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - r.Deptno is NULL", cnt)
				log.Printf(errM)
				return nil, errors.New(errM)
			}
			deptno := int(r.Deptno.Int64)
			// Запрос в транзакции
			deptExists, err := p.DeptExistsTx(deptno, tx)
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - p.DeptExistsTx(deptno, tx), args = '%v'", cnt, deptno)
				log.Printf(errM)
				return nil, errors.WithMessage(err, errM)
			}
			if !deptExists {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - dept '%v' does not exist", cnt, deptno)
				log.Printf(errM)
				return nil, errors.New(errM)
			}
			//mylog.PrintfDebug("[DEBUG] %v - dept %v exists", cnt, deptno)
		}
		{ //Если Emp существует, то игнорируем
			exists, err := p.EmpExistsTx(r.Empno, tx)
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - p.EmpExistsTx(r.Empno, tx), args = '%v'", cnt, r.Empno)
				log.Printf(errM)
				return nil, errors.WithMessage(err, errM)
			}
			// Если запись существует, то ни чего не делаем, возвращем, что пришло на вход
			if exists {
				errM := fmt.Sprintf("[WARN] %v - WARN - emp '%v' already exist - nothing to do", cnt, r.Empno)
				log.Printf(errM)
				return nil, nil
			}
			//mylog.PrintfDebug("[DEBUG] %v - emp '%v' does not exist - can be created", cnt, r.Empno)
		}
	}
	// =====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	stmText := sqlInsertEmpText
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
	// считаем данные обновленные данные - в БД могли быть тригера, которые поменяли данные
	// если запустили без проверок, то можно не возвращать результат - он будет запрошен уровнем выше
	if isValidate {
		v, err := p.GetEmpTx(r.Empno, tx)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetEmpTx(r.Empno, tx), args = '%v'", cnt, r.Empno)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
		r = v
	}
	// =====================================================================

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return r, nil
}

// UpdateEmpTx update emp. Transaction mode
// =====================================================================
func (p *PgDb) UpdateEmpTx(r *model.Emp, tx *sqlx.Tx, isValidate bool) (*model.Emp, error) {
	cnt := "postgres (p *PgDb) UpdateEmpTx" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// Проверяем определен ли контекст транзакции
	if tx == nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - tx *sqlx.Tx is NULL", cnt)
		log.Printf(errM)
		return nil, errors.New(errM)
	}

	//=====================================================================
	// Добавить валидацию входной структуры
	//=====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Если запускаем с проверками
	if isValidate {
		{ // проверки dept
			//Если Dept NULL или НЕ существует, то ошибка
			if r.Deptno.Valid {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - r.Deptno is NULL", cnt)
				log.Printf(errM)
				return nil, errors.New(errM)
			}
			deptno := int(r.Deptno.Int64)
			// Запрос в транзакции
			deptExists, err := p.DeptExistsTx(deptno, tx)
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - p.DeptExistsTx(deptno, tx), args = '%v'", cnt, deptno)
				log.Printf(errM)
				return nil, errors.Wrap(err, errM)
			}
			if !deptExists {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - dept '%v' does not exist", cnt, deptno)
				log.Printf(errM)
				return nil, errors.New(errM)
			}
			//mylog.PrintfDebug("[DEBUG] %v - dept %v exists", cnt, deptno)
		}
		{ // Проверки emp
			//Если строка НЕ существует, то игнорируем
			exists, err := p.EmpExistsTx(r.Empno, tx)
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - p.EmpExistsTx(r.Empno, tx), args = '%v'", cnt, r.Empno)
				log.Printf(errM)
				return nil, errors.Wrap(err, errM)
			}
			if !exists {
				//mylog.PrintfDebug("[DEBUG] %v - ERROR - emp '%v' does not exist - nothing to do", cnt, r.Empno)
				return nil, nil
			}
			//mylog.PrintfDebug("[DEBUG] %v - emp '%v' exist - can be updated", cnt, r.Empno)
		}
	}
	// =====================================================================

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	stmText := sqlUpdateEmpText
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
	// считаем данные обновленные данные - в БД могли быть тригера, которые поменяли данные
	// если запустили без проверок, то можно не возвращать результат - он будет запрошен уровнем выше
	if isValidate {
		v, err := p.GetEmpTx(r.Empno, tx)
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetEmpTx(r.Empno, tx), args = '%v'", cnt, r.Empno)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
		r = v
	}
	// =====================================================================

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)

	return r, nil
}
