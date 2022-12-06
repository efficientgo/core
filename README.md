# core

[![core module docs](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/efficientgo/core)

Go module with set of core packages **every** Go project needs. Minimal API, battle-tested, strictly versioned and with only two transient dependencies-- [davecgh/go-spew](https://github.com/davecgh/go-spew) and [google/go-cmp](https://github.com/google/go-cmp).

Maintained by experienced Go developers, including author of the [Efficient Go book](https://www.oreilly.com/library/view/efficient-go/9781098105709/).

Import it using `go get "github.com/efficientgo/core@latest`.

This module contains packages around the following functionalities:

> NOTE: Click on each package to see usage examples in pkg.go.dev!

## Error Handling

* [github.com/efficientgo/core/errors](https://pkg.go.dev/github.com/efficientgo/core/errors) is an improved and minimal version of the popular [`pkg/errors`](https://github.com/pkg/errors) package (archived) allowing reliable wrapping of errors with stacktrace. Unfortunately, the standard library recommended error wrapping using `%+w` is prone to errors and does not support stacktraces. For example:

```go
var (
	err     = errors.New("root error while doing A")
	wrapped = errors.Wrap(err, "while doing B")
)  
```

* [github.com/efficientgo/core/merrors](https://pkg.go.dev/github.com/efficientgo/core/merrors) implements type safe collection of multiple errors. It presentings them in a unified way as a single `error` interface.

```go
func CloseAll(closers []io.Closer) error {
	errs := merrors.New()
	for _, c := range closers {
		errs.Add(c.Close())
	}
	return errs.Err()
}
```

* [github.com/efficientgo/core/errcapture](https://pkg.go.dev/github.com/efficientgo/core/errcapture) offers readable and robust error handling in defer statement using Go return arguments.

```go
func DoAndClose(f *os.File) (err error) {
	defer errcapture.Do(&err, f.Close, "close file at the end")

	// Do something...
	if err := do(); err != nil {
		return err
	}

	return nil
}   
```

* [github.com/efficientgo/core/logerrcapture](https://pkg.go.dev/github.com/efficientgo/core/logerrcapture) is similar to `errcapture`, but instead of appending potential error it logs it using logger interface.

```go
func DoAndClose(f *os.File, logger logerrcapture.Logger) error {
	defer logerrcapture.Do(logger, f.Close, "close file at the end")

	// Do something...
	if err := do(); err != nil {
		return err
	}

	return nil
}   
```

## Waiting and Retrying

* [github.com/efficientgo/core/runutil](https://pkg.go.dev/github.com/efficientgo/core/runutil) offers `Retry` and `Repeat` functions which is often need in Go production code (e.g. repeating operation periodically) as well as tests (e.g. waiting on eventual results instead of sleeping).

```go
// Repeat every 1 second until context is done (e.g. cancel or timeout) or
// function returns error.
err := runutil.Repeat(1*time.Second, ctx.Done(), func() error {
	// ...
	return err // Ups, error - don't repeat anymore!
})

// Retry every 1 second until context is done (e.g. cancel or timeout) or
// function returns nil.
err := runutil.Retry(1*time.Second, ctx.Done(), func() error {
	// ...
	return nil // Done, no need to retry!
}) 
```

* [github.com/efficientgo/core/backoff](https://pkg.go.dev/github.com/efficientgo/core/backoff) offers backoff timers which increases wait time on every retry, incredibly useful in distributed system timeout functionalities.

## Testing

* [github.com/efficientgo/core/testutil](https://pkg.go.dev/github.com/efficientgo/core/testutil) is a minimal testing utility with only few functions like `Assert`, `Ok`, `NotOk` for errors and `Equals`. It's an alternative to [testify](https://github.com/stretchr/testify) project which has a bit more bloated interface and larger dependencies.

```go
func TestSomething(t *testing.T) {
	got, err := something()
	testutil.Ok(t, err)
	testutil.Equals(t, expected, got, "expected different thing from something")
}
```

## Initial Authors

* [`@bisakhmondal`](https://github.com/bisakhmondal) (`errors` package).
* [`@bwplotka`](https://bwplotka.dev)
* [`@GiedriusS`](https://github.com/GiedriusS) (`errcapture` and `logerrcapture`)
