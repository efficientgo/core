// Initially copied from Thanos and contributed by https://github.com/bisakhmondal.
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package errors_test

import (
	//lint:ignore faillint Custom errors package tests need to import standard library errors.
	stderrors "errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
)

const msg = "test_error_message"
const wrapper = "test_wrapper"

func TestNewf(t *testing.T) {
	err := errors.Newf(msg)
	testutil.Equals(t, err.Error(), msg, "the root error message must match")

	// The %+v triggers stacktrace print.
	reg := regexp.MustCompile(msg + `[ \n]+> github\.com\/efficientgo\/core\/errors_test\.TestNewf	.*\/errors\/errors_test\.go:\d+`)
	testutil.Equals(t, reg.MatchString(fmt.Sprintf("%+v", err)), true, "matching stacktrace in errors.Newf")
}

func TestNewfFormatted(t *testing.T) {
	fmtMsg := msg + " key=%v"
	expectedMsg := msg + " key=value"

	err := errors.Newf(fmtMsg, "value")
	testutil.Equals(t, err.Error(), expectedMsg, "the root error message must match")

	fmt.Printf("%+v", err)

	reg := regexp.MustCompile(msg + `[ \n]+> github\.com\/efficientgo\/core\/errors_test\.TestNewfFormatted	.*\/errors\/errors_test\.go:\d+`)
	testutil.Equals(t, reg.MatchString(fmt.Sprintf("%+v", err)), true, "matching stacktrace in errors.TestNewfFormatted")
}

func TestWrapf(t *testing.T) {
	err := errors.Newf(msg)
	err = errors.Wrapf(err, wrapper)

	expectedMsg := wrapper + ": " + msg
	testutil.Equals(t, err.Error(), expectedMsg, "the root error message must match")

	reg := regexp.MustCompile(`test_wrapper[ \n]+> github\.com\/thanos-io\/thanos\/pkg\/errors\.TestWrapf	.*\/pkg\/errors\/errors_test\.go:\d+
[[:ascii:]]+test_error_message[ \n]+> github\.com\/thanos-io\/thanos\/pkg\/errors\.TestWrapf	.*\/pkg\/errors\/errors_test\.go:\d+`)

	testutil.Equals(t, reg.MatchString(fmt.Sprintf("%+v", err)), true, "matching stacktrace in errors.Wrap")
}

func TestUnwrap(t *testing.T) {
	// test with base error
	err := errors.Newf(msg)

	for i, tc := range []struct {
		err      error
		expected string
		isNil    bool
	}{
		{
			// no wrapping
			err:   err,
			isNil: true,
		},
		{
			err:      errors.Wrapf(err, wrapper),
			expected: "test_error_message",
		},
		{
			err:      errors.Wrapf(errors.Wrapf(err, wrapper), wrapper),
			expected: "test_wrapper: test_error_message",
		},
		// check primitives errors
		{
			err:   stderrors.New("std-error"),
			isNil: true,
		},
		{
			err:      errors.Wrapf(stderrors.New("std-error"), wrapper),
			expected: "std-error",
		},
		{
			err:   nil,
			isNil: true,
		},
	} {
		t.Run("TestCase"+strconv.Itoa(i), func(t *testing.T) {
			unwrapped := errors.Unwrap(tc.err)
			if tc.isNil {
				testutil.Equals(t, unwrapped, nil)
				return
			}
			testutil.Equals(t, unwrapped.Error(), tc.expected, "Unwrap must match expected output")
		})
	}
}

func TestCause(t *testing.T) {
	// test with base error that implements interface containing Unwrap method
	err := errors.Newf(msg)

	for i, tc := range []struct {
		err      error
		expected string
		isNil    bool
	}{
		{
			// no wrapping
			err:   err,
			isNil: true,
		},
		{
			err:   errors.Wrapf(err, wrapper),
			isNil: true,
		},
		{
			err:   errors.Wrapf(errors.Wrapf(err, wrapper), wrapper),
			isNil: true,
		},
		// check primitives errors
		{
			err:      stderrors.New("std-error"),
			expected: "std-error",
		},
		{
			err:      errors.Wrapf(stderrors.New("std-error"), wrapper),
			expected: "std-error",
		},
		{
			err:   nil,
			isNil: true,
		},
	} {
		t.Run("TestCase"+strconv.Itoa(i), func(t *testing.T) {
			cause := errors.Cause(tc.err)
			if tc.isNil {
				testutil.Equals(t, cause, nil)
				return
			}
			testutil.Equals(t, cause.Error(), tc.expected, "Cause must match expected output")
		})
	}
}
