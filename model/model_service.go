package model

// DeptService represent basic interface for Dept
// =====================================================================
type DeptService interface {
	DeptExists(id int) (bool, error)
	GetDept(id int) (*Dept, error) // Только родительский объект - Dept
	//GetDeptEmps(id int) (*Dept, error) // Dept и все вложенные Emps
	GetDepts() ([]*Dept, error)
	GetDeptsPK() ([]*DeptPK, error)
	CreateDept(d *Dept) (*Dept, error)
	UpdateDept(d *Dept) (*Dept, error)

	RandomGetDept() (*Dept, error)    // Для целей нагрузочного тестирования
	RandomUpdateDept() (*Dept, error) // Для целей нагрузочного тестирования
	//	DeleteDept(id int) error
	IsUseCache() bool // признак, что релизация с кешированием
}

// EmpService represent basic interface for Emp
// =====================================================================
type EmpService interface {
	EmpExists(id int) (bool, error)
	GetEmp(id int) (*Emp, error)
	GetEmpsByDept(id int) ([]*Emp, error)
	GetEmpsPKByDept(id int) ([]*EmpPK, error)
	//	Emps() ([]*Emp, error)
	CreateEmp(e *Emp) (*Emp, error)
	UpdateEmp(e *Emp) (*Emp, error)
	//	DeleteEmp(id int) error
	IsUseCache() bool // признак, что релизация с кешированием
}
