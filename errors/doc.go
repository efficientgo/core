// Initially copied from Thanos and contributed by https://github.com/bisakhmondal.
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package errors

// Package errors provides basic utilities to manipulate errors with a useful stacktrace. It combines the
// benefits of errors.New and fmt.Errorf world into a single package.
//
// The idea of writing errors package in thanos is highly motivated from the Tast project of Chromium OS Authors. However, instead of
// copying the package, we end up writing our own simplified logic borrowing some ideas from the errors and github.com/pkg/errors.
// A big thanks to all of them.
//
// Example 1:
//
//
// Example 2:
//
