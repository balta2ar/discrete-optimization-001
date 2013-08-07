package main

import "testing"

func BenchmarkLocalSearch(b *testing.B) {
    filename := "data/vrp_26_8_1"
    ctx := createContext(filename)
    solution := ctx.solveRandom()
    temperature := float64(100)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ctx.localSearch(solution, temperature)
    }
}
