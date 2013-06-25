package main

import "sort"
import "fmt"
import "os"
import "container/heap"

// functions which Go developers should have implemented but happened
// to be too lazy and religious to do so

func max(a int32, b int32) (r int32) {
    if a > b {
        return a
    } else {
        return b
    }
}

type Node struct {
    index int32   // index in the input data
    value int32
    weight int32
    bound float32 // this is used as priority
    selected byte
    sel []byte
}

// Priority queue -------------------------------------------------------------

type Items []Node

func (self Items) Len() int { return len(self) }
func (self Items) Less(i, j int) bool { return self[i].bound < self[j].bound }
func (self Items) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
func (self *Items) Push(x interface{}) { *self = append(*self, x.(Node)) }
func (self *Items) Pop() (popped interface{}) {
    popped = (*self)[len(*self)-1]
    *self = (*self)[:len(*self)-1]
    return
}

// Sorting --------------------------------------------------------------------

type ByValuePerWeight Items
func (self ByValuePerWeight) Len() int { return len(self) }
func (self ByValuePerWeight) Less(i, j int) bool {
    a := float32(self[i].value) / float32(self[i].weight)
    b := float32(self[j].value) / float32(self[j].weight)
    return a > b
}
func (self ByValuePerWeight) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

type ByIndex Items
func (self ByIndex) Len() int { return len(self) }
func (self ByIndex) Less(i, j int) bool { return self[i].index < self[j].index }
func (self ByIndex) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

// Branch and Bound -----------------------------------------------------------

func (node *Node) estimate(K int32, N int32, items Items) float32 {
    var j, k int32
    var totweight int32
    var result float32

    if node.weight >= K {
        return 0
    }

    result = float32(node.value)
    totweight = node.weight
    j = node.index + 1

    for j < N && totweight + items[j].weight <= K {
        totweight += items[j].weight
        result += float32(items[j].value)
        j++
    }

    k = j
    if k < N {
        result += float32((K - totweight) * items[k].value / items[k].weight)
    }

    return result
}

// see
// http://books.google.ru/books?id=QrvsNy9paOYC&pg=PA235&lpg=PA235&dq=knapsack+problem+branch+and+bound+C%2B%2B&source=bl&ots=e6ok2kODMN&sig=Yh5__d3iAFa5rEkaCoBJ2JAWybk&hl=en&sa=X&ei=k1EDULDrHIfKqgHqtYyxDA&redir_esc=y#v=onepage&q&f=true

func knapsackBranchAndBound(K int32, items Items, maxvalue *int32) []byte {
    var N int32 = int32(len(items))
    var u, v Node
    //var x = make([]byte, N) // currently selected items
    var bestset = make([]byte, N) // best selected items
    pq := &Items{}

    heap.Init(pq)
    *maxvalue = 0

    // initialize root
    u = Node{0, 0, 0, 0, 0, make([]byte, N)}
    // index = -1, start with fake root node
    v = Node{-1, 0, 0, 0, 0, make([]byte, N)}
    v.bound = v.estimate(K, N, items)
    heap.Push(pq, v)

    for pq.Len() != 0 {
        v = heap.Pop(pq).(Node)
        if v.bound > float32(*maxvalue) {
            // make child that includes the item
            u = Node{v.index+1,
                     v.value + items[v.index+1].value,
                     v.weight + items[v.index+1].weight,
                     0,
                     0,
                     make([]byte, N)}

            copy(u.sel, v.sel)
            u.sel[u.index] = 1

            if u.weight <= K && u.value > *maxvalue {
                *maxvalue = u.value
                copy(bestset, u.sel)
            }
            u.bound = u.estimate(K, N, items)
            if u.bound > float32(*maxvalue) {
                heap.Push(pq, u)
            }

            // make child that does not include the item
            u = Node{v.index+1,
                     v.value,
                     v.weight,
                     0,
                     0,
                     make([]byte, N)}
            u.bound = u.estimate(K, N, items)
            copy(u.sel, v.sel)
            u.sel[u.index] = 0

            if u.bound > float32(*maxvalue) {
                heap.Push(pq, u)
            }
        }
    }

    return bestset
}

func solveBranchAndBound(K int32, v []int32, w []int32) {
    N := len(v)
    items := make([]Node, N)
    for i := 0; i < N; i++ {
        items[i] = Node{int32(i), v[i], w[i], -1, 0, make([]byte, N)}
    }
    sort.Sort(ByValuePerWeight(items))

    var maxvalue int32 = -1
    bestset := knapsackBranchAndBound(K, items, &maxvalue)

    fmt.Println(maxvalue, 1) // not always optimal (1), actually

    // restore indexes
    for i := 0; i < N; i++ {
        items[i].selected = bestset[i]
    }
    sort.Sort(ByIndex(items))
    for i := 0; i < N; i++ {
        if items[i].selected == 1 {
            fmt.Printf("1")
        } else {
            fmt.Printf("0")
        }
        if i != N-1 { // last?
            fmt.Printf(" ")
        }
    }
    fmt.Printf("\n")
}

// Dynamic Programming --------------------------------------------------------

func solveDynamicProgramming(K int32, v []int32, w []int32) {
    var N int32 = int32(len(v))

    var O = make([][]int32, K+1)
    var x = make([]byte, N)
    // create O lookup table

    var k int32
    var j int32
    var i int32

    for k = 0; k <= K; k++ {
        O[k] = make([]int32, N+1)
        O[k][0] = 0
    }
    // reset item-in-use table
    for j = 0; j < N; j++ {
        x[j] = 0
        O[0][j] = 0
    }
    // O(k,j) denotes the optimal solution to the
    // knapsack problem with capacity k and
    // items [1..j]

    // for all capacities
    for k = 1; k <= K; k++ {
        // for all items
        for j = 1; j <= N; j++ {
            if w[j-1] <= k {
                O[k][j] = max(O[k][j-1], v[j-1] + O[k-w[j-1]][j-1])
            } else {
                O[k][j] = O[k][j-1]
            }
        }
    }

    fmt.Println(O[K][N], 1)
    k = K
    for i = N; i > 0; i-- {
        if O[k][i] != O[k][i-1] {
            x[i-1] = 1
            k -= w[i-1]
        }
    }
    for i = 0; i < N; i++ {
        if i == N-1 {
            fmt.Printf("%d", x[i])
        } else {
            fmt.Printf("%d ", x[i])
        }
    }
    fmt.Printf("\n")
}

func solveFile(filename string, alg string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Cannot open file:", filename, err)
        return
    }
    defer file.Close()

    var n int32
    var K int32
    var i int32

    fmt.Fscanf(file, "%d %d", &n, &K)

    v := make([]int32, n)
    w := make([]int32, n)

    for i = 0; i < n; i++ {
        fmt.Fscanf(file, "%d %d", &v[i], &w[i])
    }

    fmt.Println("DP estimated memory usage, MB:", (int(K+1) * int(n+1) * 4 + int(n)) / 1024 / 1024)
    return
    //fmt.Println(n, K)
    //fmt.Println(v, w)
    switch {
    case alg == "dp":
        solveDynamicProgramming(K, v, w)
    case alg == "bnb":
        solveBranchAndBound(K, v, w)
    default:
        //fmt.Println("Auto select algorithm (BnB for now)")
        solveBranchAndBound(K, v, w)
    }
}

func main() {
    //fmt.Println("Solving file", os.Args[1])
    alg := "auto"
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    solveFile(os.Args[1], alg)
}
