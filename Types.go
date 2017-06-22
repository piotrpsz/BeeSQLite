//
// File: Types.go
// Project: BeeSQLite
//
// Created by Piotr Pszczółkowski on 22/06/2017
// Copyright 2017 Piotr Pszczółkowski
//

package BeeSQLite

// Row - container for fields in one row of the table
type Row map[string]Field

// Result - container for all rows fetched from database
type Result []Row

// Count - number of rows in result
func (r *Result) Count() int {
	return len(*r)
}

// IsNotEmpty - checks if result contains rows
func (r *Result) IsNotEmpty() bool {
	if r.Count() > 0 {
		return true
	}
	return false
}

// First - returns first row in result
func (r *Result) First() Row {
	if r.IsNotEmpty() {
		return (*r)[0]
	}
	return nil
}

// Last - returns last row in result
func (r *Result) Last() Row {
	n := r.Count()
	if n > 0 {
		return (*r)[n-1]
	}
	return nil
}
