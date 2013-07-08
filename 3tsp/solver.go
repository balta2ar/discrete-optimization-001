package main

import "fmt"
import "os"

// functions which Go developers should have implemented but happened
// to be too lazy and religious to do so

func max(a int32, b int32) (r int32) {
    if a > b {
        return a
    } else {
        return b
    }
}

//
// Internal ADT
//

type Point struct {
    x, y float64
}

type Points []Point

// type ByInt32 []int32
// func (self ByInt32) Len() int { return len(self) }
// func (self ByInt32) Less(i, j int) bool { return self[i] < self[j] }
// func (self ByInt32) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

func solveGreedy(ps Points) int {
    N := len(ps)
    fmt.Println(N)
    return 0
}

func solveFile(filename string, alg string) int {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Cannot open file:", filename, err)
        return 2
    }
    defer file.Close()

    var N int
    fmt.Fscanf(file, "%d", &N)

    ps := make([]Point, N)

    for i := 0; i < N; i++ {
        fmt.Fscanf(file, "%f %f", &ps[i].x, &ps[i].y)
    }

    switch {
    case alg == "greedy":
        return solveGreedy(ps)
    case alg == "csp":
        return 1
    default:
        return solveGreedy(ps)
    }

    return 0
}

func main() {
    alg := "auto"
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    os.Exit(solveFile(os.Args[1], alg))
}
