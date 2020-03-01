package cache

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/pkg/errors"

	"github.com/romapres2010/resttest/model"
)

var testDeptPK []*model.DeptPK

// RandomGetDept return a random row for a given id
// =====================================================================
func (c *Cache) RandomGetDept() (*model.Dept, error) {
	cnt := "- (c *Cache) RandomGetDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)
	var err error

	// Считаем все PK
	if testDeptPK == nil {
		// считаем все PK для dept
		testDeptPK, err = c.GetDeptsPK()
		if err != nil {
			errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptsPK()", cnt)
			log.Printf(errM)
			return nil, errors.WithMessage(err, errM)
		}
	}
	deptsPKlen := len(testDeptPK)

	// случайным образом выбираем PK
	par := testDeptPK[rand.Intn(deptsPKlen)].Deptno
	v, err := c.GetDeptEmps(par)
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - p.GetDeptEmps(par), args = '%v'", cnt, par)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, nil
}

// RandomUpdateDept update a random row
// =====================================================================
func (c *Cache) RandomUpdateDept() (*model.Dept, error) {
	cnt := "cache (c *Cache) UpdateDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - START", cnt)

	v, err := c.deptService.RandomUpdateDept()
	if err != nil {
		errM := fmt.Sprintf("[ERROR] %v - ERROR - c.deptService.RandomUpdateDept()", cnt)
		log.Printf(errM)
		return nil, errors.WithMessage(err, errM)
	}

	//mylog.PrintfDebug("[DEBUG] %v - SUCCESS", cnt)
	return v, err
}
