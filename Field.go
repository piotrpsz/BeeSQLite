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

// File: Field.go
// Project: BeeSQLite

package BeeSQLite

import (
	"log"
)

type ValueType int

// Types of value in field
const (
	Null ValueType = iota
	Integer
	Float
	Text
	Blob
)

// Field - information about field
type Field struct {
	Name      string
	valueType ValueType
	value     interface{}
}

// BindName - computed name used in binding
func (f *Field) BindName() string {
	return ":" + f.Name
}

// SetValue - converts & assigns value of the filed
func (f *Field) SetValue(v interface{}) {
	// fmt.Println(reflect.ValueOf(v).Kind())

	f.value = v
	switch v.(type) {
	case string:
		f.valueType = Text
	case int:
		f.valueType = Integer
	case float32:
		f.valueType = Float
	case float64:
		f.valueType = Float
	case []byte:
		f.valueType = Blob
	default:
		f.valueType = Null
	}
}

// Int - return value as int
func (f *Field) Int() int {
	v, ok := f.value.(int)
	if !ok {
		log.Fatal()
	}
	return v
}

// Float - returns value as float64
func (f *Field) Float() float64 {
	v, ok := f.value.(float64)
	if !ok {
		log.Fatal()
	}
	return v
}

// String - return value as string
func (f *Field) String() string {
	v, ok := f.value.(string)
	if !ok {
		log.Fatal()
	}
	return v
}

// Blob - return value as array of bytes
func (f *Field) Blob() []byte {
	v, ok := f.value.([]byte)
	if !ok {
		log.Fatal()
	}
	return v
}
