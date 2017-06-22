//
// File: SQLiteManager.go
// Project: SQLite for Go
//
// Created by Piotr Pszczółkowski on 21/06/2017
// Copyright 2017 Piotr Pszczółkowski
//

package BeeSQLite

/*
#include <stdlib.h>
#include <sqlite3.h>

#cgo LDFLAGS: -lsqlite3

int bind_text(sqlite3_stmt *stmt, int index, const char* txt) {
	return sqlite3_bind_text(stmt, index, txt, -1, SQLITE_TRANSIENT);
}

const char* column_text(sqlite3_stmt *stmt, int index) {
	return (const char *)sqlite3_column_text(stmt, index);
}
*/
import "C"
import "unsafe"

const (
	Ok         = iota // Successful result
	Error             // SQL error or missing database
	Internal          // Internal logic error in SQLite
	Perm              // Access permission denied
	Abort             // Callback routine requested an abort
	Busy              // The database file is locked
	Locked            // A table in the database is locked
	NoMem             // A malloc() failed
	ReadOnly          // Attempt to write a readonly database
	Interrupt         // Operation terminated by sqlite3_interrupt()
	IoErr             // Some kind of disk I/O error occurred
	Corrupt           // The database disk image is malformed
	NotFound          // NOT USED. Table or record not found
	Full              // Insertion failed because database is full
	CantOpen          // Unable to open the database file
	Protocol          // NOT USED. Database lock protocol error
	Empty             // Database is empty
	Schema            // The database schema changed
	TooBig            // String or BLOB exceeds size limit
	Constraint        // Abort due to constraint violation
	Mismatch          // Data type mismatch
	Misuse            // Library used incorrectly
	NoLfs             // Uses OS features not supported on host
	Auth              // Authorization denied
	Format            // Auxiliary database format error
	Range             // 2nd parameter to sqlite3_bind out of range
	NotADb            // File opened that is not a database file
	StatusRow  = 100  // sqlite3_step() has another row ready
	StatusDone = 101  // sqlite3_step() has finished executing
)

type Statement struct {
	stmt *C.sqlite3_stmt
}

func (s *Statement) step() int {
	return int(C.sqlite3_step(s.stmt))
}

func (s *Statement) finalize() int {
	return int(C.sqlite3_finalize(s.stmt))
}

func (s *Statement) reset() int {
	retv := C.sqlite3_reset(s.stmt)
	if retv == C.SQLITE_OK {
		return int(C.sqlite3_clear_bindings(s.stmt))
	}
	return int(retv)
}

func (s *Statement) columnType(idx int) ValueType {
	ct := C.sqlite3_column_type(s.stmt, C.int(idx))
	switch ct {
	case C.SQLITE_INTEGER:
		return Integer
	case C.SQLITE_FLOAT:
		return Float
	case C.SQLITE_TEXT:
		return Text
	case C.SQLITE_BLOB:
		return Blob
	}
	return Null
}

func (s *Statement) columnCount() int {
	return int(C.sqlite3_column_count(s.stmt))
}

func (s *Statement) columnIndex(columnName string) int {
	cstring := C.CString(columnName)
	defer C.free(unsafe.Pointer(cstring))
	return int(C.sqlite3_bind_parameter_index(s.stmt, cstring))
}

func (s *Statement) columnName(columnIndex int) string {
	return C.GoString(C.sqlite3_column_name(s.stmt, C.int(columnIndex)))
}

func (s *Statement) selectQuery(query string) Result {
	var result Result
	n := s.columnCount()
	if n > 0 {
		for s.step() == StatusRow {
			var row = Row{}
			for i := 0; i < n; i++ {
				name := s.columnName(i)
				var field = Field{Name: name}
				switch s.columnType(i) {
				case Null:
					field.SetValue(nil)
				case Integer:
					field.SetValue(s.fetchInt(i))
				case Float:
					field.SetValue(s.fetchFloat(i))
				case Text:
					field.SetValue(s.fetchString(i))
				case Blob:
					field.SetValue(s.fetchBlob(i))
				}
				row[name] = field
			}
			result = append(result, row)
		}
		s.finalize()
	}
	return result
}

func (s *Statement) bindFields(fields []Field) bool {
	for _, field := range fields {
		index := s.columnIndex(field.BindName())
		switch field.valueType {
		case Null:
			s.bindNull(index)
		case Integer:
			s.bindInt(index, field.Int())
		case Float:
			s.bindFloat(index, field.Float())
		case Text:
			s.bindString(index, field.String())
		case Blob:
			s.bindBlob(index, field.Blob())
		}

	}

	state := s.step()
	ok1 := (state == Ok || state == StatusDone)
	ok2 := (s.finalize() == Ok)
	if ok1 && ok2 {
		return true
	}
	return false
}

//###################################################################
//#                                                                 #
//#                       S E T T E R S                             #
//#                                                                 #
//###################################################################

func (s *Statement) bindNull(columnIndex int) int {
	return int(C.sqlite3_bind_null(s.stmt, C.int(columnIndex)))
}

func (s *Statement) bindInt(columnIndex int, value int) int {
	return int(C.sqlite3_bind_int64(s.stmt, C.int(columnIndex), C.sqlite3_int64(value)))
}

func (s *Statement) bindFloat(columnIndex int, value float64) int {
	return int(C.sqlite3_bind_double(s.stmt, C.int(columnIndex), C.double(value)))
}

func (s *Statement) bindString(columnIndex int, value string) int {
	cstring := C.CString(value)
	defer C.free(unsafe.Pointer(cstring))
	return int(C.bind_text(s.stmt, C.int(columnIndex), cstring))
}

func (s *Statement) bindBlob(columnIndex int, value []byte) int {
	return int(C.sqlite3_bind_blob(s.stmt, C.int(columnIndex), unsafe.Pointer(&value[0]), C.int(len(value)), nil))
}

//###################################################################
//#                                                                 #
//#                       G E T T E R S                             #
//#                                                                 #
//###################################################################

func (s *Statement) fetchInt(columnIndex int) int {
	return int(C.sqlite3_column_int64(s.stmt, C.int(columnIndex)))
}

func (s *Statement) fetchFloat(columnIndex int) float64 {
	return float64(C.sqlite3_column_double(s.stmt, C.int(columnIndex)))
}

func (s *Statement) fetchString(columnIndex int) string {
	return C.GoString(C.column_text(s.stmt, C.int(columnIndex)))
}

func (s *Statement) fetchBlob(columnIndex int) []byte {
	n := C.int(C.sqlite3_column_bytes(s.stmt, C.int(columnIndex)))
	ptr := C.sqlite3_column_blob(s.stmt, C.int(columnIndex))
	return C.GoBytes(ptr, n)
}
