package main

import "time"
import "math"
import "sort"
import "math/rand"
import "fmt"
import "log"
import "os"

const (
    MAX_SECONDS_BETWEEN_CHANGES = 120
    // SA_MAX_ITERATIONS = 100
    // LS_MAX_TRIALS = 1000
)

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
    active bool
}

type Points []Point

type Context struct {
    ps Points
    distMatrix [][]float64
    nearestToMatrix [][]int32
    N int
}

type FollowPoint struct {
    next, prev int
}

type FollowList []FollowPoint

type Solution struct {
    order []int
    cost float64
}

// calc and cache distances from each to each point
// create triangle matrix to save space
func (ctx Context) calcDistMatrix() Context {
    ctx.distMatrix = make([][]float64, ctx.N)
    for i := 1; i < ctx.N; i++ {
        // ctx.distMatrix[i] = make([]float64, ctx.N)
        // for j := 0; j < ctx.N; j++ {
        //     ctx.distMatrix[i][j] = ctx.calcDist(i, j)
        // }

        ctx.distMatrix[i] = make([]float64, i)
        for j := 0; j < i; j++ {
            ctx.distMatrix[i][j] = ctx.calcDist(i, j)
        }
    }
    return ctx
}

// used to sort indexs by distance to some point
type IndexSorter struct {
    idx []int32
    from int32
    ctx Context
}

func (s *IndexSorter) Len() int { return s.ctx.N }
func (s *IndexSorter) Swap(i, j int) { s.idx[i], s.idx[j] = s.idx[j], s.idx[i] }
func (s *IndexSorter) Less(i, j int) bool {
    a := s.ctx.calcDist(int(s.from), int(s.idx[i]))
    b := s.ctx.calcDist(int(s.from), int(s.idx[j]))
    return a < b
}

// calc and cache sorted list of distances from each to each point
// (to be used in greedy alg)
// WARNING: depends on calcDistMatrix
// ctx.distMatrix MUST be calculated before calling this function
func (ctx Context) calcNearestToMatrix() Context {
    ctx.nearestToMatrix = make([][]int32, ctx.N)
    for i := 0; i < ctx.N; i++ {
        ctx.nearestToMatrix[i] = make([]int32, ctx.N)
        for j := 0; j < ctx.N; j++ {
            ctx.nearestToMatrix[i][j] = int32(j)
        }
        sort.Sort(&IndexSorter{ctx.nearestToMatrix[i], int32(i), ctx})
    }
    return ctx
}

func (ctx Context) init() Context {
    ctx = ctx.calcDistMatrix()
    //log.Println(ctx.distMatrix)
    ctx = ctx.calcNearestToMatrix()
    return ctx
}

func (ctx Context) calcDist(i, j int) float64 {
    return float64(math.Sqrt(math.Pow(float64(ctx.ps[i].x - ctx.ps[j].x), 2) +
                             math.Pow(float64(ctx.ps[i].y - ctx.ps[j].y), 2)))
}

func (ctx Context) dist(i, j int) float64 {
    // return ctx.distMatrix[i][j]

    if i == j {
        return 0.0
    }
    if j > i {
        i, j = j, i
    }
    //log.Println(i, j, "=>", i, j-i-1)
    //return ctx.distMatrix[i][j-i-1]
    return ctx.distMatrix[i][j]
}

func (ctx Context) nearestTo(j int) int {
    for i := 0; i < ctx.N; i++ {
        // k is what i used to be before the optimization
        // k -- point index in the Points slice
        k := int(ctx.nearestToMatrix[j][i])
        if (k != j) && ctx.ps[k].active {
            return k
        }
    }

    return -1
}

// func (ctx Context) oldNearestTo(j int) int {
//     var nearest int = -1
//     var minDist float64 = math.Maxfloat64
//     for i := 0; i < ctx.N; i++ {
//         if (i == j) || (!ctx.ps[i].active) {
//             continue
//         } else if nearest == -1 {
//             nearest = i
//         } else {
//             d := ctx.dist(i, j)
//             if d < minDist {
//                 minDist = d
//                 nearest = i
//             }
//         }
//     }
//     //fmt.Println("nearest to", j, "is", nearest, "-", minDist)
//     return nearest
// }

func (ctx Context) nearestToExceptSmallerThan(j, a, b int, maxDist float64) int {
    var nearest int = -1
    var minDist float64 = math.MaxFloat64
    for i := 0; i < ctx.N; i++ {
        if (i == j) || (i == a) || (i == b) { //|| (!ps[i].active) {
            continue
        } else if nearest == -1 {
            nearest = i
        } else {
            d := ctx.dist(i, j)
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
    fmt.Printf("%f %d\n", solution.cost, 0)
    for i := 0; i < len(solution.order); i++ {
        fmt.Printf("%d ", solution.order[i])
    }
    fmt.Printf("\n")
}

func (ctx Context) setActive(val bool) {
    for i := 0; i < ctx.N; i++ {
        ctx.ps[i].active = val
    }
}

// solves the problem from the specified point
// enumerate all the points to get the best greedy solution
func (ctx Context) solveGreedyFrom(currentPoint int) Solution {
    nextPoint := 0
    var pathLen float64 = 0
    var pointOrder = make([]int, ctx.N)

    pointOrder[0] = currentPoint
    ctx.setActive(true)

    for i := 1; i < ctx.N; i++ {
        nextPoint = ctx.nearestTo(currentPoint)
        pointOrder[i] = nextPoint
        pathLen += ctx.dist(currentPoint, nextPoint)
        ctx.ps[currentPoint].active = false
        currentPoint = nextPoint
    }

    pathLen += ctx.dist(pointOrder[ctx.N-1], pointOrder[0])
    return Solution{pointOrder, pathLen}
}

// tries greedy alg for all the points in the graph and selects the best
func (ctx Context) solveGreedyBest() Solution {
    bestSolution := ctx.solveGreedyFrom(0)

    for i := 1; i < ctx.N; i++ {
        solution := ctx.solveGreedyFrom(i)
        if solution.cost < bestSolution.cost {
            bestSolution = solution
        }
    }

    return bestSolution
}

// randomly select two non-adjacent points
func (ctx Context) selectPoints(solution Solution) (int, int) {
    p1 := rand.Int() % ctx.N
    p3 := rand.Int() % ctx.N

    // p3 must not be near p1
    for p3 == p1 || p3 == (p1+1) % ctx.N || (p3+1) % ctx.N == p1 {
        p3 = rand.Int() % ctx.N
    }

    return p1, p3
}

// create new solution with swapped points and
// new cost recalculated from scratch (I used to
// have huge cumulative errors going from predictCost)
func (ctx Context) acceptSolution(p1, p3 int, solution Solution) Solution {
    acceptedSolution := reconnectPoints(p1, p3, solution)
    acceptedSolution.cost = ctx.calcCost(acceptedSolution, false)
    return acceptedSolution
}

// lightweight version of acceptSolution which
// only swaps points and sets predicted solution
// cost -- warning, this might contain huge cumulative
// error
func (ctx Context) acceptPredictedSolution(p1, p3 int, solution Solution) Solution {
    predictedCost := ctx.predictCost(p1, p3, solution)
    acceptedSolution := reconnectPoints(p1, p3, solution)
    acceptedSolution.cost = predictedCost
    return acceptedSolution
}

// run local search with Metropolis meta-heuristic
func (ctx Context) localSearch(currentSolution Solution, temperature float64) Solution {
    solution := cloneSolution(currentSolution)
    for k := 0; k < 3000; k++ {
        p1, p3 := ctx.selectPoints(solution)
        predictedCost := ctx.predictCost(p1, p3, solution)
        costDiff := predictedCost - solution.cost
        //log.Println(p1, p3, costDiff)

        if predictedCost <= solution.cost {
            //log.Println("taking predicted solution, costDiff", costDiff)
            //solution = reconnectPoints(p1, p3, solution)
            //solution.cost = predictedCost
            solution = ctx.acceptPredictedSolution(p1, p3, solution)
        } else {
            probability := math.Exp(- costDiff / temperature)
            //log.Println("prob", probability)

            if rand.Float64() < probability {
                //log.Println("taking bad solution", costDiff)
                //solution = reconnectPoints(p1, p3, solution)
                //solution.cost = predictedCost
                solution = ctx.acceptPredictedSolution(p1, p3, solution)
            }
        }
    }
    return solution
}

func (ctx Context) simulatedAnnealing() Solution {
    solution := ctx.solveGreedyFrom(0)
    bestSolution := solution
    t := 10.0
    alpha := 0.9999

    for k := 0; k < 30000; k++ {
        solution = ctx.localSearch(solution, t)
        if solution.cost < bestSolution.cost {
            log.Println("new solution", solution.cost)
            bestSolution = solution
        }
        t *= alpha
        log.Printf("t %f best cost %f\n", t, bestSolution.cost)
    }
    return bestSolution
}

func (ctx Context) calcCost(solution Solution, pr bool) float64 {
    cost := float64(0.0)
    N := len(solution.order)
    for i := 0; i < N; i++ {
        d := ctx.dist(solution.order[i], solution.order[(i+1) % N])
        if pr {
           log.Println(d)
        }
        cost += d
    }
    //cost += ctx.dist(solution.order[N-1], solution.order[0])
    return cost
}

func (ctx Context) predictCost(p1, p3 int, solution Solution) float64 {
    cost := solution.cost
    t1 := solution.order[p1 % ctx.N]
    t2 := solution.order[(p1+1) % ctx.N]
    t4 := solution.order[p3 % ctx.N]
    t3 := solution.order[(p3+1) % ctx.N]
    cost -= ctx.dist(t1, t2)
    cost -= ctx.dist(t4, t3)
    cost += ctx.dist(t1, t4)
    cost += ctx.dist(t2, t3)
    return cost
}

func cloneSolution(solution Solution) Solution {
    newSolution := solution
    newSolution.order = make([]int, len(solution.order))
    copy(newSolution.order, solution.order)
    return newSolution
}

func reconnectPoints(p1, p3 int, origSolution Solution) Solution {
    N := len(origSolution.order)

    solution := cloneSolution(origSolution)
    // solution := origSolution
    // solution.order = make([]int, N)
    // copy(solution.order, origSolution.order)

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

func (ctx Context) greedy2Opt(solution Solution) Solution {
    //log.Println("N", ctx.N)
    timestamp := time.Now().Unix()
    changed := true

    for changed {
        changed = false

        for i := 0; i < ctx.N; i++ {
            for j := i+2; j < ctx.N; j++ {
                predictedCost := ctx.predictCost(i, j, solution)
                if predictedCost < solution.cost {
                    solution = reconnectPoints(i, j, solution)

                    //diff := time.Now().Unix() - timestamp
                    //log.Println("swap", diff, "|", i, j, "|", solution.cost, "=>", predictedCost)
                    solution.cost = predictedCost

                    changed = true
                    timestamp = time.Now().Unix()
                    break
                }
            }

            if changed {
                break
            }

            if time.Now().Unix() - timestamp > MAX_SECONDS_BETWEEN_CHANGES {
                return solution
            }
        }
    }

    return solution
}

func (ctx Context) exhaustive2Opt(solution Solution) Solution {
    //log.Println("N", ctx.N)
    timestamp := time.Now().Unix()
    changed := true

    lastCost := float64(-1.0)

    for changed {
        changed = false

        bestI, bestJ := -1, -1
        var bestSwapCost float64 = -1.0

        for i := 0; i < ctx.N; i++ {
            for j := i+2; j < ctx.N; j++ {
                predictedCost := ctx.predictCost(i, j, solution)
                if predictedCost < solution.cost {
                    if bestSwapCost == -1 || predictedCost < bestSwapCost {
                        bestSwapCost = predictedCost
                        bestI, bestJ = i, j
                    }

                    changed = true
                    timestamp = time.Now().Unix()
                }
            }

            if time.Now().Unix() - timestamp > MAX_SECONDS_BETWEEN_CHANGES {
                return solution
            }
        }

        if changed {
            solution = reconnectPoints(bestI, bestJ, solution)
            //diff := time.Now().Unix() - timestamp
            //log.Println("swap", diff, "|", bestI, bestJ, "|", solution.cost, "=>", bestSwapCost)
            solution.cost = bestSwapCost

            if solution.cost < 20750.0 {
                return solution
            }

            if lastCost < 0 || (lastCost - solution.cost > 50.0) {
                //log.Println("current cost", solution.cost)
                lastCost = solution.cost
            }
        }
    }

    return solution
}

// TODO:
// 1.~select best of greedy solutions (try all points as a starting point)
// 2.+pre-compute distMatrix
// 3.+pre-compute nearestMatrix
// 4. implement
//    + 2-opt
//    _ k-opt
// 5. use double-linked slice instead of order list
// 6. use SA (Simulated Annealing)
// 7. use Metropolis meta-heuristics (to get out of local minima)
// 8. use tabu search
// 9.+implement cheaper way to predict cost after change
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
        ps[i] = Point{x, y, true}
    }

    ctx := Context{ps, nil, nil, len(ps)}
    ctx = ctx.init()

    switch {
    case alg == "greedy":
        solution := ctx.solveGreedyBest()
        printSolution(solution)

    case alg == "g2o":
        solution := ctx.solveGreedyFrom(0)
        //printSolution(solution)
        solution = ctx.greedy2Opt(solution)
        printSolution(solution)

    case alg == "e2o":
        solution := ctx.solveGreedyBest()
        //printSolution(solution)
        solution = ctx.exhaustive2Opt(solution)
        printSolution(solution)

    case alg == "g2oall":
        bestSolution := ctx.solveGreedyFrom(0)
        bestSolution = ctx.greedy2Opt(bestSolution)

        for i := 1; i < ctx.N; i++ {
            //printSolution(solution)
            solution := ctx.solveGreedyFrom(i)
            solution = ctx.greedy2Opt(solution)
            if solution.cost < bestSolution.cost {
                log.Printf("NEW BEST SOLUTION %f\n", solution.cost)
                bestSolution = solution
            }
            log.Println("iteration", i, "done")
        }
        printSolution(bestSolution)

    case alg == "g2oex":
        bestSolution := ctx.solveGreedyFrom(0)
        bestSolution = ctx.exhaustive2Opt(bestSolution)

        for i := 1; i < ctx.N; i++ {
            //printSolution(solution)
            solution := ctx.solveGreedyFrom(i)
            solution = ctx.exhaustive2Opt(solution)
            if solution.cost < bestSolution.cost {
                log.Printf("NEW BEST SOLUTION %f\n", solution.cost)
                bestSolution = solution
            }
            log.Println("iteration", i, "done")
        }
        printSolution(bestSolution)

    default:
        //solution := ctx.solveGreedyBest()
        //solution := ctx.solveGreedyFrom(90)
        //log.Println("greedy done")
        //printSolution(solution)

        //solution = ctx.exhaustive2Opt(solution)
        //solution = ctx.greedy2Opt(solution)
        //printSolution(solution)

        solution := ctx.simulatedAnnealing()
        printSolution(solution)
        log.Printf("actual cost %f\n", ctx.calcCost(solution, false))

        // solution := ctx.solveGreedyFrom(0)
        // p1, p3 := 5, 10
        // log.Println("points", p1, p3, solution.order[p1], solution.order[p3])
        // predictedCost := ctx.predictCost(p1, p3, solution)
        // newSolution := ctx.acceptSolution(p1, p3, solution)
        // printSolution(solution)
        // printSolution(newSolution)
        // log.Println("original", solution.cost)
        // log.Println("predicted", predictedCost)
        // log.Println("actual", newSolution.cost)
    }

    return 0
}

func main() {
    rand.Seed(time.Now().UTC().UnixNano())
    alg := "auto"
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    os.Exit(solveFile(os.Args[1], alg))
}
