package cache

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/romapres2010/restdept/model"
)

// IsUseCache returns TRUE for cache
// =====================================================================
func (c *Cache) IsUseCache() bool {
	return true
}

// GetDept return a row for a given id. Cach mode
// =====================================================================
func (c *Cache) GetDept(id int) (*model.Dept, error) {
	cnt := "cache (c *Cache) GetDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// Ищем запись в кэше
	if v, ok := c.depts.get(id); ok {
		//mylog.PrintfDebug("[DEBUG] %v - c.depts.get(id) - presented in cache", cnt)
		return v, nil
	}
	// не нашли к кэше
	//mylog.PrintfDebug("[DEBUG] %v - c.depts.get(id) - NOT presented in cache. Get from DB", cnt)

	// считываем данные из БД
	v, err := c.deptService.GetDept(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.service.GetDept(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// Данные не найдены - нормальная ситация
	if v == nil {
		errM := fmt.Sprintf("[DEBUG] %v - NO_DATA_FOUND - c.service.GetDept(id)", cnt)
		log.Printf(errM)
		return nil, nil
	}

	// Запишем в кэш с блокировкой на запись
	c.depts.set(id, v)

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// DeptExists returns TRUE if row for a given id exists
// =====================================================================
func (c *Cache) DeptExists(id int) (bool, error) {
	cnt := "cache (c *Cache) DeptExists" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	v, err := c.deptService.DeptExists(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.deptService.DeptExists(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return false, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, err
}

// EmpExists returns TRUE if row for a given id exists
// =====================================================================
func (c *Cache) EmpExists(id int) (bool, error) {
	cnt := "cache (c *Cache) EmpExists" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	v, err := c.empService.EmpExists(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.empService.EmpExists(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return false, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, err
}

// GetDeptEmps return a row for a given id and all Emps. Cach mode
// =====================================================================
func (c *Cache) GetDeptEmps(id int) (*model.Dept, error) {
	cnt := "cache (c *Cache) GetDeptEmps" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// Запрашиваем родительскую запись. Данные одновременно записываются в кэш
	v, err := c.GetDept(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.GetDept(id, tx), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// Считываем подчиненные записи. Данные одновременно записываются в кэш
	v.Emps, err = c.GetEmpsByDept(v.Deptno)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.GetEmpsByDept(v.Deptno), args = '%v'", cnt, v.Deptno)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// GetEmpsByDept returns all emps for a given dept. Cach mode
// =====================================================================
func (c *Cache) GetEmpsByDept(id int) ([]*model.Emp, error) {
	cnt := "cache (c *Cache) GetEmpsByDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// =====================================================================
	// Переменная часть кода
	// =====================================================================
	// Запрашиваем родительскую запись. Данные одновременно записываются в кэш
	v, err := c.GetDept(id)
	// =====================================================================
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.GetDept(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// ранее считанные empPK для данного dept
	pks := v.EmpsPK

	// при повторном запросе emps уже заполнен
	list := v.Emps

	if len(pks) == 0 {
		log.Printf(fmt.Sprintf("[DEBUG] %v - len(pks) == 0", cnt))
		return nil, nil
	}

	// Если есть подчиненные записи
	if pks != nil {
		if list != nil {
			list = list[:0] // сбросим слайс, чтобы пересобрать его заново
		} else {
			list = model.GetEmpSlice() // возьмем новый из pool
			v.Emps = list
		}

		// Обработаем подчиненные записи по одной
		for _, pk := range pks {
			emp, err := c.GetEmp(pk.Empno)
			emp.Dept = v // для обратной связи в Кэше
			if err != nil {
				errM := fmt.Sprintf("[ERROR] %v - ERROR - c.GetEmp(pk.Empno), args = '%v'", cnt, pk.Empno)
				log.Printf(errM)
				return nil, errors.WithMessage(err, errM)
			}
			list = append(list, emp)
		}
	}
	// =====================================================================

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return list, nil
}

// GetEmpsPKByDept returns all emps PK for a given dept
// =====================================================================
func (c *Cache) GetEmpsPKByDept(id int) ([]*model.EmpPK, error) {
	cnt := "cache (c *Cache) GetEmpsPKByDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	v, err := c.empService.GetEmpsPKByDept(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.empService.GetEmpsPKByDept(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, err
}

// GetDeptsPK returns Primary Key (PK) for all depts
// =====================================================================
func (c *Cache) GetDeptsPK() ([]*model.DeptPK, error) {
	cnt := "cache (c *Cache) GetDeptsPK" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	v, err := c.deptService.GetDeptsPK()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.deptService.GetDeptsPK()", cnt)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, err
}

// GetDepts returns all depts
// =====================================================================
func (c *Cache) GetDepts() ([]*model.Dept, error) {
	cnt := "cache (c *Cache) GetDepts" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	v, err := c.deptService.GetDepts()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.deptService.GetDepts()", cnt)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, err
}

// GetEmp return a row for a given id. Cach mode
// =====================================================================
func (c *Cache) GetEmp(id int) (*model.Emp, error) {
	cnt := "cache (c *Cache) GetEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, id)

	// Найдена запись в кэше
	if v, ok := c.emps.get(id); ok {
		//mylog.PrintfDebug("[DEBUG] %v - c.emps.get(id) - presented in cache", cnt)
		return v, nil
	}
	// не нашли к кэше
	//mylog.PrintfDebug("[DEBUG] %v - c.emps.get(id) -  NOT presented in cache. Get from DB", cnt)

	// считываем данные из БД
	v, err := c.empService.GetEmp(id)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.service.GetDept(id), args = '%v'", cnt, id)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}
	// Данные не найдены - нормальная ситация
	if v == nil {
		errM := fmt.Sprintf("[DEBUG] %v - NO_DATA_FOUND - c.service.GetDept(id)", cnt)
		log.Printf(errM)
		return nil, nil
	}

	// Запишем в кэш с блокировкой на запись
	c.emps.set(id, v)

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, err
}

// CreateEmp create new emp. Cach mode
// =====================================================================
func (c *Cache) CreateEmp(r *model.Emp) (*model.Emp, error) {
	cnt := "cache (c *Cache) CreateEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// вызываем обработку из БД
	v, err := c.empService.CreateEmp(r)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.empService.CreateEmp(r), args = '%+v'", cnt, r)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// Запишем в кэш с блокировкой на запись
	if v != nil {
		c.emps.set(v.Empno, v)
		//mylog.PrintfDebug("[DEBUG] %v - c.emps.set(v.Empno, v), args = '%+v'", cnt, v.Empno)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// CreateDept create new Dept. Cach mode
// =====================================================================
func (c *Cache) CreateDept(r *model.Dept) (*model.Dept, error) {
	cnt := "cache (c *Cache) CreateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// вызываем обработку из БД
	v, err := c.deptService.CreateDept(r)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.deptService.CreateDept(r), args = '%+v'", cnt, r)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// Запишем в кэш с блокировкой на запись
	if v != nil {
		c.depts.set(v.Deptno, v)
		//mylog.PrintfDebug("[DEBUG] %v - c.depts.set(v.Deptno, v), args = '%+v'", cnt, v.Deptno)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// UpdateEmp Update emp. Cach mode
// =====================================================================
func (c *Cache) UpdateEmp(r *model.Emp) (*model.Emp, error) {
	cnt := "cache (c *Cache) UpdateEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// удалим из кеша
	c.emps.outdate(r.Empno)
	//mylog.PrintfDebug("[DEBUG] %v - c.emps.outdate(r.Empno), args = '%+v'", cnt, r.Empno)

	// вызываем обработку из БД
	v, err := c.empService.UpdateEmp(r)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.empService.UpdateEmp(r), args = '%+v'", cnt, r)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// Запишем в кэш с блокировкой на запись
	if v != nil {
		c.emps.set(v.Empno, v)
		//mylog.PrintfDebug("[DEBUG] %v - c.emps.set(v.Empno, v), args = '%+v'", cnt, v.Empno)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// UpdateDept Update Dept. Cach mode
// =====================================================================
func (c *Cache) UpdateDept(r *model.Dept) (*model.Dept, error) {
	cnt := "cache (c *Cache) UpdateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START, param: '%+v'", cnt, r)

	// удалим из кеша
	c.depts.outdate(r.Deptno)
	//mylog.PrintfDebug("[DEBUG] %v - c.depts.outdate(r.Deptno), args = '%+v'", cnt, r.Deptno)

	// вызываем обработку из БД
	v, err := c.deptService.UpdateDept(r)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.deptService.UpdateDept(r), args = '%+v'", cnt, r)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	// Запишем в кэш с блокировкой на запись
	if v != nil {
		c.depts.set(v.Deptno, v)
		//mylog.PrintfDebug("[DEBUG] %v - c.depts.set(v.Deptno, v), args = '%+v'", cnt, v.Deptno)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}
