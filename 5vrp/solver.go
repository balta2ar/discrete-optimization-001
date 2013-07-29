package main

import "time"
import "math"
//import "sort"
import "math/rand"
import "fmt"
import "log"
import "encoding/gob"
import "compress/gzip"
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

type Client struct {
    Demand int
    X, Y float32
}

type Context struct {
    N int // number of customers
    V int // number of vehicles
    C int // vehicle capacity

    Clients []Client
    // DistMatrix [][]float64
    // NearestToMatrix [][]int32
    // N int
}

type Path struct {
    Cost float32
    VertexOrder []int
}

type Solution struct {
    Cost float32
    Paths []Path
}

// --- Functions ---------------------------------------------------------------

// calc and cache distances from each to each point
// create triangle matrix to save space
// func (ctx Context) calcDistMatrix() Context {
//     ctx.DistMatrix = make([][]float64, ctx.N)
//     for i := 1; i < ctx.N; i++ {
//         // ctx.DistMatrix[i] = make([]float64, ctx.N)
//         // for j := 0; j < ctx.N; j++ {
//         //     ctx.DistMatrix[i][j] = ctx.calcDist(i, j)
//         // }
// 
//         ctx.DistMatrix[i] = make([]float64, i)
//         for j := 0; j < i; j++ {
//             ctx.DistMatrix[i][j] = ctx.calcDist(i, j)
//         }
//     }
//     return ctx
// }

// used to sort indexs by distance to some point
type IndexSorter struct {
    idx []int32
    from int32
    ctx Context
}

// func (s *IndexSorter) Len() int { return s.ctx.N }
// func (s *IndexSorter) Swap(i, j int) { s.idx[i], s.idx[j] = s.idx[j], s.idx[i] }
// func (s *IndexSorter) Less(i, j int) bool {
//     a := s.ctx.calcDist(int(s.from), int(s.idx[i]))
//     b := s.ctx.calcDist(int(s.from), int(s.idx[j]))
//     return a < b
// }

// calc and cache sorted list of distances from each to each point
// (to be used in greedy alg)
// WARNING: depends on calcDistMatrix
// ctx.DistMatrix MUST be calculated before calling this function
// func (ctx Context) calcNearestToMatrix() Context {
//     ctx.NearestToMatrix = make([][]int32, ctx.N)
//     for i := 0; i < ctx.N; i++ {
//         ctx.NearestToMatrix[i] = make([]int32, ctx.N)
//         for j := 0; j < ctx.N; j++ {
//             ctx.NearestToMatrix[i][j] = int32(j)
//         }
//         sort.Sort(&IndexSorter{ctx.NearestToMatrix[i], int32(i), ctx})
//     }
//     return ctx
// }

func (ctx Context) init() Context {
    //ctx = ctx.calcDistMatrix()
    //log.Println(ctx.DistMatrix)
    //ctx = ctx.calcNearestToMatrix()
    return ctx
}

func (ctx Context) calcDist(i, j int) float32 {
    return float32(math.Sqrt(math.Pow(float64(ctx.Clients[i].X - ctx.Clients[j].X), 2) +
                             math.Pow(float64(ctx.Clients[i].Y - ctx.Clients[j].Y), 2)))
}

func (ctx Context) dist(i, j int) float32 {
    // return ctx.DistMatrix[i][j]

    if i == j {
        return 0.0
    }
    if j > i {
        i, j = j, i
    }
    //log.Println(i, j, "=>", i, j-i-1)
    //return ctx.DistMatrix[i][j-i-1]
    return ctx.calcDist(i, j)
    //return ctx.DistMatrix[i][j]
}

// func (ctx Context) nearestTo(j int) int {
//     for i := 0; i < ctx.N; i++ {
//         // k is what i used to be before the optimization
//         // k -- point index in the Points slice
//         k := int(ctx.NearestToMatrix[j][i])
//         if (k != j) && ctx.Ps[k].Active {
//             return k
//         }
//     }
// 
//     return -1
// }

// func (ctx Context) oldNearestTo(j int) int {
//     var nearest int = -1
//     var minDist float64 = math.Maxfloat64
//     for i := 0; i < ctx.N; i++ {
//         if (i == j) || (!ctx.Ps[i].Active) {
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

// func (ctx Context) nearestToExceptSmallerThan(j, a, b int, maxDist float64) int {
//     var nearest int = -1
//     var minDist float64 = math.MaxFloat64
//     for i := 0; i < ctx.N; i++ {
//         if (i == j) || (i == a) || (i == b) { //|| (!Ps[i].Active) {
//             continue
//         } else if nearest == -1 {
//             nearest = i
//         } else {
//             d := ctx.dist(i, j)
//             if d >= maxDist {
//                 continue
//             }
//             if d < minDist {
//                 minDist = d
//                 nearest = i
//             }
//         }
//     }
//     //fmt.Println("nearest to", j, "is", nearest, "-", minDist)
//     return nearest
// }

func printSolution(solution Solution) {
    fmt.Printf("%f %d\n", solution.Cost, 0)
    for i := 0; i < len(solution.Paths); i++ {
        for j := 0; j < len(solution.Paths[i].VertexOrder); j++ {
            fmt.Printf("%d", solution.Paths[i].VertexOrder[j])
            if j != len(solution.Paths[i].VertexOrder)-1 {
                fmt.Printf(" ")
            }
        }
        fmt.Printf("\n")
    }
    //fmt.Printf("\n")
}

// func (ctx Context) setActive(val bool) {
//     for i := 0; i < ctx.N; i++ {
//         ctx.Ps[i].Active = val
//     }
// }

// tries greedy alg for all the points in the graph and selects the best
// func (ctx Context) solveGreedyBest() Solution {
//     bestSolution := ctx.solveGreedyFrom(0)
// 
//     for i := 1; i < ctx.N; i++ {
//         solution := ctx.solveGreedyFrom(i)
//         if solution.Cost < bestSolution.Cost {
//             bestSolution = solution
//         }
//     }
// 
//     return bestSolution
// }

// func (ctx Context) solveGreedyRandom() Solution {
//     return ctx.solveGreedyFrom(rand.Int() % ctx.N)
// }

// randomly select two non-adjacent points
// func (ctx Context) selectPoints(solution Solution) (int, int) {
//     p1 := rand.Int() % ctx.N
//     p3 := rand.Int() % ctx.N
// 
//     // p3 must not be near p1
//     for p3 == p1 || p3 == (p1+1) % ctx.N || (p3+1) % ctx.N == p1 {
//         p3 = rand.Int() % ctx.N
//     }
// 
//     return p1, p3
// }

// create new solution with swapped points and
// new cost recalculated from scratch (I used to
// have huge cumulative errors going from predictCost)
// func (ctx Context) acceptSolution(p1, p3 int, solution Solution) Solution {
//     acceptedSolution := reconnectPoints(p1, p3, solution)
//     acceptedSolution.Cost = ctx.calcCost(acceptedSolution, false)
//     return acceptedSolution
// }

// lightweight version of acceptSolution which
// only swaps points and sets predicted solution
// cost -- warning, this might contain huge cumulative
// error
// func (ctx Context) acceptPredictedSolution(p1, p3 int, solution Solution) Solution {
//     predictedCost := ctx.predictCost(p1, p3, solution)
//     acceptedSolution := reconnectPoints(p1, p3, solution)
//     acceptedSolution.Cost = predictedCost
//     return acceptedSolution
// }

// run local search with Metropolis meta-heuristic
// func (ctx Context) localSearch(currentSolution Solution, temperature float64) Solution {
//     solution := cloneSolution(currentSolution)
//     for k := 0; k < 10000; k++ {
//         p1, p3 := ctx.selectPoints(solution)
//         predictedCost := ctx.predictCost(p1, p3, solution)
//         costDiff := predictedCost - solution.Cost
//         //log.Println(p1, p3, costDiff)
// 
//         if predictedCost <= solution.Cost {
//             //log.Println("taking predicted solution, costDiff", costDiff)
//             //solution = reconnectPoints(p1, p3, solution)
//             //solution.Cost = predictedCost
//             solution = ctx.acceptPredictedSolution(p1, p3, solution)
//         } else {
//             probability := math.Exp(- costDiff / temperature)
//             //log.Println("prob", probability)
// 
//             if rand.Float64() < probability {
//                 //log.Println("taking bad solution", costDiff)
//                 //solution = reconnectPoints(p1, p3, solution)
//                 //solution.Cost = predictedCost
//                 solution = ctx.acceptPredictedSolution(p1, p3, solution)
//             }
//         }
//     }
//     return solution
// }
// 
// func (ctx Context) simulatedAnnealing() Solution {
//     //solution := ctx.solveGreedyFrom(0)
//     var solution Solution
//     ptr := loadSolution("solution.greedy.best.bin")
//     if ptr == nil {
//         //solution = ctx.solveGreedyRandom()
//         solution = ctx.solveGreedyBest()
//         saveSolution(&solution, "solution.greedy.best.bin")
//     } else {
//         solution = *ptr
//     }
// 
//     solution = *loadSolution("solution.last.bin")
// 
//     bestSolution := solution
//     t := 75.0
//     // 0.99991 -- 327K
//     alpha := 0.99999
// 
//     log.Println("start solution, t", t, "cost", solution.Cost)
//     //for k := 0; k < 200000; k++ {
//     for t > 1.5 {
//         if t < 40.0 {
//             alpha = 0.999999
//         }
// 
//         solution = ctx.localSearch(solution, t)
//         if solution.Cost < bestSolution.Cost {
//             diff := bestSolution.Cost - solution.Cost
//             log.Printf("1 | new solution, t %f cost %f diff %f\n", t, solution.Cost, diff)
//             bestSolution = solution
// 
//             saveSolution(&solution, "solution.current.bin")
//         }
//         t *= alpha
//         //log.Printf("t %f best cost %f\n", t, bestSolution.Cost)
//     }
//     log.Println("last solution, t", t, "cost", bestSolution.Cost)
// 
// 
//     // t = 50.0
//     // alpha = 0.99991
// 
//     // solution = bestSolution
//     // log.Println("start solution, t", t, "cost", solution.Cost)
//     // for k := 0; k < 30000; k++ {
//     //     solution = ctx.localSearch(solution, t)
//     //     if solution.Cost < bestSolution.Cost {
//     //         diff := bestSolution.Cost - solution.Cost
//     //         log.Println("2 | new solution, t", t, "cost", solution.Cost, "diff", diff)
//     //         bestSolution = solution
//     //     }
//     //     t *= alpha
//     //     //log.Printf("t %f best cost %f\n", t, bestSolution.Cost)
//     // }
//     // log.Println("last solution, t", t, "cost", bestSolution.Cost)
// 
//     return bestSolution
// }

// func (ctx Context) calcCost(solution Solution, pr bool) float64 {
//     cost := float64(0.0)
//     N := len(solution.Order)
//     for i := 0; i < N; i++ {
//         d := ctx.dist(solution.Order[i], solution.Order[(i+1) % N])
//         if pr {
//            log.Println(d)
//         }
//         cost += d
//     }
//     //cost += ctx.dist(solution.Order[N-1], solution.Order[0])
//     return cost
// }
// 
// func (ctx Context) predictCost(p1, p3 int, solution Solution) float64 {
//     cost := solution.Cost
//     t1 := solution.Order[p1 % ctx.N]
//     t2 := solution.Order[(p1+1) % ctx.N]
//     t4 := solution.Order[p3 % ctx.N]
//     t3 := solution.Order[(p3+1) % ctx.N]
//     cost -= ctx.dist(t1, t2)
//     cost -= ctx.dist(t4, t3)
//     cost += ctx.dist(t1, t4)
//     cost += ctx.dist(t2, t3)
//     return cost
// }
// 
// func cloneSolution(solution Solution) Solution {
//     newSolution := solution
//     newSolution.Order = make([]int, len(solution.Order))
//     copy(newSolution.Order, solution.Order)
//     return newSolution
// }
// 
// func reconnectPoints(p1, p3 int, origSolution Solution) Solution {
//     N := len(origSolution.Order)
// 
//     solution := cloneSolution(origSolution)
//     // solution := origSolution
//     // solution.Order = make([]int, N)
//     // copy(solution.Order, origSolution.Order)
// 
//     //t1 := solution.Order[p1]
//     t2 := solution.Order[(p1+1) % N]
// 
//     t3 := solution.Order[(p3+1) % N]
//     t4 := solution.Order[p3]
// 
//     // t3InOrder := findInSlice(t3, solution.Order)
//     // t3InOrderPrev := (t3InOrder-1) % N
//     // if t3InOrderPrev < 0 {
//     //     // stupid Go
//     //     t3InOrderPrev = N + t3InOrderPrev
//     // }
// 
//     //log.Println("t3InOrderPrev", t3InOrderPrev)
//     //t4 := solution.Order[t3InOrderPrev]
//     //log.Println("t3InOrder", t3InOrder, "t4", t4)
// 
//     t3InOrder := (p3+1) % N
//     //t3InOrderPrev := p3
// 
//     selected := p1
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
//         orderPart[i] = solution.Order[(from+i) % N]
//     }
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
//     solution.Order[ptr % N] = t4
//     ptr++
// 
//     // insert reversed part order
//     for i := 0; i < len(orderPart); i++ {
//         solution.Order[ptr % N] = orderPart[i]
//         ptr++
//     }
// 
//     // insert t2 => t3 connection
//     solution.Order[ptr % N] = t2
//     ptr++
//     solution.Order[ptr % N] = t3
// 
//     return solution
// }
// 
// func (ctx Context) greedy2Opt(solution Solution) Solution {
//     //log.Println("N", ctx.N)
//     timestamp := time.Now().Unix()
//     changed := true
// 
//     for changed {
//         changed = false
// 
//         for i := 0; i < ctx.N; i++ {
//             for j := i+2; j < ctx.N; j++ {
//                 predictedCost := ctx.predictCost(i, j, solution)
//                 if predictedCost < solution.Cost {
//                     solution = reconnectPoints(i, j, solution)
// 
//                     //diff := time.Now().Unix() - timestamp
//                     //log.Println("swap", diff, "|", i, j, "|", solution.Cost, "=>", predictedCost)
//                     solution.Cost = predictedCost
// 
//                     changed = true
//                     timestamp = time.Now().Unix()
//                     break
//                 }
//             }
// 
//             if changed {
//                 break
//             }
// 
//             if time.Now().Unix() - timestamp > MAX_SECONDS_BETWEEN_CHANGES {
//                 return solution
//             }
//         }
//     }
// 
//     return solution
// }
// 
// func (ctx Context) exhaustive2Opt(solution Solution) Solution {
//     //log.Println("N", ctx.N)
//     timestamp := time.Now().Unix()
//     changed := true
// 
//     lastCost := float64(-1.0)
// 
//     for changed {
//         changed = false
// 
//         bestI, bestJ := -1, -1
//         var bestSwapCost float64 = -1.0
// 
//         for i := 0; i < ctx.N; i++ {
//             for j := i+2; j < ctx.N; j++ {
//                 predictedCost := ctx.predictCost(i, j, solution)
//                 if predictedCost < solution.Cost {
//                     if bestSwapCost == -1 || predictedCost < bestSwapCost {
//                         bestSwapCost = predictedCost
//                         bestI, bestJ = i, j
//                     }
// 
//                     changed = true
//                     timestamp = time.Now().Unix()
//                 }
//             }
// 
//             if time.Now().Unix() - timestamp > MAX_SECONDS_BETWEEN_CHANGES {
//                 return solution
//             }
//         }
// 
//         if changed {
//             solution = reconnectPoints(bestI, bestJ, solution)
//             //diff := time.Now().Unix() - timestamp
//             //log.Println("swap", diff, "|", bestI, bestJ, "|", solution.Cost, "=>", bestSwapCost)
//             solution.Cost = bestSwapCost
// 
//             if solution.Cost < 20750.0 {
//                 return solution
//             }
// 
//             if lastCost < 0 || (lastCost - solution.Cost > 50.0) {
//                 //log.Println("current cost", solution.Cost)
//                 lastCost = solution.Cost
//             }
//         }
//     }
// 
//     return solution
// }

// - load / save ---------------------------------------------------------------

//
// Solution
//
func saveSolution(solution *Solution, name string) {
    file, err := os.Create(name)
    if err != nil {
        log.Println("Cannot save to file", name, err)
        return
    }
    defer file.Close()

    zip := gzip.NewWriter(file)
    defer zip.Close()

    encoder := gob.NewEncoder(zip)
    encoder.Encode(solution)
    //log.Println("Saved to file", name)
}

func loadSolution(name string) *Solution {
    file, err := os.Open(name)
    if err != nil {
        log.Println("Cannot open file", name, err)
        return nil
    }
    defer file.Close()

    unzip, _ := gzip.NewReader(file)
    defer unzip.Close()

    var solution Solution
    decoder := gob.NewDecoder(unzip)
    decoder.Decode(&solution)
    //log.Println("Loaded from file", name)
    return &solution
}

//
// Context
//
func saveContext(ctx *Context, name string) {
    // Ps Points
    // DistMatrix [][]float64
    // NearestToMatrix [][]int32
    // N int
    file, err := os.Create(name)
    if err != nil {
        log.Println("Cannot save to file", name, err)
        return
    }
    defer file.Close()

    zip := gzip.NewWriter(file)
    defer zip.Close()

    encoder := gob.NewEncoder(zip)
    encoder.Encode(ctx)
    log.Println("Saved to file", name)
}

func loadContext(name string) *Context {
    file, err := os.Open(name)
    if err != nil {
        log.Println("Cannot open file", name, err)
        return nil
    }
    defer file.Close()

    unzip, _ := gzip.NewReader(file)
    defer unzip.Close()

    var ctx Context
    decoder := gob.NewDecoder(unzip)
    decoder.Decode(&ctx)
    log.Println("Loaded from file", name)
    return &ctx
}

// TODO:

func initContextFromFile(filename string) Context {
    file, err := os.Open(filename)
    if err != nil {
        panic(fmt.Sprintf("Cannot open file %s: %s", filename, err))
    }
    defer file.Close()

    var N, V, C int
    var demand int
    var x, y float32
    fmt.Fscanf(file, "%d %d %d", &N, &V, &C)

    clients := []Client(make([]Client, N))

    for i := 0; i < N; i++ {
        fmt.Fscanf(file, "%d %f %f", &demand, &x, &y)
        clients[i] = Client{demand, x, y}
    }

    ctx := Context{N, V, C, clients}
    ctx = ctx.init()
    return ctx
}

func createContext(filename string) Context {
    // ctx := initContextFromFile(filename)
    // return ctx

    ctx := initContextFromFile(filename)
    // var ctx Context
    // ptr := loadContext("context.bin")
    // if ptr == nil {
    //     ctx = initContextFromFile(filename)
    //     saveContext(&ctx, "context.bin")
    // } else {
    //     ctx = *ptr
    // }
    return ctx
}

// Solution --------------------------------------------------------------------

func (ctx Context) nearestCustomer(v int, active []bool, capacityLeft int) int {
    bestCustomer := 0
    bestDist := float32(0)
    for i := 1; i < ctx.N; i++ {
        if !active[i] || capacityLeft < ctx.Clients[i].Demand {
            continue
        }
        dist := ctx.dist(v, i)
        if (bestCustomer == 0) || (dist < bestDist) {
            bestCustomer = i
            bestDist = dist
        }
    }
    return bestCustomer
}

// solves the problem from the specified point
// enumerate all the points to get the best greedy solution
func (ctx Context) solveGreedyFrom() Solution {
    paths := []Path{}
    totalCost := float32(0.0)
    activeClients := make([]bool, ctx.N)
    numVisitedCustomers := 0

    for i := 0; i < ctx.N; i++ {
        activeClients[i] = true
    }
    //log.Println(activeClients)

    // fill route for each vehicle while capacity remains
    for v := 0; v < ctx.V; v++ {
        //cost := float32(0.0)
        //route := []int{0}
        path := Path{float32(0), []int{0}}
        capacityLeft := ctx.C

        for {
            lastCustomer := path.VertexOrder[len(path.VertexOrder)-1]
            nextCustomer := ctx.nearestCustomer(lastCustomer,
                                                activeClients,
                                                capacityLeft)
            path.Cost += ctx.dist(lastCustomer, nextCustomer)
            path.VertexOrder = append(path.VertexOrder, nextCustomer)
            if nextCustomer == 0 { // warehouse returned => no customer can fit
                break
            } else {
                activeClients[nextCustomer] = false
                capacityLeft -= ctx.Clients[nextCustomer].Demand
                numVisitedCustomers += 1
            }
        }

        totalCost += path.Cost
        paths = append(paths, path)
    }

    return Solution{totalCost, paths}

    // nextPoint := 0
    // var pathLen float64 = 0
    // var pointOrder = make([]int, ctx.N)

    // pointOrder[0] = currentPoint
    // ctx.setActive(true)

    // for i := 1; i < ctx.N; i++ {
    //     nextPoint = ctx.nearestTo(currentPoint)
    //     pointOrder[i] = nextPoint
    //     pathLen += ctx.dist(currentPoint, nextPoint)
    //     ctx.Ps[currentPoint].Active = false
    //     currentPoint = nextPoint
    // }

    // pathLen += ctx.dist(pointOrder[ctx.N-1], pointOrder[0])
    // return Solution{pointOrder, pathLen}
}

func solveFile(filename string, alg string) int {
    ctx := createContext(filename)

    //log.Println(ctx)
    solution := ctx.solveGreedyFrom()
    printSolution(solution)

    // switch {
    // case alg == "greedy":
    //     solution := ctx.solveGreedyBest()
    //     printSolution(solution)

    // case alg == "g2o":
    //     solution := ctx.solveGreedyFrom(0)
    //     //printSolution(solution)
    //     solution = ctx.greedy2Opt(solution)
    //     printSolution(solution)

    // case alg == "e2o":
    //     solution := ctx.solveGreedyBest()
    //     //printSolution(solution)
    //     solution = ctx.exhaustive2Opt(solution)
    //     printSolution(solution)

    // case alg == "g2oall":
    //     bestSolution := ctx.solveGreedyFrom(0)
    //     bestSolution = ctx.greedy2Opt(bestSolution)

    //     for i := 1; i < ctx.N; i++ {
    //         //printSolution(solution)
    //         solution := ctx.solveGreedyFrom(i)
    //         solution = ctx.greedy2Opt(solution)
    //         if solution.Cost < bestSolution.Cost {
    //             log.Printf("NEW BEST SOLUTION %f\n", solution.Cost)
    //             bestSolution = solution
    //         }
    //         log.Println("iteration", i, "done")
    //     }
    //     printSolution(bestSolution)

    // case alg == "g2oex":
    //     bestSolution := ctx.solveGreedyFrom(0)
    //     bestSolution = ctx.exhaustive2Opt(bestSolution)

    //     for i := 1; i < ctx.N; i++ {
    //         //printSolution(solution)
    //         solution := ctx.solveGreedyFrom(i)
    //         solution = ctx.exhaustive2Opt(solution)
    //         if solution.Cost < bestSolution.Cost {
    //             log.Printf("NEW BEST SOLUTION %f\n", solution.Cost)
    //             bestSolution = solution
    //         }
    //         log.Println("iteration", i, "done")
    //     }
    //     printSolution(bestSolution)

    // default:
    //     //solution := ctx.solveGreedyBest()
    //     //solution := ctx.solveGreedyFrom(90)
    //     //log.Println("greedy done")
    //     //printSolution(solution)

    //     //solution = ctx.exhaustive2Opt(solution)
    //     //solution = ctx.greedy2Opt(solution)
    //     //printSolution(solution)

    //     solution := ctx.simulatedAnnealing()
    //     printSolution(solution)
    //     log.Printf("actual cost %f\n", ctx.calcCost(solution, false))

    //     // solution := ctx.solveGreedyFrom(0)
    //     // p1, p3 := 5, 10
    //     // log.Println("points", p1, p3, solution.Order[p1], solution.Order[p3])
    //     // predictedCost := ctx.predictCost(p1, p3, solution)
    //     // newSolution := ctx.acceptSolution(p1, p3, solution)
    //     // printSolution(solution)
    //     // printSolution(newSolution)
    //     // log.Println("original", solution.Cost)
    //     // log.Println("predicted", predictedCost)
    //     // log.Println("actual", newSolution.Cost)
    // }

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
