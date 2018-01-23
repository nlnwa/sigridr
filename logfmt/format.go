// Copyright 2014 Alan Shreve
// Copyright 2018 National Library of Norway
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// 	See the License for the specific language governing permissions and
// limitations under the License.

package logfmt

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/inconshreveable/log15"
)

const timeFormat = time.RFC3339
const floatFormat = 'f'
const errorKey = "LOG_ERROR_KEY_NOT_STRING"

type Format interface {
	Format(r *log15.Record) []byte
}

// FormatFunc returns a new Format object which uses
// the given function to perform record formatting.
func FormatFunc(f func(*log15.Record) []byte) Format {
	return formatFunc(f)
}

type formatFunc func(*log15.Record) []byte

func (f formatFunc) Format(r *log15.Record) []byte {
	return f(r)
}

type level struct {
	log15.Lvl
}

// Returns the name of a Lvl
func (l level) String() string {
	switch l.Lvl {
	case log15.LvlDebug:
		return "debug"
	case log15.LvlInfo:
		return "info"
	case log15.LvlWarn:
		return "warn"
	case log15.LvlCrit:
		fallthrough
	case log15.LvlError:
		return "error"
	default:
		return "bad"
	}
}

func LogbackFormat() Format {
	return FormatFunc(func(r *log15.Record) []byte {
		b := &bytes.Buffer{}
		lvl := strings.ToUpper(level{r.Lvl}.String())
		now := r.Time.Format(time.RFC3339)
		msg := r.Msg
		thread := "main"
		logger := "undefined"

		m := delineate(r.Ctx)

		if t, ok := m["thread"]; ok {
			thread = t
			delete(m, "thread")
		}
		if l, ok := m["logger"]; ok {
			logger = l
			delete(m, "logger")
		} else if caller, ok := m["fn"]; ok {
			logger = caller
			delete(m, "fn")
		}

		// timestamp
		b.WriteString(now)
		b.WriteByte(' ')

		// thread
		b.WriteByte('[')
		b.WriteString(thread)
		b.WriteByte(']')
		b.WriteByte(' ')

		// level
		b.WriteString(lvl)
		b.WriteByte(' ')

		// logger
		b.WriteString(logger)

		// -
		b.WriteString(" - ")

		// fields
		i := len(m)

		b.WriteByte('{')
		for key, value := range m {
			b.WriteString(key)
			b.WriteByte('=')
			b.WriteString(value)
			if i--; i > 0 {
				b.WriteString(", ")
			}
		}
		b.WriteByte('}')
		b.WriteByte(' ')

		// msg
		b.WriteString(msg)

		b.WriteByte('\n')

		return b.Bytes()
	})
}

func delineate(ctx []interface{}) map[string]string {
	m := make(map[string]string, len(ctx))
	for i := 0; i < len(ctx); i += 2 {
		key, ok := ctx[i].(string)
		v := formatValue(ctx[i+1])
		if !ok {
			m[errorKey] = v
		} else {
			m[key] = v
		}
	}
	return m
}

func formatShared(value interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr && v.IsNil() {
				result = "nil"
			} else {
				panic(err)
			}
		}
	}()

	switch v := value.(type) {
	case time.Time:
		return v.Format(timeFormat)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

// formatValue formats a value for serialization
func formatValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	if t, ok := value.(time.Time); ok {
		// Performance optimization: No need for escaping since the provided
		// timeFormat doesn't have any escape characters, and escaping is
		// expensive.
		return t.Format(timeFormat)
	}
	value = formatShared(value)
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), floatFormat, 3, 64)
	case float64:
		return strconv.FormatFloat(v, floatFormat, 3, 64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", value)
	case string:
		return escapeString(v)
	default:
		return escapeString(fmt.Sprintf("%+v", value))
	}
}

var stringBufPool = sync.Pool{
	New: func() interface{} { return new(bytes.Buffer) },
}

func escapeString(s string) string {
	needsQuotes := false
	needsEscape := false
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' {
			needsQuotes = true
		}
		if r == '\\' || r == '"' || r == '\n' || r == '\r' || r == '\t' {
			needsEscape = true
		}
	}
	if needsEscape == false && needsQuotes == false {
		return s
	}
	e := stringBufPool.Get().(*bytes.Buffer)
	e.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\', '"':
			e.WriteByte('\\')
			e.WriteByte(byte(r))
		case '\n':
			e.WriteString("\\n")
		case '\r':
			e.WriteString("\\r")
		case '\t':
			e.WriteString("\\t")
		default:
			e.WriteRune(r)
		}
	}
	e.WriteByte('"')
	var ret string
	if needsQuotes {
		ret = e.String()
	} else {
		ret = string(e.Bytes()[1 : e.Len()-1])
	}
	e.Reset()
	stringBufPool.Put(e)
	return ret
}
