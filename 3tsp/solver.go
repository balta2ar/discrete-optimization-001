package main

import "time"
import "math"
import "math/rand"
import "fmt"
import "log"
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

type FollowPoint struct {
    next, prev int
}

type FollowList []FollowPoint

type Solution struct {
    order []int
    cost float64
}

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

func (ps Points) nearestToExceptSmallerThan(j, a, b int, maxDist float64) int {
    var nearest int = -1
    var minDist float64 = math.MaxFloat64
    for i := 0; i < len(ps); i++ {
        if (i == j) || (i == a) || (i == b) { //|| (!ps[i].active) {
            continue
        } else if nearest == -1 {
            nearest = i
        } else {
            d := ps.dist(i, j)
            if d >= maxDist {
                continue
            }
            if d < minDist {
                minDist = d
                nearest = i
            }
        }
    }
    //fmt.Println("nearest to", j, "is", nearest, "-", minDist)
    return nearest
}

func printSolution(solution Solution) {
    fmt.Println(solution.cost, 0)
    for i := 0; i < len(solution.order); i++ {
        fmt.Printf("%d ", solution.order[i])
    }
    fmt.Printf("\n")
}

func (ps Points) solveGreedy() Solution {
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

    pathLen += ps.dist(pointOrder[N-1], pointOrder[0])

    //follow := make(FollowList, N)
    //log.Println(follow)

    //fmt.Println(ps)
    //return pointOrder
    return Solution{pointOrder, pathLen}
}

func findInSlice(what int, where []int) int {
    pos := -1
    for i := 0; i < len(where); i++ {
        if where[i] == what {
            return i
        }
    }
    return pos
}

func (ps Points) calcCost(solution Solution, pr bool) float64 {
    cost := 0.0
    N := len(solution.order)
    for i := 0; i < N; i++ {
        d := ps.dist(solution.order[i], solution.order[(i+1) % N])
        if pr {
           log.Println(d)
        }
        cost += d
    }
    //cost += ps.dist(solution.order[N-1], solution.order[0])
    return cost
}

// connect t1->t4, t2-t3, and reverse path between t2 and t4
// func reconnectPoints(selected, t1, t2, t3 int, origSolution Solution) Solution {
//     N := len(origSolution.order)
// 
//     solution := origSolution
//     solution.order = make([]int, N)
//     copy(solution.order, origSolution.order)
// 
//     t3InOrder := findInSlice(t3, solution.order)
//     t3InOrderPrev := (t3InOrder-1) % N
//     if t3InOrderPrev < 0 {
//         // stupid Go
//         t3InOrderPrev = N + t3InOrderPrev
//     }
//     //log.Println("t3InOrderPrev", t3InOrderPrev)
//     t4 := solution.order[t3InOrderPrev]
//     //log.Println("t3InOrder", t3InOrder, "t4", t4)
// 
//     // there is a part of graph order which needs to be reversed
//     // from next(t2) == selected+2 (inclusive)
//     // to t4 == t3InOrder-1 (not inclusive)
//     from := selected+2 // inclusive
//     to := t3InOrder-1 // not inclusive
//     var length int
//     if from <= to {
//         length = to-from
//     } else {
//         length = (N-from) + to
//     }
//     orderPart := make([]int, length)
//     for i := 0; i < length; i++ {
//         orderPart[i] = solution.order[(from+i) % N]
//     }
//     //log.Println("order part", orderPart)
// 
//     // reverse
//     for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
//         orderPart[i], orderPart[j] = orderPart[j], orderPart[i]
//     }
// 
//     // now fix solution order
//     ptr := selected+1
// 
//     // t1 - - -> t4
//     solution.order[ptr % N] = t4
//     ptr++
// 
//     // insert reversed part order
//     for i := 0; i < len(orderPart); i++ {
//         solution.order[ptr % N] = orderPart[i]
//         ptr++
//     }
// 
//     // insert t2 => t3 connection
//     solution.order[ptr % N] = t2
//     ptr++
//     solution.order[ptr % N] = t3
// 
//     return solution
// }

func reconnectPoints(p1, p3 int, origSolution Solution) Solution {
    N := len(origSolution.order)

    solution := origSolution
    solution.order = make([]int, N)
    copy(solution.order, origSolution.order)

    //t1 := solution.order[p1]
    t2 := solution.order[(p1+1) % N]

    t3 := solution.order[(p3+1) % N]
    t4 := solution.order[p3]

    // t3InOrder := findInSlice(t3, solution.order)
    // t3InOrderPrev := (t3InOrder-1) % N
    // if t3InOrderPrev < 0 {
    //     // stupid Go
    //     t3InOrderPrev = N + t3InOrderPrev
    // }

    //log.Println("t3InOrderPrev", t3InOrderPrev)
    //t4 := solution.order[t3InOrderPrev]
    //log.Println("t3InOrder", t3InOrder, "t4", t4)

    t3InOrder := (p3+1) % N
    //t3InOrderPrev := p3

    selected := p1

    // there is a part of graph order which needs to be reversed
    // from next(t2) == selected+2 (inclusive)
    // to t4 == t3InOrder-1 (not inclusive)
    from := selected+2 // inclusive
    to := t3InOrder-1 // not inclusive
    var length int
    if from <= to {
        length = to-from
    } else {
        length = (N-from) + to
    }
    orderPart := make([]int, length)
    for i := 0; i < length; i++ {
        orderPart[i] = solution.order[(from+i) % N]
    }
    //log.Println("order part", orderPart)

    // reverse
    for i, j := 0, length-1; i < j; i, j = i+1, j-1 {
        orderPart[i], orderPart[j] = orderPart[j], orderPart[i]
    }

    // now fix solution order
    ptr := selected+1

    // t1 - - -> t4
    solution.order[ptr % N] = t4
    ptr++

    // insert reversed part order
    for i := 0; i < len(orderPart); i++ {
        solution.order[ptr % N] = orderPart[i]
        ptr++
    }

    // insert t2 => t3 connection
    solution.order[ptr % N] = t2
    ptr++
    solution.order[ptr % N] = t3

    return solution
}

func (ps Points) kOpt(solution Solution) Solution {
    N := len(solution.order)
    log.Println("N", N)

    //ns := solution


    //solution
    //return solution

    // selected := 0 //rand.Int() % N // 3
    // t1 := solution.order[selected]
    // t2 := solution.order[(selected+1) % N]

    // x1 := ps.dist(t1, t2)
    // t2Next := solution.order[(selected+2) % N]
    // t3 := ps.nearestToExceptSmallerThan(t2, t1, t2Next, x1)

    //log.Println("current cost", solution.cost, "N", N)
    //log.Println("selected", selected, "t1", t1, "t2", t2)
    //log.Println("x1", x1, "t2Next", t2Next, "t3", t3)

    for i := 0; i < N; i++ {
        for j := i+2; j < N; j++ {
            newSolution := reconnectPoints(i, j, solution)
            newSolution.cost = ps.calcCost(newSolution, false)

            //log.Println(newSolution)

            if newSolution.cost < solution.cost {
                solution = newSolution
                solution.cost = newSolution.cost
                log.Println("BETTER SOLUTION FOUND", solution.cost)
                //return newSolution
            }
        }
    }

    // if t3 != -1 {
    //     //newSolution := reconnectPoints(selected, t1, t2, t3, solution)
    //     newSolution := reconnectPoints(i, j, solution)
    //     newSolution.cost = ps.calcCost(solution)

    //     if newSolution.cost < solution.cost {
    //         log.Println("BETTER SOLUTION FOUND")
    //         //return newSolution
    //         solution = newSolution
    //     }
    // }

    return solution
}

// TODO:
// 1. select best of greedy solutions (try all points as a starting point)
// 2. pre-compute distMatrix
// 3. pre-compute nearestMatrix
// 4. implement 2-opt, k-opt
// 5. use double-linked slice instead of order list
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
        solution := ps.solveGreedy()
        printSolution(solution)
    case alg == "csp":
        return 1
    default:
        solution := ps.solveGreedy()
        //c := ps.calcCost(solution, false)
        //fmt.Println(c)
        //printSolution(solution)

        solution = ps.kOpt(solution)
        // log.Println(solution.cost)
        printSolution(solution)

        // c = ps.calcCost(solution, false)
        // fmt.Println(c)
    }

    return 0
}

func main() {
    //log.Println((3+3) % 5)
    //log.Println((1-3) % 5)
    //return

    rand.Seed(time.Now().UTC().UnixNano())
    alg := "auto"
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    os.Exit(solveFile(os.Args[1], alg))
}
