package main

import "testing"
import "math"
// import "log"

func BenchmarkLocalSearch(b *testing.B) {
    filename := "data/vrp_26_8_1"
    ctx := createContext(filename)
    solution := ctx.solveRandom()
    temperature := float64(100)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ctx.localSearch(solution, temperature, 100000, 1.0)
    }
}

func BenchmarkSelectCustomerMove(b *testing.B) {
    filename := "data/vrp_26_8_1"
    ctx := createContext(filename)
    solution := ctx.solveRandom()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ctx.selectCustomerMove(solution)
    }
}

func TestOverCapacityAfterMove(b *testing.T) {
    filename := "data/vrp_26_8_1"
    ctx := createContext(filename)
    solution := ctx.solveRandom()
    // ctx.printSolution(solution)
    ctx.overCapacity(solution)

    for i := 0; i < 10000; i++ {
        move := ctx.selectCustomerMove(solution)
        // log.Println("new move i", i, move)
        incrOverCapacity := move.NewOverCapacity
        solution = ctx.applyMove(move, solution)
        calcOverCapacity := ctx.overCapacity(solution)
        // ctx.printSolution(solution)

        for _, path := range solution.Paths {
            incrDemand := path.Demand
            calcDemand := ctx.pathDemand(path)
            // log.Println("i", i, "j", j, "demand incr", incrDemand, "calc", calcDemand)
            if incrDemand != calcDemand {
                b.FailNow()
            }
        }

        if incrOverCapacity != calcOverCapacity {
            b.Fatal("incrOverCapacity", incrOverCapacity,
                    "!= calcOverCapacity", calcOverCapacity, "i", i)
        }
        // log.Println("---")
    }
}

func TestCostAfterMove(b *testing.T) {
    filename := "data/vrp_26_8_1"
    ctx := createContext(filename)
    solution := ctx.solveRandom()
    // ctx.printSolution(solution)

    for i := 0; i < 100000; i++ {
        move := ctx.selectCustomerMove(solution)
        // log.Println(move)
        incrCost := move.NewCost
        solution = ctx.applyMove(move, solution)
        calcCost := ctx.solutionCost(solution)
        // ctx.printSolution(solution)
        if math.Abs(float64(incrCost - calcCost)) > 0.1 {
            b.Fatal("incrCost", incrCost, "!= calcCost", calcCost, "i", i)
        }
    }
}

func TestLocalSearch(b *testing.T) {
    filename := "data/vrp_26_8_1"
    ctx := createContext(filename)
    solution := ctx.solveRandom()
    lastSolution := ctx.localSearch(cloneSolution(solution), 100, 14, 1.0)
    // ctx.printSolution(lastSolution)

    if math.Abs(float64(lastSolution.Cost - ctx.solutionCost(lastSolution))) > 0.1 {
        b.Fatal("incr Cost", lastSolution.Cost, "!= calc Cost", ctx.solutionCost(lastSolution))
    }

    if lastSolution.OverCapacity != ctx.overCapacity(lastSolution) {
        b.Fatal("incr OverCapacity", lastSolution.OverCapacity,
                "!= calc OverCapacity", ctx.overCapacity(lastSolution))
    }
}
