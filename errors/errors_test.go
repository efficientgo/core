// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

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
	"testing"

	"github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/testutil"
)

func ExampleWrap() {
	if err := func() error {
		// Do something...
		return errors.New("error!")
	}(); err != nil {
		wrapped := errors.Wrap(err, "doing something")
		fmt.Println(wrapped)
	}

	// Output: doing something: error!
}

func ExampleWrapf() {
	if err := func() error {
		// Do something...
		return errors.Newf("I am surprised %d + %d equals %v", 2, 3, 5)
	}(); err != nil {
		wrapped := errors.Wrapf(err, "doing some math %v times", 10)
		fmt.Println(wrapped)
	}

	// Output: doing some math 10 times: I am surprised 2 + 3 equals 5
}

const msg = "test_error_message"
const wrapper = "test_wrapper"

func TestNew(t *testing.T) {
	err := errors.New(msg)
	testutil.NotOk(t, err)
	testutil.Equals(t, err.Error(), msg, "the root error message must match")

	// The %+v triggers stacktrace print.
	reg := regexp.MustCompile(msg + `[ \n]+> github\.com\/efficientgo\/core\/errors_test\.TestNew	.*\/errors\/errors_test\.go:\d+`)
	testutil.Assert(t, reg.MatchString(fmt.Sprintf("%+v", err)), "matching stacktrace in errors.New")
}

func TestNewf(t *testing.T) {
	fmtMsg := msg + " key=%v"
	expectedMsg := msg + " key=value"

	err := errors.Newf(fmtMsg, "value")
	testutil.NotOk(t, err)
	testutil.Equals(t, err.Error(), expectedMsg, "the root error message must match")

	reg := regexp.MustCompile(expectedMsg + `[ \n]+> github\.com\/efficientgo\/core\/errors_test\.TestNewf	.*\/errors\/errors_test\.go:\d+`)
	testutil.Assert(t, reg.MatchString(fmt.Sprintf("%+v", err)), "matching stacktrace in errors.Newf")
}

func TestWrap(t *testing.T) {
	err := errors.New(msg)
	err = errors.Wrap(err, wrapper)
	testutil.NotOk(t, err)

	expectedMsg := wrapper + ": " + msg
	testutil.Equals(t, err.Error(), expectedMsg, "the root error message must match")

	reg := regexp.MustCompile(`test_wrapper[ \n]+> github\.com\/efficientgo\/core\/errors_test\.TestWrap	.*\/errors\/errors_test\.go:\d+
[[:ascii:]]+test_error_message[ \n]+> github\.com\/efficientgo\/core\/errors_test\.TestWrap	.*\/errors\/errors_test\.go:\d+`)
	testutil.Assert(t, reg.MatchString(fmt.Sprintf("%+v", err)), "not matching stacktrace in errors.Wrap")

	testutil.Ok(t, errors.Wrap(nil, wrapper))
}

func TestWrapf(t *testing.T) {
	err := errors.New(msg)

	fmtWrapper := wrapper + " key=%v"
	err = errors.Wrapf(err, fmtWrapper, "value")
	testutil.NotOk(t, err)

	expectedMsg := wrapper + " key=value: " + msg
	testutil.Equals(t, err.Error(), expectedMsg, "the root error message must match")

	reg := regexp.MustCompile(`[ \n]+> github\.com\/efficientgo\/core\/errors_test\.TestWrapf	.*\/errors\/errors_test\.go:\d+
[[:ascii:]]+test_error_message[ \n]+> github\.com\/efficientgo\/core\/errors_test\.TestWrapf	.*\/errors\/errors_test\.go:\d+`)
	testutil.Assert(t, reg.MatchString(fmt.Sprintf("%+v", err)), "not matching stacktrace in errors.Wrapf")

	testutil.Ok(t, errors.Wrapf(nil, fmtWrapper))
}

func TestUnwrap(t *testing.T) {
	baseErr := errors.Newf(msg)

	for _, tc := range []struct {
		err            error
		expectedUnwrap string
	}{
		{
			err: baseErr,
		},
		{
			err:            errors.Wrap(baseErr, wrapper),
			expectedUnwrap: "test_error_message",
		},
		{
			err:            errors.Wrap(errors.Wrap(baseErr, wrapper), wrapper),
			expectedUnwrap: "test_wrapper: test_error_message",
		},
		{
			err: stderrors.New("std-error"),
		},
		{
			err:            errors.Wrap(stderrors.New("std-error"), wrapper),
			expectedUnwrap: "std-error",
		},
		{
			err: nil,
		},
	} {
		t.Run("", func(t *testing.T) {
			unwrapped := errors.Unwrap(tc.err)
			if tc.expectedUnwrap == "" {
				testutil.Ok(t, unwrapped)
				return
			}
			testutil.NotOk(t, unwrapped)
			testutil.Equals(t, tc.expectedUnwrap, unwrapped.Error())
		})
	}
}

func TestCause(t *testing.T) {
	baseErr := errors.Newf(msg)

	for _, tc := range []struct {
		err           error
		expectedCause string
	}{
		{
			// no wrapping
			err: baseErr,
		},
		{
			err: errors.Wrap(baseErr, wrapper),
		},
		{
			err: errors.Wrap(errors.Wrapf(baseErr, wrapper), wrapper),
		},
		{
			err:           stderrors.New("std-error"),
			expectedCause: "std-error",
		},
		{
			err:           errors.Wrap(stderrors.New("std-error"), wrapper),
			expectedCause: "std-error",
		},
		{
			err: nil,
		},
	} {
		t.Run("", func(t *testing.T) {
			cause := errors.Cause(tc.err)
			if tc.expectedCause == "" {
				testutil.Ok(t, cause)
				return
			}
			testutil.NotOk(t, cause)
			testutil.Equals(t, tc.expectedCause, cause.Error())
		})
	}
}
