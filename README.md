# core

Go module with set of core packages **every** Go project needs. Minimal API, battle-tested, strictly versioned and with only one transient dependency--[davecgh/go-spew](https://github.com/davecgh/go-spew).

Maintained by experienced Go developers, including author of the [Efficient Go book](https://www.oreilly.com/library/view/efficient-go/9781098105709/).

Import it using `go get "github.com/efficientgo/core@latest`.

This module contains packages around the following functionalities:

> NOTE: Click on each package to see usage examples in pkg.go.dev!

## Error Handling

* `github.com/efficientgo/core/errors` is an improved and minimal version of the popular [`pkg/errors`](https://github.com/pkg/errors) package (archived) allowing reliable wrapping of errors with stacktrace. Unfortunately, the standard library recommended error wrapping using `%+w` is prone to errors and does not support stacktraces.
* `github.com/efficientgo/core/merrors` implements type safe collection of multiple errors. It presentings them in a unified way as a single `error` interface.
* `github.com/efficientgo/core/errcapture` offers readable and robust error handling in defer statement using Go return arguments.
* `github.com/efficientgo/core/logerrcapture` is similar to `errcapture`, but instead of appending potential error it logs it using logger interface.

## Waiting and Retrying

* `github.com/efficientgo/core/runutil` offers `Retry` and `Repeat` functions which is often need in Go production code (e.g. repeating operation periodically) as well as tests (e.g. waiting on eventual results instead of sleeping).
* `github.com/efficientgo/core/backoff` offers backoff timers which increases wait time on every retry, incredibly useful in distributed system timeout functionalities.

## Testing

* `github.com/efficientgo/core/testutil` is a minimal testing utility with only few functions like `Assert`, `Ok`, `NotOk` for errors and `Equals`. It's an alternative to [testify](https://github.com/stretchr/testify) project which has a bit more bloated interface and larger dependencies.

## Initial Authors

* [`@bisakhmondal`](https://github.com/bisakhmondal) (`errors` package).
* [`@bwplotka`](https://bwplotka.dev)
* [`@GiedriusS`](https://github.com/GiedriusS) (`errcapture` and `logerrcapture`)
