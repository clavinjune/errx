// Copyright 2022 clavinjune/errutil
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errutil

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

const (
	skipRuntimeCaller int = 3
)

// Err is a custom error
type Err struct {
	// Caused is the error that caused the current error
	Caused error

	// FileLine is the file and line number of the error
	FileLine string

	// FuncName is the name of the function that caused the error
	FuncName string

	// Message is the message of the error
	Message string
}

// Unwrap returns internal error
func (e *Err) Unwrap() error {
	return e.Caused
}

// Err returns error message
func (e *Err) Error() string {
	var msg string

	inner := e.Caused.Error()
	if strings.HasPrefix(inner, "{") {
		msg = fmt.Sprintf(`"caused":%s`, inner)
	} else {
		msg = fmt.Sprintf(`"caused":%q`, inner)
	}

	msg += fmt.Sprintf(`,"funcname":"%s","fileline":"%s"`, e.FuncName, e.FileLine)

	if e.Message != "" {
		msg += fmt.Sprintf(`,"message":%q`, e.Message)
	}
	return fmt.Sprintf(`{%s}`, msg)
}

// Wrap creates *Err by wrapping error
func Wrap(err error) *Err {
	e := createError()
	e.Caused = err
	return e
}

// New creates *Err from text
func New(text string) *Err {
	e := createError()
	e.Caused = errors.New(text)
	return e
}

// WrapWithMsg creates *Err by wrapping error with a custom message
func WrapWithMsg(err error, msg string) *Err {
	e := createError()
	e.Caused = err
	e.Message = msg
	return e
}

func createError() *Err {
	fl, fn := getFlFn()

	return &Err{
		FileLine: fl,
		FuncName: fn,
	}
}

func getFlFn() (fl, fn string) {
	fl, fn = "?", "?"
	pc, file, line, ok := runtime.Caller(skipRuntimeCaller)
	if ok {
		fl = constructFileLine(file, line)
		fn = getFuncName(pc)
	}

	return
}

func constructFileLine(file string, line int) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}

func getFuncName(pc uintptr) string {
	f := runtime.FuncForPC(pc)

	split := strings.Split(f.Name(), "/")
	return split[len(split)-1]
}
