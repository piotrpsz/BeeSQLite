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
	"bytes"
	"encoding/binary"
	"fmt"
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
	data      []byte
}

// BindName - computed name used in binding
func (f *Field) BindName() string {
	return ":" + f.Name
}

// SetValue - converts & assigns value of the filed
func (f *Field) SetValue(v interface{}) {
	// fmt.Println(reflect.ValueOf(v).Kind())

	switch x := v.(type) {
	case string:
		f.data = []byte(x)
		f.valueType = Text
	case int:
		f.data = convert(int64(x))
		f.valueType = Integer
	case float32:
		f.data = convert(float64(x))
		f.valueType = Float
	case float64:
		f.data = convert(x)
		f.valueType = Float
	case bool:
		bi := 1
		if !x {
			bi = 0
		}
		f.data = convert(int64(bi))
		f.valueType = Integer
	case []byte:
		f.data = x
		f.valueType = Blob
	default:
		f.data = []byte{}
		f.valueType = Null
	}
}

func convert(v interface{}) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.LittleEndian, v)
	if err != nil {
		fmt.Println("Field.convert: " + err.Error())
		return []byte{}
	}
	return buff.Bytes()

}

// Int - return value as int
func (f *Field) Int() int {
	var value int64
	buff := bytes.NewReader(f.data)
	err := binary.Read(buff, binary.LittleEndian, &value)
	if err != nil {
		log.Fatal(err)
	}
	return int(value)
}

// Float - returns value as float64
func (f *Field) Float() float64 {
	var value float64
	buff := bytes.NewReader(f.data)
	err := binary.Read(buff, binary.LittleEndian, &value)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

// String - return value as string
func (f *Field) String() string {
	return string(f.data)
}

// Blob - return value as array of bytes
func (f *Field) Blob() []byte {
	return f.data
}

// Bool - return value as bool
func (f *Field) Bool() bool {
	v := f.Int()
	if v == 1 {
		return true
	}
	return false
}
