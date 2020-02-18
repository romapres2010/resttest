package model

import (
	"sync"
)

// deptsPool represent depts pooling
var deptsPool = sync.Pool{
	New: func() interface{} { return new(Dept) },
}

// empsPool represent emps pooling
var empsPool = sync.Pool{
	New: func() interface{} { return new(Emp) },
}

// deptsPKPool represent empsPK pooling
var deptsPKPool = sync.Pool{
	New: func() interface{} {
		v := make([]*DeptPK, 0)
		return DeptPKs(v)
	},
}

// EmpsPKPool represent empsPK pooling
var empsPKPool = sync.Pool{
	New: func() interface{} {
		v := make([]*EmpPK, 0)
		return EmpPKs(v)
	},
}

// empSlicePool represent emps Slice pooling
var empSlicePool = sync.Pool{
	New: func() interface{} {
		v := make([]*Emp, 0)
		return EmpSlice(v)
	},
}

// Reset reset all fields in structure - use for sync.Pool
func (p *Dept) Reset() {
	p.Deptno = 0
	p.Dname = ""
	p.Loc.String = ""
	p.Loc.Valid = false
	EmpSlice(p.Emps).Reset()
	p.Emps = nil
	EmpPKs(p.EmpsPK).Reset()
	p.EmpsPK = nil
}

// Reset reset all fields in structure - use for sync.Pool
func (p *Emp) Reset() {
	p.Empno = 0
	p.Ename.String = ""
	p.Ename.Valid = false
	p.Job.String = ""
	p.Job.Valid = false
	p.Mgr.Int64 = 0
	p.Mgr.Valid = false
	p.Hiredate.String = ""
	p.Hiredate.Valid = false
	p.Sal.Int64 = 0
	p.Sal.Valid = false
	p.Comm.Int64 = 0
	p.Comm.Valid = false
	p.Deptno.Int64 = 0
	p.Deptno.Valid = false
	p.Dept = nil
}

// Reset reset all fields in structure - use for sync.Pool
func (p EmpPKs) Reset() {
	for i := range p {
		p[i].Empno = 0
	}
}

// Reset reset all fields in structure - use for sync.Pool
func (p DeptPKs) Reset() {
	for i := range p {
		p[i].Deptno = 0
	}
}

// Reset reset all fields in structure - use for sync.Pool
func (p EmpSlice) Reset() {
	for i := range p {
		p[i].Reset()
		PutEmp(p[i])
	}
}

// GetDept allocates a new struct or grabs a cached one
func GetDept() *Dept {
	//cnt := "model sync.Pool.GetDept" // Имя текущего метода для логирования

	p := deptsPool.Get().(*Dept)
	p.Reset()
	//mylog.PrintfDebug("[DEBUG] %v - p = '%p'", cnt, p)
	return p
}

// PutDept return struct to cache
func PutDept(p *Dept, isCascad bool) {
	//cnt := "model sync.Pool.PutDept" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - p = '%p', isCascad = '%v'", cnt, p, isCascad)

	if p != nil {
		p.IsPutToPool = true // установим признак, что структура была возвращена в pool - для анализа при последующем получении
		PutEmpSlice(p.Emps, isCascad)
		p.Emps = nil
		PutEmpPK(p.EmpsPK) // Вернем в кэш
		p.EmpsPK = nil
		deptsPool.Put(p)
	}
}

// GetEmp allocates a new struct or grabs a cached one
func GetEmp() *Emp {
	//cnt := "model sync.Pool.GetEmp" // Имя текущего метода для логирования

	p := empsPool.Get().(*Emp)
	p.Reset()

	//mylog.PrintfDebug("[DEBUG] %v - p = '%p'", cnt, p)
	return p
}

// PutEmp return struct to cache
func PutEmp(p *Emp) {
	//cnt := "model sync.Pool.PutEmp" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v - p = '%p'", cnt, p)

	if p != nil {
		p.IsPutToPool = true // установим признак, что структура была возвращена в pool - для анализа при последующем получении
		empsPool.Put(p)
	}
}

// GetDeptPK allocates a new struct or grabs a cached one
func GetDeptPK() DeptPKs {
	//cnt := "model sync.Pool.GetDeptPK" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v", cnt)

	p := deptsPKPool.Get().(DeptPKs)
	p.Reset()
	return p
}

// PutDeptPK return struct to cache
func PutDeptPK(p DeptPKs) {
	//cnt := "model sync.Pool.PutDeptPK" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v", cnt)

	p = p[:0] // сброс
	deptsPKPool.Put(p)
}

// GetEmpPK allocates a new struct or grabs a cached one
func GetEmpPK() EmpPKs {
	//cnt := "model sync.Pool.GetEmpPK" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v", cnt)

	p := empsPKPool.Get().(EmpPKs)
	p.Reset()
	return p
}

// PutEmpPK return struct to cache
func PutEmpPK(p EmpPKs) {
	//cnt := "model sync.Pool.PutEmpPK" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v", cnt)

	p = p[:0] // сброс
	empsPKPool.Put(p)
}

// GetEmpSlice allocates a new struct or grabs a cached one
func GetEmpSlice() EmpSlice {
	//cnt := "model sync.Pool.GetEmpSlice" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v", cnt)

	p := empSlicePool.Get().(EmpSlice)
	p.Reset()
	return p
}

// PutEmpSlice return struct to cache
func PutEmpSlice(p EmpSlice, isCascad bool) {
	//cnt := "model sync.Pool.PutEmpSlice" // Имя текущего метода для логирования
	//mylog.PrintfDebug("[DEBUG] %v isCascad = '%v'", cnt, isCascad)

	if p != nil {
		if isCascad {
			for i := range p {
				PutEmp(p[i])
			}
		}

		p = p[:0] // сброс указателя слайса
		empSlicePool.Put(p)
	}
}
