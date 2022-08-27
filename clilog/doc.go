// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package clilog

// Logging formatter that transforms structure log entry into human readable, clean friendly entry
// suitable more for CLI tools.
//
// In details this means:
//
// * No special sign escaping.
// * No key printing.
// * Values separated with ': '
// * Support for pretty printing multi errors (including nested ones) in format of (<something>: <err1>; <err2>; ...; <errN>)
// * TODO(bwplotka): Support for multiple multilines.
//
// Compatible with `github.com/go-kit/kit/log.Logger`
