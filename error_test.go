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

package errutil_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net"
	"strings"
	"testing"

	"github.com/clavinjune/errutil"
	"github.com/stretchr/testify/require"
)

const (
	expectedFileLine = "error_test.go:"
	expectedFuncName = "errutil_test.TestErr_Error"
)

func helperCompareNestedError(r *require.Assertions, actual map[string]any, expected *errutil.Err) {
	r.Equal(expectedFuncName, actual["funcname"])
	r.True(strings.HasPrefix(actual["fileline"].(string), expectedFileLine))
	if expected.Message != "" {
		r.Equal(expected.Message, actual["message"])
	}

	switch caused := actual["caused"].(type) {
	case string:
		r.Equal(expected.Caused.Error(), caused)
	case map[string]any:
		helperCompareNestedError(r, caused, expected.Caused.(*errutil.Err))
	}
}

func TestErr_Error(t *testing.T) {
	type testCase struct {
		name     string
		error    *errutil.Err
		expected *errutil.Err
	}

	runFunc := func(tc testCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Parallel()
			r := require.New(t)

			var actual map[string]any
			r.NoError(json.Unmarshal([]byte(tc.error.Error()), &actual))
			helperCompareNestedError(r, actual, tc.expected)
		}
	}

	tt := []testCase{
		{
			name:  "new",
			error: errutil.New("this is an error"),
			expected: &errutil.Err{
				Caused: errors.New("this is an error"),
			},
		},
		{
			name:  "new with double quote",
			error: errutil.New(`"this is an error"`),
			expected: &errutil.Err{
				Caused: errors.New(`"this is an error"`),
			},
		},
		{
			name:  "wrap simple error",
			error: errutil.Wrap(sql.ErrNoRows),
			expected: &errutil.Err{
				Caused: sql.ErrNoRows,
			},
		},
		{
			name:  "wrap simple error with message",
			error: errutil.WrapWithMsg(sql.ErrNoRows, "wrap simple error with message"),
			expected: &errutil.Err{
				Caused:  sql.ErrNoRows,
				Message: "wrap simple error with message",
			},
		},
		{
			name:  "wrap simple error with double quoted message",
			error: errutil.WrapWithMsg(sql.ErrNoRows, `"wrap simple error with double quoted message"`),
			expected: &errutil.Err{
				Caused:  sql.ErrNoRows,
				Message: `"wrap simple error with double quoted message"`,
			},
		},
		{
			name: "nested",
			error: errutil.WrapWithMsg(
				errutil.WrapWithMsg(
					errutil.New("inner error"),
					"inner error message"),
				`"wrap simple error with double quoted message"`),
			expected: &errutil.Err{
				Caused: &errutil.Err{
					Caused: &errutil.Err{
						Caused: errors.New("inner error"),
					},
					Message: "inner error message",
				},
				Message: `"wrap simple error with double quoted message"`,
			},
		},
	}

	for i := range tt {
		tc := tt[i]
		t.Run(tc.name, runFunc(tc))
	}
}

func TestErr_Unwrap(t *testing.T) {
	r := require.New(t)

	var inner error = &net.AddrError{
		Err:  "test err",
		Addr: "test addr",
	}

	var err error = errutil.WrapWithMsg(inner, "wrap simple error with message")

	r.ErrorIs(err, inner)
	r.NotErrorIs(err, sql.ErrTxDone)

	var target *errutil.Err
	r.ErrorAs(err, &target)
	r.Equal(target, err)

	var innerTarget *net.AddrError
	r.ErrorAs(err, &innerTarget)
	r.Equal(innerTarget, inner)
}
