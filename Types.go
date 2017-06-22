// MIT License
//
// Copyright (c) 2017 Piotr Pszczółkowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// File: Types.go
// Project: BeeSQLite

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
