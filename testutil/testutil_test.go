// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package testutil

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestContains(t *testing.T) {
	tests := map[string]struct {
		haystack    []string
		needle      []string
		shouldMatch bool
	}{
		"empty haystack": {
			haystack:    []string{},
			needle:      []string{"key1"},
			shouldMatch: false,
		},

		"empty needle": {
			haystack:    []string{"key1", "key2", "key3"},
			needle:      []string{},
			shouldMatch: false,
		},

		"single value needle": {
			haystack:    []string{"key1", "key2", "key3"},
			needle:      []string{"key1"},
			shouldMatch: true,
		},

		"multiple value needle": {
			haystack:    []string{"key1", "key2", "key3"},
			needle:      []string{"key1", "key2"},
			shouldMatch: true,
		},

		"same size needle as haystack": {
			haystack:    []string{"key1", "key2", "key3"},
			needle:      []string{"key1", "key2", "key3"},
			shouldMatch: true,
		},

		"larger needle than haystack": {
			haystack:    []string{"key1", "key2", "key3"},
			needle:      []string{"key1", "key2", "key3", "key4"},
			shouldMatch: false,
		},

		"needle not contained": {
			haystack:    []string{"key1", "key2", "key3"},
			needle:      []string{"key4"},
			shouldMatch: false,
		},

		"haystack ends before needle": {
			haystack:    []string{"key1", "key2", "key3"},
			needle:      []string{"key3", "key4"},
			shouldMatch: false,
		},
	}

	for testName, testData := range tests {
		t.Run(testName, func(t *testing.T) {
			if testData.shouldMatch != contains(testData.haystack, testData.needle) {
				t.Fatalf("unexpected result testing contains() with %#v", testData)
			}
		})
	}
}

type child struct {
	Val float64
}

type parent struct {
	C child
}

type unexp struct {
	val float64
}

func TestEqualsWithNaN(t *testing.T) {
	for _, tc := range []struct {
		name string
		a    interface{}
		b    interface{}
		opts cmp.Options
	}{
		{
			a:    math.NaN(),
			b:    math.NaN(),
			name: "Simple NaN value comparison",
		},
		{
			a:    child{Val: math.NaN()},
			b:    child{Val: math.NaN()},
			name: "NaN value as struct member comparison",
		},
		{
			a:    parent{C: child{Val: math.NaN()}},
			b:    parent{C: child{Val: math.NaN()}},
			name: "NaN value in nested struct comparison",
		},
		{
			a:    unexp{val: math.NaN()},
			b:    unexp{val: math.NaN()},
			name: "NaN value as unexported struct member comparison",
			opts: cmp.Options{cmp.AllowUnexported(unexp{})},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			EqualsWithNaN(t, tc.a, tc.b, tc.opts)
		})
	}
}
