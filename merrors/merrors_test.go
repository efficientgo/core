// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package merrors_test

import (
	"errors"
	"testing"

	corerrors "github.com/efficientgo/core/errors"
	"github.com/efficientgo/core/merrors"
	"github.com/efficientgo/core/testutil"
)

func TestNilMultiError(t *testing.T) {
	testutil.Ok(t, merrors.New().Err())
	testutil.Ok(t, merrors.New(nil, nil, nil).Err())

	e := merrors.New()
	e.Add()
	testutil.Ok(t, e.Err())

	e = merrors.New(nil, nil, nil)
	e.Add()
	testutil.Ok(t, e.Err())

	e = merrors.New()
	e.Add(nil, nil, nil)
	testutil.Ok(t, e.Err())

	e = merrors.New(nil, nil, nil)
	e.Add(nil, nil, nil)
	testutil.Ok(t, e.Err())
}

func TestMultiError(t *testing.T) {
	err := errors.New("test1")
	testutil.NotOk(t, merrors.New(err).Err())
	testutil.NotOk(t, merrors.New(nil, err, nil).Err())

	e := merrors.New(err)
	e.Add()
	testutil.NotOk(t, e.Err())

	e = merrors.New(nil, nil, nil)
	e.Add(err)
	testutil.NotOk(t, e.Err())

	e = merrors.New(err)
	e.Add(nil, nil, nil)
	testutil.NotOk(t, e.Err())

	e = merrors.New(nil, nil, nil)
	e.Add(nil, err, nil)
	testutil.NotOk(t, e.Err())

	testutil.NotOk(t, func() error {
		return e.Err()
	}())

	testutil.Ok(t, func() error {
		return merrors.New(nil, nil, nil).Err()
	}())
}

func TestMultiError_Error(t *testing.T) {
	err := errors.New("test1")

	testutil.Equals(t, "test1", merrors.New(err).Err().Error())
	testutil.Equals(t, "test1", merrors.New(err, nil).Err().Error())
	testutil.Equals(t, "4 errors: test1; test1; test2; test3", merrors.New(
		err,
		err,
		errors.New("test2"),
		nil,
		errors.New("test3"),
	).Err().Error())
}

type customErr struct{ error }

type customErr2 struct{ error }

type customErr3 struct{ error }

func TestMultiError_As(t *testing.T) {
	err := customErr{error: errors.New("err1")}

	testutil.Assert(t, errors.As(err, &err))
	testutil.Assert(t, errors.As(err, &customErr{}))

	testutil.Assert(t, !errors.As(err, &customErr2{}))
	testutil.Assert(t, !errors.As(err, &customErr3{}))

	// This is just to show limitation of std As.
	testutil.Assert(t, !errors.As(&err, &err))
	testutil.Assert(t, !errors.As(&err, &customErr{}))
	testutil.Assert(t, !errors.As(&err, &customErr2{}))
	testutil.Assert(t, !errors.As(&err, &customErr3{}))

	e := merrors.New(err).Err()
	testutil.Assert(t, errors.As(e, &customErr{}))
	same := merrors.New(err).Err()
	testutil.Assert(t, errors.As(e, &same))
	testutil.Assert(t, !errors.As(e, &customErr2{}))
	testutil.Assert(t, !errors.As(e, &customErr3{}))

	e2 := merrors.New(err, customErr3{error: errors.New("some")}).Err()
	testutil.Assert(t, errors.As(e2, &customErr{}))
	testutil.Assert(t, errors.As(e2, &customErr3{}))
	testutil.Assert(t, !errors.As(e2, &customErr2{}))

	// Wrapped.
	e3 := corerrors.Wrap(merrors.New(err, customErr3{}).Err(), "wrap")
	testutil.Assert(t, errors.As(e3, &customErr{}))
	testutil.Assert(t, errors.As(e3, &customErr3{}))
	testutil.Assert(t, !errors.As(e3, &customErr2{}))

	// This is just to show limitation of std As.
	e4 := merrors.New(err, &customErr3{}).Err()
	testutil.Assert(t, !errors.As(e4, &customErr2{}))
	testutil.Assert(t, !errors.As(e4, &customErr3{}))
}

func TestMultiError_Is(t *testing.T) {
	err := customErr{error: errors.New("err1")}

	testutil.Assert(t, errors.Is(err, err))
	testutil.Assert(t, errors.Is(err, customErr{error: err.error}))
	testutil.Assert(t, !errors.Is(err, &err))
	testutil.Assert(t, !errors.Is(err, customErr{}))
	testutil.Assert(t, !errors.Is(err, customErr{error: errors.New("err1")}))
	testutil.Assert(t, !errors.Is(err, customErr2{}))
	testutil.Assert(t, !errors.Is(err, customErr3{}))

	testutil.Assert(t, errors.Is(&err, &err))
	testutil.Assert(t, !errors.Is(&err, &customErr{error: err.error}))
	testutil.Assert(t, !errors.Is(&err, &customErr2{}))
	testutil.Assert(t, !errors.Is(&err, &customErr3{}))

	e := merrors.New(err).Err()
	testutil.Assert(t, errors.Is(e, err))
	testutil.Assert(t, errors.Is(err, customErr{error: err.error}))
	testutil.Assert(t, errors.Is(e, e))
	testutil.Assert(t, errors.Is(e, merrors.New(err).Err()))
	testutil.Assert(t, !errors.Is(e, &err))
	testutil.Assert(t, !errors.Is(err, customErr{}))
	testutil.Assert(t, !errors.Is(e, customErr2{}))
	testutil.Assert(t, !errors.Is(e, customErr3{}))

	e2 := merrors.New(err, customErr3{}).Err()
	testutil.Assert(t, errors.Is(e2, err))
	testutil.Assert(t, errors.Is(e2, customErr3{}))
	testutil.Assert(t, errors.Is(e2, merrors.New(err, customErr3{}).Err()))
	testutil.Assert(t, !errors.Is(e2, merrors.New(customErr3{}, err).Err()))
	testutil.Assert(t, !errors.Is(e2, customErr{}))
	testutil.Assert(t, !errors.Is(e2, customErr2{}))

	// Wrapped.
	e3 := corerrors.Wrap(merrors.New(err, customErr3{}).Err(), "wrap")
	testutil.Assert(t, errors.Is(e3, err))
	testutil.Assert(t, errors.Is(e3, customErr3{}))
	testutil.Assert(t, !errors.Is(e3, customErr{}))
	testutil.Assert(t, !errors.Is(e3, customErr2{}))

	exact := &customErr3{}
	e4 := merrors.New(err, exact).Err()
	testutil.Assert(t, errors.Is(e4, err))
	testutil.Assert(t, errors.Is(e4, exact))
	testutil.Assert(t, errors.Is(e4, merrors.New(err, exact).Err()))
	testutil.Assert(t, !errors.Is(e4, customErr{}))
	testutil.Assert(t, !errors.Is(e4, customErr2{}))
	testutil.Assert(t, !errors.Is(e4, &customErr3{}))
}

func TestMultiError_Count(t *testing.T) {
	err := customErr{error: errors.New("err1")}
	merr := merrors.New()
	merr.Add(customErr3{})

	m, ok := merrors.AsMulti(merr.Err())
	testutil.Assert(t, ok)
	testutil.Equals(t, 0, m.Count(err))
	testutil.Equals(t, 1, m.Count(customErr3{}))

	merr.Add(customErr3{})
	merr.Add(customErr3{})

	m, ok = merrors.AsMulti(merr.Err())
	testutil.Assert(t, ok)
	testutil.Equals(t, 0, m.Count(err))
	testutil.Equals(t, 3, m.Count(customErr3{}))

	// Nest multi errors with wraps.
	merr2 := merrors.New()
	merr2.Add(customErr3{})
	merr2.Add(customErr3{})
	merr2.Add(customErr3{})

	merr3 := merrors.New()
	merr3.Add(customErr3{})
	merr3.Add(customErr3{})

	// Wrap it so Add cannot add inner errors in.
	merr2.Add(corerrors.Wrap(merr3.Err(), "wrap"))
	merr.Add(corerrors.Wrap(merr2.Err(), "wrap"))

	m, ok = merrors.AsMulti(merr.Err())
	testutil.Assert(t, ok)
	testutil.Equals(t, 0, m.Count(err))
	testutil.Equals(t, 8, m.Count(customErr3{}))
}

func TestAsMulti(t *testing.T) {
	err := customErr{error: errors.New("err1")}
	merr := merrors.New(err, customErr3{}).Err()
	wrapped := corerrors.Wrap(merr, "wrap")

	_, ok := merrors.AsMulti(err)
	testutil.Assert(t, !ok)

	m, ok := merrors.AsMulti(merr)
	testutil.Assert(t, ok)
	testutil.Assert(t, errors.Is(m, merr))

	m, ok = merrors.AsMulti(wrapped)
	testutil.Assert(t, ok)
	testutil.Assert(t, errors.Is(m, merr))
}
