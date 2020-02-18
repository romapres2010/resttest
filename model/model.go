package model

import (
	//	"gopkg.in/guregu/sql.Nullv3"
	"database/sql"
)

// Dept - основная структура для хранения объекта
type Dept struct {
	Deptno      int            `db:"deptno" json:"deptNumber" validate:"required"`
	Dname       string         `db:"dname" json:"deptName"`
	Loc         sql.NullString `db:"loc" json:"deptLocation, nullempty"`
	Emps        []*Emp         `json:"emps"` // массив указателей на дочерние emp
	EmpsPK      []*EmpPK       `json:"-"`    // массив PK дочерних emp
	IsPutToPool bool           `json:"-"`    // признак, что структура возвращена в pool
}

// Emp - основная структура для хранения объекта
type Emp struct {
	Empno       int            `db:"empno" json:"empNo"`
	Ename       sql.NullString `db:"ename" json:"empName, nullempty"`
	Job         sql.NullString `db:"job" json:"job, nullempty"`
	Mgr         sql.NullInt64  `db:"mgr" json:"mgr, omitempty"`
	Hiredate    sql.NullString `db:"hiredate" json:"hiredate, nullempty"`
	Sal         sql.NullInt64  `db:"sal" json:"sal, nullempty"`
	Comm        sql.NullInt64  `db:"comm" json:"comm, nullempty"`
	Deptno      sql.NullInt64  `db:"deptno" json:"deptNumber, nullempty"`
	Dept        *Dept          `json:"-"` // связь с Dept для кеширования
	IsPutToPool bool           `json:"-"` // признак, что структура возвращена в pool
}

// DeptPK - вспомогательная структура для запроса PK объекта
type DeptPK struct {
	Deptno int `db:"deptno" json:"deptNumber" validate:"required"`
}

// EmpPK - вспомогательная структура для запроса PK объекта
type EmpPK struct {
	Empno int `db:"empno" json:"empNo"`
}

//DeptPKs - колекция указателей на PK
type DeptPKs []*DeptPK

//EmpPKs - колекция указателей на PK
type EmpPKs []*EmpPK

//Emps - колекция указателей на emp
type EmpSlice []*Emp
