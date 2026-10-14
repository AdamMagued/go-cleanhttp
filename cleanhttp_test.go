// Copyright IBM Corp. 2015, 2026
// SPDX-License-Identifier: MPL-2.0

package cleanhttp

import (
	"runtime"
	"testing"
)

func TestCalcMaxIdleConnsPerHost(t *testing.T) {
	cases := []struct {
		name     string
		procs    int
		expected int
	}{
		{
			name:     "negative procs floor",
			procs:    -1,
			expected: 4,
		},
		{
			name:     "zero procs floor",
			procs:    0,
			expected: 4,
		},
		{
			name:     "single core floor",
			procs:    1,
			expected: 4,
		},
		{
			name:     "dual core floor",
			procs:    2,
			expected: 4,
		},
		{
			name:     "three cores boundary",
			procs:    3,
			expected: 4,
		},
		{
			name:     "four cores scaling",
			procs:    4,
			expected: 5,
		},
		{
			name:     "eight cores scaling",
			procs:    8,
			expected: 9,
		},
		{
			name:     "sixteen cores scaling",
			procs:    16,
			expected: 17,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := calcMaxIdleConnsPerHost(tc.procs)
			if actual != tc.expected {
				t.Fatalf("calcMaxIdleConnsPerHost(%d) = %d; want %d", tc.procs, actual, tc.expected)
			}
		})
	}
}

func TestDefaultPooledTransport_MaxIdleConnsPerHost(t *testing.T) {
	orig := runtime.GOMAXPROCS(0)
	defer runtime.GOMAXPROCS(orig)

	cases := []struct {
		name     string
		setProcs int
		expected int
	}{
		{
			name:     "procs 1 floor",
			setProcs: 1,
			expected: 4,
		},
		{
			name:     "procs 2 floor",
			setProcs: 2,
			expected: 4,
		},
		{
			name:     "procs 4 computed",
			setProcs: 4,
			expected: 5,
		},
		{
			name:     "procs 8 computed",
			setProcs: 8,
			expected: 9,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runtime.GOMAXPROCS(tc.setProcs)
			tr := DefaultPooledTransport()
			if tr.MaxIdleConnsPerHost != tc.expected {
				t.Fatalf("DefaultPooledTransport().MaxIdleConnsPerHost = %d; want %d", tr.MaxIdleConnsPerHost, tc.expected)
			}
		})
	}
}

func TestDefaultTransport(t *testing.T) {
	tr := DefaultTransport()
	if !tr.DisableKeepAlives {
		t.Fatal("expected DisableKeepAlives to be true")
	}
	if tr.MaxIdleConnsPerHost != -1 {
		t.Fatalf("expected MaxIdleConnsPerHost to be -1, got %d", tr.MaxIdleConnsPerHost)
	}
}

func TestDefaultClient(t *testing.T) {
	c := DefaultClient()
	if c == nil || c.Transport == nil {
		t.Fatal("expected non-nil client and transport")
	}
}

func TestDefaultPooledClient(t *testing.T) {
	c := DefaultPooledClient()
	if c == nil || c.Transport == nil {
		t.Fatal("expected non-nil client and transport")
	}
}
