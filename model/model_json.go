// Code generated by Gojay. DO NOT EDIT.

package model

import (
	"database/sql"

	"github.com/francoispqt/gojay"
)

type EmpsPtr []*Emp

func (s *EmpsPtr) UnmarshalJSONArray(dec *gojay.Decoder) error {
	var value = &Emp{}
	if err := dec.Object(value); err != nil {
		return err
	}
	*s = append(*s, value)
	return nil
}

func (s EmpsPtr) MarshalJSONArray(enc *gojay.Encoder) {
	for i := range s {
		enc.Object(s[i])
	}
}

func (s EmpsPtr) IsNil() bool {
	return len(s) == 0
}

// MarshalJSONObject implements MarshalerJSONObject
func (d *Dept) MarshalJSONObject(enc *gojay.Encoder) {
	enc.IntKey("deptNumber", d.Deptno)
	enc.StringKey("deptName", d.Dname)
	enc.SQLNullStringKeyNullEmpty("deptLocation", &d.Loc)
	var empsSlice = EmpsPtr(d.Emps)
	enc.ArrayKey("emps", empsSlice)
}

// IsNil checks if instance is nil
func (d *Dept) IsNil() bool {
	return d == nil
}

// UnmarshalJSONObject implements gojay's UnmarshalerJSONObject
func (d *Dept) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {

	switch k {
	case "deptNumber":
		return dec.Int(&d.Deptno)

	case "deptName":
		return dec.String(&d.Dname)

	case "deptLocation":
		var value = sql.NullString{}
		err := dec.SQLNullString(&value)
		if err == nil {
			d.Loc = value
		}
		return err

	case "emps":
		var aSlice = EmpsPtr{}
		err := dec.Array(&aSlice)
		if err == nil && len(aSlice) > 0 {
			d.Emps = []*Emp(aSlice)
		}
		return err

	}
	return nil
}

// NKeys returns the number of keys to unmarshal
func (d *Dept) NKeys() int { return 4 }

// MarshalJSONObject implements MarshalerJSONObject
func (e *Emp) MarshalJSONObject(enc *gojay.Encoder) {
	enc.IntKey("empNo", e.Empno)
	enc.SQLNullStringKeyNullEmpty("empName", &e.Ename)
	enc.SQLNullStringKeyNullEmpty("job", &e.Job)
	enc.SQLNullInt64KeyOmitEmpty("mgr", &e.Mgr)
	enc.SQLNullStringKeyNullEmpty("hiredate", &e.Hiredate)
	enc.SQLNullInt64KeyNullEmpty("sal", &e.Sal)
	enc.SQLNullInt64KeyNullEmpty("comm", &e.Comm)
	enc.SQLNullInt64KeyNullEmpty("deptNumber", &e.Deptno)
}

// IsNil checks if instance is nil
func (e *Emp) IsNil() bool {
	return e == nil
}

// UnmarshalJSONObject implements gojay's UnmarshalerJSONObject
func (e *Emp) UnmarshalJSONObject(dec *gojay.Decoder, k string) error {

	switch k {
	case "empNo":
		return dec.Int(&e.Empno)

	case "empName":
		var value = sql.NullString{}
		err := dec.SQLNullString(&value)
		if err == nil {
			e.Ename = value
		}
		return err

	case "job":
		var value = sql.NullString{}
		err := dec.SQLNullString(&value)
		if err == nil {
			e.Job = value
		}
		return err

	case "mgr":
		var value = sql.NullInt64{}
		err := dec.SQLNullInt64(&value)
		if err == nil {
			e.Mgr = value
		}
		return err

	case "hiredate":
		var value = sql.NullString{}
		err := dec.SQLNullString(&value)
		if err == nil {
			e.Hiredate = value
		}
		return err

	case "sal":
		var value = sql.NullInt64{}
		err := dec.SQLNullInt64(&value)
		if err == nil {
			e.Sal = value
		}
		return err

	case "comm":
		var value = sql.NullInt64{}
		err := dec.SQLNullInt64(&value)
		if err == nil {
			e.Comm = value
		}
		return err

	case "deptNumber":
		var value = sql.NullInt64{}
		err := dec.SQLNullInt64(&value)
		if err == nil {
			e.Deptno = value
		}
		return err

	}
	return nil
}

// NKeys returns the number of keys to unmarshal
func (e *Emp) NKeys() int { return 8 }
