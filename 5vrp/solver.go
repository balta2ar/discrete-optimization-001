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

func max(a int, b int) (r int) {
    if a > b {
        return a
    } else {
        return b
    }
}

//
// Internal ADT
//

type Float float64
const (
    MaxFloat = math.MaxFloat64
)

type Client struct {
    Demand int
    X, Y Float
}

type Context struct {
    N int // number of customers
    V int // number of vehicles
    C int // vehicle capacity

    Clients []Client
    DistMatrix [][]Float
    // NearestToMatrix [][]int32
    // N int
}

const (
    // move customer from one route to another, (possibly the same route)
    // CustomerFrom -- customer pos in PathFrom to extract
    // CustomerTo -- customer after which to insert CustomerFrom to PathTo
    MT_MOVE = iota
    // swap two customers (possibly in the same route)
    MT_SWAP
    // reverse part of the same route
    // (range including both CustomerFrom, CustomerTo)
    MT_REVERSE
    // max value for convenience (to use in rand)
    MT_LAST
)

type CustomerMove struct {
    MoveType int
    PathFrom, CustomerFrom int
    PathTo, CustomerTo int
    NewCost Float
    NewOverCapacity int
    Feasible bool
}

type Path struct {
    Cost Float
    Demand int
    VertexOrder []int
}

type Solution struct {
    Cost Float
    OverCapacity int
    Paths []Path
}

// --- Functions ---------------------------------------------------------------

// calc and cache distances from each to each point
// create triangle matrix to save space
func (ctx Context) calcDistMatrix() Context {
    ctx.DistMatrix = make([][]Float, ctx.N)
    for i := 1; i < ctx.N; i++ {
        // ctx.DistMatrix[i] = make([]float64, ctx.N)
        // for j := 0; j < ctx.N; j++ {
        //     ctx.DistMatrix[i][j] = ctx.calcDist(i, j)
        // }

        ctx.DistMatrix[i] = make([]Float, i)
        for j := 0; j < i; j++ {
            ctx.DistMatrix[i][j] = ctx.calcDist(i, j)
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

func (ctx Context) calcDist(i, j int) Float {
    return Float(math.Sqrt(math.Pow(float64(ctx.Clients[i].X - ctx.Clients[j].X), 2) +
                             math.Pow(float64(ctx.Clients[i].Y - ctx.Clients[j].Y), 2)))
}

func (ctx Context) dist(i, j int) Float {
    // return ctx.DistMatrix[i][j]

    if i == j {
        return 0.0
    }
    if j > i {
        i, j = j, i
    }
    //log.Println(i, j, "=>", i, j-i-1)
    // return ctx.calcDist(i, j)
    return ctx.DistMatrix[i][j]
}

func (ctx Context) averageDist() float64 {
    N := 0
    D := 0.0

    for i := 0; i < ctx.N; i++ {
        for j := 0; j < ctx.N; j++ {
            if i != j {
                D += float64(ctx.dist(i, j))
                N += 1
            }
        }
    }

    return D/float64(N)
}

func (ctx Context) averageDistInSolution(solution Solution) float64 {
    N := 0
    D := 0.0

    for _, path := range solution.Paths {
        for i, _ := range path.VertexOrder {
            if i != len(path.VertexOrder)-1 {
                D += float64(ctx.dist(path.VertexOrder[i], path.VertexOrder[i+1]))
                N += 1
            }
        }
    }

    return D/float64(N)
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

func (ctx Context) printSolution(solution Solution) {
    ctx.printSolutionToFile(solution, "stdout")
}

func (ctx Context) printSolutionToFile(solution Solution, filename string) {
    var file *os.File
    var err error

    if filename == "stdout" {
        file = os.Stdout
    } else {
        file, err = os.Create(filename)
        if err != nil {
            log.Println("Cannot open file", filename, err)
            return
        }
        defer file.Close()
    }

    // log.Println("OverCapacity incr", solution.OverCapacity, "calc", ctx.overCapacity(solution))
    // log.Println("Cost incr", solution.Cost, "calc", ctx.solutionCost(solution))
    fmt.Fprintf(file, "%f %d\n", solution.Cost, 0)
    for i := 0; i < len(solution.Paths); i++ {
        for j := 0; j < len(solution.Paths[i].VertexOrder); j++ {
            fmt.Fprintf(file, "%d", solution.Paths[i].VertexOrder[j])
            if j != len(solution.Paths[i].VertexOrder)-1 {
                fmt.Fprintf(file, " ")
            }
        }
        fmt.Fprintf(file, "\n")
    }
    //fmt.Printf("\n")
}

// func (ctx Context) setActive(val bool) {
//     for i := 0; i < ctx.N; i++ {
//         ctx.Ps[i].Active = val
//     }
// }

//func randomMoveType() int { return rand.Int() % MT_LAST }
// func randomMoveType() int { return MT_SWAP }
func randomMoveType() int { return rand.Int() % MT_REVERSE }

// func (ctx Context) selectFeasibleMove(solution Solution) CustomerMove {
//     move := ctx.selectCustomerMove(randomMoveType(), solution)
//     for !ctx.isFeasibleMove(move, solution) {
//         move = ctx.selectCustomerMove(randomMoveType(), solution)
//     }
//     return move
// }

func (ctx Context) selectRandomPoints(solution Solution,
                                      pathFrom *int,
                                      customerFrom *int,
                                      pathTo *int,
                                      customerTo *int) {
    *pathFrom = rand.Int() % ctx.V
    *pathTo = rand.Int() % ctx.V
    *customerFrom = rand.Int() % len(solution.Paths[*pathFrom].VertexOrder)
    *customerTo = rand.Int() % len(solution.Paths[*pathTo].VertexOrder)
}

func (ctx Context) selectMTMove(solution Solution) CustomerMove {
    pathFrom := rand.Int() % ctx.V
    for len(solution.Paths[pathFrom].VertexOrder) < 3 {
        pathFrom = rand.Int() % ctx.V
    }

    vFrom := solution.Paths[pathFrom].VertexOrder
    customerFrom := (rand.Int() % (len(vFrom) - 2)) + 1

    pathTo := rand.Int() % ctx.V
    vTo := solution.Paths[pathTo].VertexOrder
    customerTo := rand.Int() % (len(vTo) - 1)
    for {
        if pathFrom != pathTo {
            break
        } else {
            if customerFrom != customerTo {
                break
            }
        }
        pathTo = rand.Int() % ctx.V
        vTo = solution.Paths[pathTo].VertexOrder
        customerTo = rand.Int() % (len(vTo) - 1)
    }

    return CustomerMove{MT_MOVE, pathFrom, customerFrom, pathTo, customerTo, 0, 0, false}
}

func (ctx Context) selectMTSwap(solution Solution) CustomerMove {
    pathFrom, pathTo, customerFrom, customerTo := 0, 0, 0, 0
    for {
        ctx.selectRandomPoints(solution, &pathFrom, &customerFrom, &pathTo, &customerTo)
        lenFrom := len(solution.Paths[pathFrom].VertexOrder)
        lenTo := len(solution.Paths[pathTo].VertexOrder)
        if ((pathFrom == pathTo) && (customerFrom == customerTo)) ||
            (customerFrom == 0) ||
            (customerFrom == lenFrom-1) ||
            (customerTo == 0) ||
            (customerTo == lenTo-1) {
            continue
        }
        break
    }

    return CustomerMove{MT_SWAP, pathFrom, customerFrom, pathTo, customerTo, 0, 0, false}
}

func (ctx Context) selectMTReverse(solution Solution) CustomerMove {
    return CustomerMove{MT_REVERSE, 0, 0, 0, 0, 0, 0, false}
}

// randomly select two customers
// return: route1, customer1, route2, customer2
func (ctx Context) selectCustomerMove(moveType int, solution Solution) CustomerMove {
    var move CustomerMove
    switch moveType {
    case MT_MOVE:
        move = ctx.selectMTMove(solution)
    case MT_SWAP:
        move = ctx.selectMTSwap(solution)
    case MT_REVERSE:
        move = ctx.selectMTReverse(solution)
    default:
        panic(fmt.Sprintf("Unknown move type %v", moveType))
    }

    // move := CustomerMove{moveType, pathFrom, customerFrom, pathTo, customerTo, 0, 0, false}
    // move := CustomerMove{1, 2, 0, 1, 0, 0, false}
    move.NewCost = ctx.costAfterMove(move, solution)
    move.NewOverCapacity = ctx.overCapacityAfterMove(move, solution)
    move.Feasible = move.NewOverCapacity == 0 //ctx.isFeasibleMove(move, solution)
    //return CustomerMove{2, 1, 2, 2, 0, true}
    //return CustomerMove{2, 1, 1, 5, 0, true}
    return move
}

func (ctx Context) pathDemand(path Path) int {
    demand := 0
    for i := 0; i < len(path.VertexOrder); i++ {
        demand += ctx.Clients[path.VertexOrder[i]].Demand
    }
    return demand
}

func (ctx Context) pathCost(path Path) Float {
    cost := Float(0)
    for i := 0; i < len(path.VertexOrder) - 1; i++ {
        cost += ctx.dist(path.VertexOrder[i], path.VertexOrder[i+1])
    }
    return cost
}

func (ctx Context) overCapacity(solution Solution) int {
    over := 0
    for _, path := range solution.Paths {
        current := max(0, ctx.pathDemand(path) - ctx.C)
        over += current
        // log.Println("overCapacity i", i, "demand incr", path.Demand,
        //             "calc", ctx.pathDemand(path),
        //             "current incr", max(0, path.Demand - ctx.C),
        //             "calc", current, "over", over)
    }
    return over
}

func (ctx Context) solutionCost(solution Solution) Float {
    cost := Float(0)
    for i := 0; i < len(solution.Paths); i++ {
        cost += ctx.pathCost(solution.Paths[i])
    }
    return cost
}

// func (ctx Context) isFeasibleMove(move CustomerMove, solution Solution) bool {
//     if move.PathFrom == move.PathTo {
//         return true
//     }
//     path := solution.Paths[move.PathTo]
//     newDemand := ctx.pathDemand(path)
//     newDemand += ctx.Clients[path.VertexOrder[move.CustomerTo]].Demand
//     // log.Println("newDemand", newDemand, "C", ctx.C)
//     return newDemand <= ctx.C
// }

func (ctx Context) isFeasibleSolution(solution Solution) bool {
    for _, path := range solution.Paths {
        if ctx.pathDemand(path) > ctx.C {
            return false
        }
    }

    return true

    // if move.PathFrom == move.PathTo {
    //     return true
    // }
    // path := solution.Paths[move.PathTo]
    // newDemand := ctx.pathDemand(path)
    // newDemand += ctx.Clients[path.VertexOrder[move.CustomerTo]].Demand
    // // log.Println("newDemand", newDemand, "C", ctx.C)
    // return newDemand < ctx.C
}

func (ctx Context) costAfterMTMove(move CustomerMove, solution Solution) Float {
    cost := solution.Cost

    // special cases
    if (move.PathFrom == move.PathTo) && (move.CustomerFrom+1 == move.CustomerTo) {
        // ..., CustomerFrom, CustomerTo, ...
        // swap CustomerFrom and CustomerTo
    }

    if (move.PathFrom == move.PathTo) && (move.CustomerTo+1 == move.CustomerFrom) {
        // ..., CustomerTo, CustomerFrom ...
        // do nothing, CustomerFrom will be inserted into the same position
        return cost
    }

    // remove "from" customer
    cFrom := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]
    cFromNext := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom+1]
    cFromPrev := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom-1]
    cost -= ctx.dist(cFrom, cFromNext)
    cost -= ctx.dist(cFrom, cFromPrev)
    cost += ctx.dist(cFromPrev, cFromNext)

    customerTo := move.CustomerTo
    customerToNext := customerTo + 1
    cToNext := solution.Paths[move.PathTo].VertexOrder[customerToNext]
    cToPrev := solution.Paths[move.PathTo].VertexOrder[customerTo]

    // add "to" customer
    cost -= ctx.dist(cToPrev, cToNext)
    cost += ctx.dist(cFrom, cToPrev)
    cost += ctx.dist(cFrom, cToNext)

    return cost
}

func (ctx Context) costAfterMTSwap(move CustomerMove, solution Solution) Float {
    cost := solution.Cost

    cFrom := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]
    cFromPrev := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom-1]
    cFromNext := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom+1]

    cTo := solution.Paths[move.PathTo].VertexOrder[move.CustomerTo]
    cToPrev := solution.Paths[move.PathTo].VertexOrder[move.CustomerTo-1]
    cToNext := solution.Paths[move.PathTo].VertexOrder[move.CustomerTo+1]

    cost -= ctx.dist(cFromPrev, cFrom)
    cost -= ctx.dist(cFrom, cFromNext)
    cost -= ctx.dist(cToPrev, cTo)
    cost -= ctx.dist(cTo, cToNext)

    solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom],
    solution.Paths[move.PathTo].VertexOrder[move.CustomerTo] =
        solution.Paths[move.PathTo].VertexOrder[move.CustomerTo],
        solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]

    cFromPrev = solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom-1]
    cFromNext = solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom+1]
    cToPrev = solution.Paths[move.PathTo].VertexOrder[move.CustomerTo-1]
    cToNext = solution.Paths[move.PathTo].VertexOrder[move.CustomerTo+1]

    cost += ctx.dist(cFromPrev, cTo)
    cost += ctx.dist(cTo, cFromNext)
    cost += ctx.dist(cToPrev, cFrom)
    cost += ctx.dist(cFrom, cToNext)

    solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom],
    solution.Paths[move.PathTo].VertexOrder[move.CustomerTo] =
        solution.Paths[move.PathTo].VertexOrder[move.CustomerTo],
        solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]

    return cost
}

func (ctx Context) costAfterMTReverse(move CustomerMove, solution Solution) Float {
    return 0
}

func (ctx Context) costAfterMove(move CustomerMove, solution Solution) Float {
    switch move.MoveType {
    case MT_MOVE:
        return ctx.costAfterMTMove(move, solution)
    case MT_SWAP:
        return ctx.costAfterMTSwap(move, solution)
    case MT_REVERSE:
        return ctx.costAfterMTReverse(move, solution)
    default:
        panic(fmt.Sprintf("Unknown move type %v", move.MoveType))
    }
}

func (ctx Context) printPathDemands(solution Solution) {
    for i := 0; i < len(solution.Paths); i++ {
        over := max(0, solution.Paths[i].Demand - ctx.C)
        log.Println("path i", i, "demand", solution.Paths[i].Demand, "over", over)
    }
}

// calc demands for From and To paths
func (ctx Context) pathDemandsAfterMTMove(move CustomerMove, solution Solution) (int, int) {
    pathFromDemand := solution.Paths[move.PathFrom].Demand
    pathToDemand := solution.Paths[move.PathTo].Demand
    if move.PathFrom == move.PathTo {
        return pathFromDemand, pathToDemand
    }
    customerId := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]
    customerDemand := ctx.Clients[customerId].Demand

    from, to := pathFromDemand - customerDemand, pathToDemand + customerDemand
    return from, to
}

func (ctx Context) pathDemandsAfterMTSwap(move CustomerMove, solution Solution) (int, int) {
    pathFromDemand := solution.Paths[move.PathFrom].Demand
    pathToDemand := solution.Paths[move.PathTo].Demand
    if move.PathFrom == move.PathTo {
        return pathFromDemand, pathToDemand
    }
    customerFrom := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]
    customerFromDemand := ctx.Clients[customerFrom].Demand
    customerTo := solution.Paths[move.PathTo].VertexOrder[move.CustomerTo]
    customerToDemand := ctx.Clients[customerTo].Demand

    from := pathFromDemand - customerFromDemand + customerToDemand
    to := pathToDemand - customerToDemand + customerFromDemand
    return from, to
}

func (ctx Context) pathDemandsAfterMTReverse(move CustomerMove, solution Solution) (int, int) {
    return 0, 0
}

func (ctx Context) overCapacityAfterMTMove(move CustomerMove, solution Solution) int {
    over := solution.OverCapacity

    if move.PathFrom == move.PathTo {
        return over
    }

    pathFromDemand := solution.Paths[move.PathFrom].Demand
    pathToDemand := solution.Paths[move.PathTo].Demand
    customerId := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]
    customerDemand := ctx.Clients[customerId].Demand

    oldOverFrom := max(0, pathFromDemand - ctx.C)
    newOverFrom := max(0, pathFromDemand - customerDemand - ctx.C)
    oldOverTo := max(0, pathToDemand - ctx.C)
    newOverTo := max(0, pathToDemand + customerDemand - ctx.C)

    over -= oldOverFrom
    over += newOverFrom
    over -= oldOverTo
    over += newOverTo

    return over
}

func (ctx Context) overCapacityAfterMTSwap(move CustomerMove, solution Solution) int {
    over := solution.OverCapacity

    if move.PathFrom == move.PathTo {
        return over
    }

    pathFromDemand := solution.Paths[move.PathFrom].Demand
    pathToDemand := solution.Paths[move.PathTo].Demand
    customerFrom := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]
    customerTo := solution.Paths[move.PathTo].VertexOrder[move.CustomerTo]
    customerFromDemand := ctx.Clients[customerFrom].Demand
    customerToDemand := ctx.Clients[customerTo].Demand

    oldOverFrom := max(0, pathFromDemand - ctx.C)
    newOverFrom := max(0, pathFromDemand - customerFromDemand + customerToDemand - ctx.C)
    oldOverTo := max(0, pathToDemand - ctx.C)
    newOverTo := max(0, pathToDemand - customerToDemand + customerFromDemand - ctx.C)

    over -= oldOverFrom
    over += newOverFrom
    over -= oldOverTo
    over += newOverTo

    return over
}

func (ctx Context) overCapacityAfterMTReverse(move CustomerMove, solution Solution) int {
    return 0
}

func (ctx Context) overCapacityAfterMove(move CustomerMove, solution Solution) int {
    switch move.MoveType {
    case MT_MOVE:
        return ctx.overCapacityAfterMTMove(move, solution)
    case MT_SWAP:
        return ctx.overCapacityAfterMTSwap(move, solution)
    case MT_REVERSE:
        return ctx.overCapacityAfterMTReverse(move, solution)
    default:
        panic(fmt.Sprintf("Unknown move type %v", move.MoveType))
    }
}

func (ctx Context) applyMTMove(move CustomerMove,
                               origSolution Solution) Solution {
    solution := cloneSolution(origSolution)

    from, to := ctx.pathDemandsAfterMTMove(move, solution)
    solution.Paths[move.PathFrom].Demand = from
    solution.Paths[move.PathTo].Demand = to
    // log.Println("applyMove", from, to)

    customerFromValue := solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]
    // log.Println("customerFromValue", customerFromValue)
    // log.Println("old pathFrom", solution.Paths[move.PathFrom].VertexOrder)

    // remove "from" customer
    pathFrom := make([]int, len(solution.Paths[move.PathFrom].VertexOrder) - 1)
    copy(pathFrom,
         solution.Paths[move.PathFrom].VertexOrder[:move.CustomerFrom])
    copy(pathFrom[move.CustomerFrom:],
         solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom+1:])
    solution.Paths[move.PathFrom].VertexOrder = pathFrom
    // log.Println("new pathFrom", pathFrom)

    // fix destination customer if on the same path
    customerTo := move.CustomerTo
    // log.Println("old customerTo", customerTo)
    if (move.PathFrom == move.PathTo) && (move.CustomerFrom < move.CustomerTo) {
        customerTo -= 1
    }
    // log.Println("new customerTo", customerTo)

    // log.Println("old pathTo", solution.Paths[move.PathTo].VertexOrder)
    // add "to" customer
    // log.Println("old pathTo len", len(solution.Paths[move.PathTo].VertexOrder))
    pathTo := make([]int, len(solution.Paths[move.PathTo].VertexOrder) + 1)
    copy(pathTo,
         solution.Paths[move.PathTo].VertexOrder[:customerTo+1])
    pathTo[customerTo+1] = customerFromValue
    copy(pathTo[customerTo+2:],
         solution.Paths[move.PathTo].VertexOrder[customerTo+1:])
    solution.Paths[move.PathTo].VertexOrder = pathTo
    // log.Println("new pathTo len", len(solution.Paths[move.PathTo].VertexOrder))
    // log.Println("new pathTo", pathTo)

    solution.Cost = move.NewCost
    solution.OverCapacity = move.NewOverCapacity
    return solution
}

func (ctx Context) applyMTSwap(move CustomerMove,
                               origSolution Solution) Solution {
    solution := cloneSolution(origSolution)

    from, to := ctx.pathDemandsAfterMTSwap(move, solution)
    solution.Paths[move.PathFrom].Demand = from
    solution.Paths[move.PathTo].Demand = to

    solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom],
    solution.Paths[move.PathTo].VertexOrder[move.CustomerTo] =
        solution.Paths[move.PathTo].VertexOrder[move.CustomerTo],
        solution.Paths[move.PathFrom].VertexOrder[move.CustomerFrom]

    solution.Cost = move.NewCost
    solution.OverCapacity = move.NewOverCapacity
    return solution
}

func (ctx Context) applyMTReverse(move CustomerMove,
                                  origSolution Solution) Solution {
    return origSolution
}

func (ctx Context) applyMove(move CustomerMove,
                             origSolution Solution) Solution {
    switch move.MoveType {
    case MT_MOVE:
        return ctx.applyMTMove(move, origSolution)
    case MT_SWAP:
        return ctx.applyMTSwap(move, origSolution)
    case MT_REVERSE:
        return ctx.applyMTReverse(move, origSolution)
    default:
        panic(fmt.Sprintf("Unknown move type %v", move.MoveType))
    }
}

func cloneSolution(solution Solution) Solution {
    newSolution := solution
    newSolution.Paths = make([]Path, len(solution.Paths))
    for i := 0; i < len(solution.Paths); i++ {
        newSolution.Paths[i] = solution.Paths[i]
        newSolution.Paths[i].VertexOrder = make([]int, len(newSolution.Paths[i].VertexOrder))
        copy(newSolution.Paths[i].VertexOrder, solution.Paths[i].VertexOrder)
    }
    return newSolution
}

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

func color(msg string, color int) string {
    return fmt.Sprintf("\033[%dm%s\033[0m", color, msg)
}

func green(msg string) string { return color(msg, 32) }
func red(msg string) string { return color(msg, 31) }

// run local search with Metropolis meta-heuristic
func (ctx Context) localSearch(solution Solution, temperature float64,
                               K int, penalty float64) Solution {
    // solution := cloneSolution(currentSolution)
    // log.Println("starting with solution cost", solution.Cost)
    for k := 0; k < K; k++ {
        //move := ctx.selectFeasibleMove(solution)
        move := ctx.selectCustomerMove(randomMoveType(), solution)
        // log.Println(move)
        predictedCost := move.NewCost
        demandExcess := move.NewOverCapacity //ctx.overCapacity(solution)
        predictedCost += Float(penalty * float64(demandExcess))
        costDiff := float64(predictedCost - solution.Cost)

        if predictedCost <= solution.Cost {
            log.Println(costDiff, "=", predictedCost, "-", solution.Cost)
            log.Println("taking predicted solution, costDiff", costDiff)
            solution = ctx.applyMove(move, solution)
            //calculatedCost := ctx.solutionCost(solution)
            //solution.Cost = calculatedCost
            // log.Println("cost after accept", solution.Cost, "calculatedCost", calculatedCost)
        } else {
            probability := math.Exp(- costDiff / temperature)
            //log.Println("prob", probability)

            if rand.Float64() < probability {
                // log.Println(demandExcess, solution.Cost, move.NewCost, predictedCost, costDiff)
                // log.Println(costDiff, "=", predictedCost, "-", solution.Cost)
                // log.Println("taking bad solution", probability)
                solution = ctx.applyMove(move, solution)
                // calculatedCost := ctx.solutionCost(solution)
                // if math.Abs(float64(solution.Cost - calculatedCost)) > 10.0 {
                //     log.Println("suspicious move", move)
                //     log.Println("cost after accept", solution.Cost, "calculatedCost", calculatedCost)
                // }
                // solution.Cost = calculatedCost
            }
        }
    }
    // log.Println("ending with solution cost", solution.Cost)
    return solution
}

func (ctx Context) simulatedAnnealing() Solution {
    //solution := ctx.solveGreedyFrom(0)
    // var solution Solution
    // ptr := loadSolution("solution.bin")
    // if ptr == nil {
    //     //solution = ctx.solveGreedyRandom()
    //     // solution = ctx.solveGreedyBest()
    //     solution = ctx.solveRandom()
    //     saveSolution(&solution, "solution.bin")
    // } else {
    //     solution = *ptr
    // }

    // solution = *loadSolution("solution.last.bin")

    // solution := ctx.solveGreedyBest()
    solution := ctx.solveRandom()
    ctx.printSolutionToFile(solution, "current.sol")
    // return solution
    // ctx.printSolution(solution)

    bestSolution := solution
    // t0 := ctx.averageDistInSolution(solution)
    t0 := 200.0
    minT := 1.0 //0.5
    K := 10000
    // 0.99991 -- 327K
    alpha := 0.99999
    penalty := 50.0

    tStep := 5.0
    oldT := t0

    log.Println("annealing params t0", t0, "K", K,
                "alpha", alpha, "penalty", penalty)

    // solution = ctx.localSearch(cloneSolution(solution), t, K, penalty)
    // return solution

    // reheat := 0

    //for k := 0; k < 200000; k++ {
    t := t0
    log.Println("start solution, t", t, "cost", solution.Cost)
    for t > minT {
        // if (reheat == 0) && (t < 1.0) {
        //     log.Println("reheat", reheat)
        //     reheat = 1
        //     t = 10.0
        // }

        // if t < 30.0 {
        //     alpha = 0.99995
        //     K = 15000
        //     //K = 10000 + int((6.0 - t) * 5000)
        // }
        // if t < 5.0 {
        //     alpha = 0.99999
        //     K = 20000
        // }
        // penalty = (1 + (t0 - t) / 5)
        penalty = (1 + (t0 - t) / 4)
        // K = 5000 + int(t0 - t) * 150

        solution = ctx.localSearch(cloneSolution(solution), t, K, penalty)
        // feasible := ctx.isFeasibleSolution(solution)
        feasible := solution.OverCapacity == 0
        if (solution.Cost < bestSolution.Cost) && feasible {
            ctx.printSolutionToFile(solution, "current.sol")
        // if (solution.Cost < bestSolution.Cost) {
            diff := bestSolution.Cost - solution.Cost
            msg := fmt.Sprintf("t %f cost %f penalty %f K %v diff %f",
                               t, solution.Cost, penalty, K, diff)
            if feasible {
                msg = green(msg)
            } else {
                msg = red(msg)
            }
            log.Printf("%s\n", msg)
            bestSolution = solution

            //saveSolution(&solution, "solution.current.bin")
            oldT = t
        }
        t *= alpha

        if (oldT - t) > tStep {
            log.Printf("t %f cost %f penalty %f K %v\n",
                       t, bestSolution.Cost, penalty, K)
            oldT = t
        }
    }
    log.Println("last solution, t", t, "cost", bestSolution.Cost)


    // t = 50.0
    // alpha = 0.99991

    // solution = bestSolution
    // log.Println("start solution, t", t, "cost", solution.Cost)
    // for k := 0; k < 30000; k++ {
    //     solution = ctx.localSearch(solution, t)
    //     if solution.Cost < bestSolution.Cost {
    //         diff := bestSolution.Cost - solution.Cost
    //         log.Println("2 | new solution, t", t, "cost", solution.Cost, "diff", diff)
    //         bestSolution = solution
    //     }
    //     t *= alpha
    //     //log.Printf("t %f best cost %f\n", t, bestSolution.Cost)
    // }
    // log.Println("last solution, t", t, "cost", bestSolution.Cost)

    return bestSolution
}

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
    var x, y Float
    fmt.Fscanf(file, "%d %d %d", &N, &V, &C)

    clients := []Client(make([]Client, N))

    for i := 0; i < N; i++ {
        fmt.Fscanf(file, "%d %f %f", &demand, &x, &y)
        clients[i] = Client{demand, x, y}
    }

    ctx := Context{N, V, C, clients, nil}
    ctx = ctx.init()
    return ctx
}

func createContext(filename string) Context {
    // ctx := initContextFromFile(filename)
    // return ctx

    ctx := initContextFromFile(filename)
    ctx = ctx.calcDistMatrix()

    // var ctx Context
    // ptr := loadContext("context.bin")
    // if ptr == nil {
    //     ctx = initContextFromFile(filename)
    //     ctx = ctx.calcDistMatrix()
    //     saveContext(&ctx, "context.bin")
    // } else {
    //     ctx = *ptr
    // }

    return ctx
}

// Solution --------------------------------------------------------------------

func (ctx Context) nearestCustomer(v int, active []bool, capacityLeft int) int {
    bestCustomer := 0
    bestDist := Float(0)
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

func (path *Path) appendCustomer(ctx *Context, lastCustomer, customer int) {
    //lastCustomer := path.VertexOrder[len(path.VertexOrder)-1]
    path.Cost += ctx.dist(lastCustomer, customer)
    path.Demand += ctx.Clients[customer].Demand
    path.VertexOrder = append(path.VertexOrder, customer)
}

// solves the problem from the specified point
// enumerate all the points to get the best greedy solution
func (ctx Context) solveGreedyFrom(startingCustomer int) Solution {
    paths := []Path{}
    totalCost := Float(0.0)
    activeCustomers := make([]bool, ctx.N)
    numVisitedCustomers := 0
    localStartingCustomer := startingCustomer

    for i := 0; i < ctx.N; i++ {
        activeCustomers[i] = true
    }
    //log.Println(activeCustomers)

    // lastCustomer := startingCustomer
    // path := Path{Float(0), []int{0}}
    // path.appendCustomer(&ctx, 0, lastCustomer)
    // capacityLeft := ctx.C - ctx.Clients[lastCustomer].Demand

    // fill route for each vehicle while capacity remains
    for v := 0; v < ctx.V; v++ {
        path := Path{Float(0), 0, []int{0}}
        capacityLeft := ctx.C

        for {
            if localStartingCustomer > 0 {
                path.appendCustomer(&ctx, 0, localStartingCustomer)
                capacityLeft -= ctx.Clients[localStartingCustomer].Demand
                activeCustomers[localStartingCustomer] = false
                localStartingCustomer = -1
                numVisitedCustomers += 1
            }

            lastCustomer := path.VertexOrder[len(path.VertexOrder)-1]
            nextCustomer := ctx.nearestCustomer(lastCustomer,
                                                activeCustomers,
                                                capacityLeft)
            path.appendCustomer(&ctx, lastCustomer, nextCustomer)
            if nextCustomer == 0 { // warehouse returned => no customer can fit
                break
            }

            activeCustomers[nextCustomer] = false
            capacityLeft -= ctx.Clients[nextCustomer].Demand
            numVisitedCustomers += 1
        }

        totalCost += path.Cost
        paths = append(paths, path)
    }

    if numVisitedCustomers < ctx.N-1 {
        log.Printf("Not all customers were visited: %d of %d (%d, %f)\n",
                    numVisitedCustomers, ctx.N-1, startingCustomer,
                    totalCost)
        return Solution{MaxFloat, 0, nil}
    } else {
        // log.Printf("Solution: %f (%d)\n", totalCost, startingCustomer)
    }
    return Solution{totalCost, 0, paths}
}

// just random customer assignement
func (ctx Context) solveRandom() Solution {
    solution := Solution{0.0, 0, make([]Path, ctx.V)}
    for i := 0; i < ctx.V; i++ {
        solution.Paths[i] = Path{Float(0), 0, []int{0}}
    }

    for i := 1; i < ctx.N; i++ {
        // assign customer to random path
        path := rand.Int() % ctx.V
        last := solution.Paths[path].VertexOrder[len(solution.Paths[path].VertexOrder)-1]
        solution.Paths[path].appendCustomer(&ctx, last, i)
    }

    // close path back to warehouse
    for i := 0; i < ctx.V; i++ {
        last := solution.Paths[i].VertexOrder[len(solution.Paths[i].VertexOrder)-1]
        solution.Paths[i].appendCustomer(&ctx, last, 0)
        solution.Cost += solution.Paths[i].Cost
        solution.OverCapacity += max(0, solution.Paths[i].Demand - ctx.C)
    }

    return solution
}

// tries greedy alg for all the points in the graph and selects the best
func (ctx Context) solveGreedyBest() Solution {
    bestSolution := ctx.solveGreedyFrom(1)

    for i := 2; i < ctx.N; i++ {
        solution := ctx.solveGreedyFrom(i)
        //log.Println("another solution", solution.Cost)
        if solution.Cost < bestSolution.Cost {
            bestSolution = solution
        }
    }

    return bestSolution
}

func (ctx Context) solveGreedyRandom() Solution {
    return ctx.solveGreedyFrom((rand.Int() % (ctx.N - 1)) + 1)
}

func solveFile(filename string, alg string) int {
    ctx := createContext(filename)

    solution := ctx.simulatedAnnealing()
    // calcCost := ctx.solutionCost(solution)
    // log.Println("actual cost", solution.Cost, "calc cost", calcCost)
    ctx.printSolution(solution)

    // //log.Println(ctx)
    // // solution := ctx.solveGreedyFrom(12)
    // solution := ctx.solveGreedyBest()
    // printSolution(solution)
    // // move := ctx.selectCustomerMove(solution)
    // move := ctx.selectFeasibleMove(solution)
    // // move = CustomerMove{0, 5, 2, 4, 0}
    // // log.Println("move", move, "feasible", ctx.isFeasibleMove(move, solution))
    // solution = ctx.applyMove(move, solution)
    // printSolution(solution)

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
    // log.Printf("\033[32m 1 | new solution, t cost diff feasible \033[0m\n")
    // return
    rand.Seed(time.Now().UTC().UnixNano())
    alg := "auto"
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    os.Exit(solveFile(os.Args[1], alg))
}
