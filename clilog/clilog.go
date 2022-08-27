// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package clilog

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/efficientgo/core/merrors"
)

// Logger interface compatible with go-kit/logger.
type Logger interface {
	Log(keyvals ...interface{}) error
}

type buf struct {
	bytes.Buffer

	*Encoder
}

func (l *buf) Reset() {
	l.Encoder.Reset()
	l.Buffer.Reset()
}

var bufPool = sync.Pool{
	New: func() interface{} {
		var b buf
		b.Encoder = NewEncoder(&b.Buffer)
		return &b
	},
}

type logger struct {
	w io.Writer
}

// New returns a logger that encodes keyvals to the Writer in
// CLI friendly format. Each log event produces no more than one call to w.Write.
// The passed Writer must be safe for concurrent use by multiple goroutines if
// the returned Logger will be used concurrently.
func New(w io.Writer) Logger {
	return &logger{w}
}

func (l logger) Log(keyvals ...interface{}) error {
	buf := bufPool.Get().(*buf)
	buf.Reset()
	defer bufPool.Put(buf)

	if err := buf.EncodeKeyvals(keyvals...); err != nil {
		return err
	}

	// Add newline to the end of the buffer
	if err := buf.EndRecord(); err != nil {
		return err
	}

	// The Logger interface requires implementations to be safe for concurrent
	// use by multiple goroutines. For this implementation that means making
	// only one call to l.w.Write() for each call to Log.
	if _, err := l.w.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// MarshalKeyvals returns the clilog encoding of keyvals, a variadic sequence
// of alternating keys and values.
func MarshalKeyvals(keyvals ...interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := NewEncoder(buf).EncodeKeyvals(keyvals...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// An Encoder writes clilog data to an output stream.
type Encoder struct {
	w       io.Writer
	scratch bytes.Buffer
	needSep bool

	errs []merrors.Error
}

// NewEncoder returns a new clilog Encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

var (
	sep     = []byte(": ")
	newline = []byte("\n")
	null    = []byte("null")
)

// EncodeKeyval writes the clilog encoding of key and value to the stream. A
// single space is written before the second and subsequent keys in a record.
// Nothing is written if a non-nil error is returned.
func (enc *Encoder) EncodeKeyval(_, value interface{}) error {
	if e, ok := value.(error); ok {
		if errs, ok := merrors.AsMulti(e); ok {
			enc.errs = append(enc.errs, errs)
			return nil
		}
	}

	enc.scratch.Reset()

	// Naive right now, print values only (:
	if enc.needSep {
		if _, err := enc.scratch.Write(sep); err != nil {
			return err
		}
	}
	if err := writeValue(&enc.scratch, value); err != nil {
		return err
	}
	_, err := enc.w.Write(enc.scratch.Bytes())
	enc.needSep = true
	return err
}

// EncodeKeyvals writes the logfmt encoding of keyvals to the stream. Keyvals
// is a variadic sequence of alternating keys and values. Keys of unsupported
// type are skipped along with their corresponding value. Values of
// unsupported type or that cause a MarshalerError are replaced by their error
// but do not cause EncodeKeyvals to return an error. If a non-nil error is
// returned some key/value pairs may not have be written.
func (enc *Encoder) EncodeKeyvals(keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if len(keyvals)%2 == 1 {
		keyvals = append(keyvals, nil)
	}
	for i := 0; i < len(keyvals); i += 2 {
		k, v := keyvals[i], keyvals[i+1]
		err := enc.EncodeKeyval(k, v)
		if err == ErrUnsupportedKeyType {
			continue
		}
		if _, ok := err.(*MarshalerError); ok || err == ErrUnsupportedValueType {
			v = err
			err = enc.EncodeKeyval(k, v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// EndRecord ends the log record.
func (enc *Encoder) EndRecord() error {
	if len(enc.errs) > 0 {
		enc.scratch.Reset()
		if enc.needSep {
			if _, err := enc.scratch.Write(sep); err != nil {
				return err
			}
		}

		merr := merrors.Merge(enc.errs)
		if err := merrors.PrettyPrint(&enc.scratch, merr); err != nil {
			return err
		}

		if _, err := enc.w.Write(enc.scratch.Bytes()); err != nil {
			return err
		}
	}

	_, err := enc.w.Write(newline)
	if err == nil {
		enc.needSep = false
	}
	return err
}

// Reset resets the Encoder to the beginning of a new record.
func (enc *Encoder) Reset() {
	enc.needSep = false
}

// MarshalerError represents an error encountered while marshaling a value.
type MarshalerError struct {
	Type reflect.Type
	Err  error
}

func (e *MarshalerError) Error() string {
	return "error marshaling value of type " + e.Type.String() + ": " + e.Err.Error()
}

// ErrUnsupportedKeyType is returned by Encoder methods if a key has an
// unsupported type.
var ErrUnsupportedKeyType = errors.New("unsupported key type")

// ErrUnsupportedValueType is returned by Encoder methods if a value has an
// unsupported type.
var ErrUnsupportedValueType = errors.New("unsupported value type")

func writeValue(w io.Writer, value interface{}) error {
	switch v := value.(type) {
	case nil:
		return writeBytesValue(w, null)
	case string:
		return writeStringValue(w, v, true)
	case []byte:
		return writeBytesValue(w, v)
	case encoding.TextMarshaler:
		vb, err := safeMarshal(v)
		if err != nil {
			return err
		}
		if vb == nil {
			vb = null
		}
		return writeBytesValue(w, vb)
	case error:
		se, ok := safeError(v)
		return writeStringValue(w, se, ok)
	case fmt.Stringer:
		ss, ok := safeString(v)
		return writeStringValue(w, ss, ok)
	default:
		rvalue := reflect.ValueOf(value)
		switch rvalue.Kind() {
		case reflect.Array, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Struct:
			return ErrUnsupportedValueType
		case reflect.Ptr:
			if rvalue.IsNil() {
				return writeBytesValue(w, null)
			}
			return writeValue(w, rvalue.Elem().Interface())
		}
		return writeStringValue(w, fmt.Sprintf("%v", v), true) //nolint
	}
}

func writeStringValue(w io.Writer, value string, ok bool) error {
	var err error
	if ok && value == "null" {
		_, err = io.WriteString(w, `"null"`)
	} else {
		_, err = io.WriteString(w, value)
	}
	return err
}

func writeBytesValue(w io.Writer, value []byte) error {
	_, err := w.Write(value)
	return err
}

func safeError(err error) (s string, ok bool) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(err); v.Kind() == reflect.Ptr && v.IsNil() {
				s, ok = "null", false
			} else {
				s, ok = fmt.Sprintf("PANIC:%v", panicVal), false
			}
		}
	}()
	s, ok = err.Error(), true
	return
}

func safeString(str fmt.Stringer) (s string, ok bool) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(str); v.Kind() == reflect.Ptr && v.IsNil() {
				s, ok = "null", false
			} else {
				s, ok = fmt.Sprintf("PANIC:%v", panicVal), true
			}
		}
	}()
	s, ok = str.String(), true
	return
}

func safeMarshal(tm encoding.TextMarshaler) (b []byte, err error) {
	defer func() {
		if panicVal := recover(); panicVal != nil {
			if v := reflect.ValueOf(tm); v.Kind() == reflect.Ptr && v.IsNil() {
				b, err = nil, nil
			} else {
				b, err = nil, fmt.Errorf("panic when marshaling: %s", panicVal)
			}
		}
	}()
	b, err = tm.MarshalText()
	if err != nil {
		return nil, &MarshalerError{
			Type: reflect.TypeOf(tm),
			Err:  err,
		}
	}
	return
}
