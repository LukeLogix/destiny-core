package bazi

import (
	"testing"
	"time"
)

func benchBirth() Birth {
	tz, _ := time.LoadLocation("Asia/Taipei")
	return Birth{Time: time.Date(1990, 5, 20, 10, 30, 0, 0, tz), Gender: Male}
}

func BenchmarkComputeNatalOnly(b *testing.B) {
	opt := Default()
	for i := 0; i < b.N; i++ {
		if _, err := Compute(benchBirth(), opt); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkComputeWithDynamicShenSha(b *testing.B) {
	opt := Default()
	opt.IncludeDynamicShenSha = true
	for i := 0; i < b.N; i++ {
		if _, err := Compute(benchBirth(), opt); err != nil {
			b.Fatal(err)
		}
	}
}
