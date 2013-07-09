package main

import "math"
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
    index int32
    x, y float64
    active bool
}

type Points []Point

// type ByInt32 []int32
// func (self ByInt32) Len() int { return len(self) }
// func (self ByInt32) Less(i, j int) bool { return self[i] < self[j] }
// func (self ByInt32) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

func (ps Points) dist(i, j int) float64 {
    return math.Sqrt(math.Pow(ps[i].x - ps[j].x, 2) +
                     math.Pow(ps[i].y - ps[j].y, 2))
}

func (ps Points) nearestTo(j int) int {
    var nearest int = -1
    var minDist float64 = math.MaxFloat64
    for i := 0; i < len(ps); i++ {
        if (i == j) || (!ps[i].active) {
            continue
        } else if nearest == -1 {
            nearest = i
        } else {
            d := ps.dist(i, j)
            if d < minDist {
                minDist = d
                nearest = i
            }
        }
    }
    //fmt.Println("nearest to", j, "is", nearest, "-", minDist)
    return nearest
}

func (ps Points) solveGreedy() int {
    N := len(ps)
    currentPoint := 0
    nextPoint := 0
    var pathLen float64 = 0
    var pointOrder = make([]int, N)

    pointOrder[0] = currentPoint

    //fmt.Println(pointOrder)
    for i := 1; i < N; i++ {
        nextPoint = ps.nearestTo(currentPoint)
        //fmt.Println(nextPoint)
        pointOrder[i] = nextPoint
        pathLen += ps.dist(currentPoint, nextPoint)
        ps[currentPoint].active = false

        //fmt.Println("turn off", currentPoint)
        //fmt.Println(ps)

        currentPoint = nextPoint
        //fmt.Println(pointOrder)
    }

    fmt.Println(pathLen, 0)
    for i := 0; i < N; i++ {
        fmt.Printf("%d ", pointOrder[i])
    }
    fmt.Printf("\n")
    //fmt.Println(ps)
    return 0
}

// TODO:
// 1. select best of greedy solutions (try all points as a starting point)
// 2. build distMatrix
// 3. build nearestMatrix
// 4. implement 2-opt, k-opt
//

func solveFile(filename string, alg string) int {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Cannot open file:", filename, err)
        return 2
    }
    defer file.Close()

    var N int
    var x, y float64
    fmt.Fscanf(file, "%d", &N)

    ps := Points(make([]Point, N))

    for i := 0; i < N; i++ {
        fmt.Fscanf(file, "%f %f", &x, &y)
        ps[i] = Point{int32(i), x, y, true}
    }

    switch {
    case alg == "greedy":
        return ps.solveGreedy()
    case alg == "csp":
        return 1
    default:
        return ps.solveGreedy()
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
