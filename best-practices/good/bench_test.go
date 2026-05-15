package main

import "testing"

// BenchmarkProcess covers the three modes against a realistic input slice.
// Run via `make bench` to capture allocs/op for benchstat comparison.
func BenchmarkProcess(b *testing.B) {
	items := []string{"alpha", "bb", "gamma", "dd", "epsilon", "ff", "eta", "hh"}
	modes := []struct {
		name string
		mode int
	}{
		{"mode1", 1},
		{"mode2", 2},
		{"default", 99},
	}
	for _, m := range modes {
		b.Run(m.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = process(items, m.mode)
			}
		})
	}
}

func BenchmarkHashSecret(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = hashSecret("hunter2")
	}
}
