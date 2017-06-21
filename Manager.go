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

*/
import "C"
import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

var header = [...]byte{0x53, 0x51, 0x4c, 0x69, 0x74, 0x65, 0x20, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x20, 0x33, 0x00}

type Row map[string]Field

type SQLite struct {
	db        *C.sqlite3
	statement Statement
}

// Version - returns version of SQLite
func (s *SQLite) Version() string {
	return C.GoString(C.sqlite3_libversion())
}

// Prepare - prepares query
func (s *SQLite) prepare(query string) int {
	cstring := C.CString(query)
	defer C.free(unsafe.Pointer(cstring))
	return int(C.sqlite3_prepare_v2(s.db, cstring, -1, &s.statement.stmt, nil))
}

// ErrorCode - returns code of last error
func (s *SQLite) ErrorCode() int {
	return int(C.sqlite3_errcode(s.db))
}

// ErrorString - returns description of last error
func (s *SQLite) ErrorString() string {
	return C.GoString(C.sqlite3_errmsg(s.db))
}

// Remove - removes database file from disk
func (s *SQLite) Remove(fpath string) bool {
	s.Close()
	err := os.Remove(fpath)
	if err != nil {
		fmt.Println(err)
	}
	return true
}

// Open - opens database (read & write)
func (s *SQLite) Open(fpath string) bool {
	if s.db != nil {
		fmt.Println("Database is already opened")
		return false
	}

	if !databaseExists(fpath) {
		fmt.Println("Database not exists")
		return false
	}

	cstring := C.CString(fpath)
	defer C.free(unsafe.Pointer(cstring))

	C.sqlite3_initialize()
	retv := C.sqlite3_open_v2(cstring, &s.db, C.SQLITE_OPEN_READWRITE, nil)
	if retv == C.SQLITE_OK {
		if s.applyPragmas() {
			return true
		}
	}
	return false
}

// Create - creates new database
func (s *SQLite) Create(fpath string) bool {
	if s.db != nil {
		fmt.Println("Database is already opened")
		return false
	}

	if databaseExists(fpath) {
		fmt.Println("Database exists")
		return false
	}

	cstring := C.CString(fpath)
	defer C.free(unsafe.Pointer(cstring))

	C.sqlite3_initialize()
	retv := C.sqlite3_open_v2(cstring, &s.db, C.SQLITE_OPEN_READWRITE|C.SQLITE_OPEN_CREATE, nil)
	if retv == C.SQLITE_OK {
		if s.applyPragmas() {
			return true
		}
	}
	s.checkError()
	return false

}

// Close - closes database and deinits library
func (s *SQLite) Close() {
	if s.db == nil {
		return
	}
	if C.sqlite3_close(s.db) == C.SQLITE_OK {
		C.sqlite3_shutdown
		s.db = nil
	}
}

// ExecQuery - executes query
func (s *SQLite) ExecQuery(query string) bool {
	cquery := C.CString(query)
	defer C.free(unsafe.Pointer(cquery))

	if C.sqlite3_exec(s.db, cquery, nil, nil, nil) == C.SQLITE_OK {
		return true
	}
	s.checkError()
	return false
}

func (s *SQLite) BeginTransaction() bool {
	return s.ExecQuery("BEGIN IMMEDIATE TRANSACTION")
}

func (s *SQLite) CommitTransaction() bool {
	return s.ExecQuery("COMMIT TRANSACTION")
}

func (s *SQLite) RollbackTransaction() bool {
	return s.ExecQuery("ROLLBACK TRANSACTION")
}

func (s *SQLite) EndTransaction(success bool) bool {
	if success {
		return s.CommitTransaction()
	}
	return s.RollbackTransaction()
}

func (s *SQLite) lastInsertRowID() int {
	return int(C.sqlite3_last_insert_rowid(s.db))
}

// Insert - inserts fields to table
func (s *SQLite) Insert(table string, fields []Field) (int, bool) {
	if len(fields) == 0 {
		return 0, false
	}

	names := ""
	binds := ""
	for _, field := range fields[:len(fields)-1] {
		names += fmt.Sprintf("%s,", field.Name)
		binds += fmt.Sprintf("%s,", field.BindName())
	}
	field := fields[len(fields)-1]
	names += field.Name
	binds += field.BindName()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, names, binds)
	if s.prepare(query) == Ok {
		if s.statement.bindFields(fields) {
			rowid := s.lastInsertRowID()
			return rowid, true
		}
	}
	s.checkError()
	return 0, false
}

// Update - updates record content
func (s *SQLite) Update(table string, fields []Field) bool {
	fmt.Println("SQLiteManager.Update")
	if len(fields) == 0 {
		return false
	}

	n := len(fields) - 1
	tokens := ""

	for i := 0; i < n; i++ {
		field := fields[i]
		tokens += fmt.Sprintf("%s=%s,", field.Name, field.BindName())
	}
	field := fields[n]
	tokens += fmt.Sprintf("%s=%s", field.Name, field.BindName())
	whereClause := fmt.Sprintf("%s=%d", fields[0].Name, fields[0].Int())
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s", table, tokens, whereClause)
	if s.prepare(query) == Ok {
		if s.statement.bindFields(fields) {
			return true
		}
	}
	s.checkError()
	return false
}

// Select - selects query
func (s *SQLite) Select(query string) ([]Row, bool) {
	if s.prepare(query) == Ok {
		result := s.statement.selectQuery(query)
		if s.ErrorCode() == Ok {
			return result, true
		}
	}
	s.checkError()
	return nil, false
}

//###################################################################
//#                                                                 #
//#                          P R I V A T E                          #
//#                                                                 #
//###################################################################

func (s *SQLite) checkError() {
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("SQLite: %s (%d): %s (%d)", fn, line, s.ErrorString(), s.ErrorCode())
}

func (s *SQLite) applyPragmas() bool {
	return s.ExecQuery("PRAGMA foreign_keys = ON")
}

func databaseExists(fpath string) bool {
	file, err := os.Open(fpath)
	defer file.Close()

	if err != nil {
		// fmt.Println(err)
		return false
	}

	n := len(header)
	data := make([]byte, n)
	count, err := file.Read(data)
	if err != nil {
		// fmt.Println(err)
		return false
	}
	if count != n {
		fmt.Println("Invalid bytes count")
		return false
	}

	for i := 0; i < n; i++ {
		if data[i] != header[i] {
			return false
		}
	}
	return true
}
