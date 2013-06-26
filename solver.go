package main

import "sort"
import "fmt"
import "os"
import "io"
import "bytes"
import "container/heap"
import "encoding/binary"
import "compress/gzip"

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

func dumpToFile(file *os.File, data []int32) int64 {
    // write data into bytes Buffer
    var buf bytes.Buffer
    binary.Write(&buf, binary.LittleEndian, data)

    // prepare packed buffer
    var packedBuf bytes.Buffer
    z := gzip.NewWriter(&packedBuf)

    // write unpacked to packed through gzip
    buf.WriteTo(z)
    z.Flush()
    z.Close()

    size := int64(packedBuf.Len())
    packedBuf.WriteTo(file)
    return size
}

func loadFromFile(file *os.File, size int, data *[]int32) {
    var packedIn bytes.Buffer
    packedIn.ReadFrom(io.LimitReader(file, int64(size)))

    unz, _ := gzip.NewReader(&packedIn)
    var unpacked bytes.Buffer
    unpacked.ReadFrom(unz)

    //var dataIn = make([]int32, N)
    binary.Read(&unpacked, binary.LittleEndian, data)
}

func solveDynamicProgramming(K int32, v []int32, w []int32) {
    var N int32 = int32(len(v))

    file, _ := os.Create("dptable.bin")
    var offsets = make([]int64, N+1) // store offsets of dumped columns in the file
    var sizes = make([]int, N+1)
    var position int64
    var lastPackedSize int64

    var O = make([][]int32, 2)
    var x = make([]byte, N)
    // create O lookup table

    var k int32
    var j int32
    var i int32

    for i = 0; i <= 1; i++ {
        O[i] = make([]int32, K+1)
        O[i][0] = 0
    }

    // reset item-in-use table
    for j = 0; j < N; j++ {
        x[j] = 0
    }
    O[0][0] = 0
    O[1][0] = 0
    // O(k,j) denotes the optimal solution to the
    // knapsack problem with capacity k and
    // items [1..j]

    lastPackedSize = dumpToFile(file, O[0])
    offsets[0] = position
    sizes[0] = int(lastPackedSize)
    position += lastPackedSize

    // for all items
    for j = 1; j <= N; j++ {
        // for all capacities
        for k = 1; k <= K; k++ {
            if w[j-1] <= k {
                O[1][k] = max(O[0][k], v[j-1] + O[0][k-w[j-1]])
            } else {
                O[1][k] = O[0][k]
            }
        }

        // dump to disk and save offset
        lastPackedSize = dumpToFile(file, O[1])
        offsets[j] = position
        sizes[j] = int(lastPackedSize)
        position += lastPackedSize

        for i := 0; i <= int(K); i++ {
            O[0][i] = O[1][i]
        }
    }

    file.Sync()
    fmt.Println(O[1][K], 1)

    // restore best set of items
    k = K
    file.Seek(offsets[N], 0)
    loadFromFile(file, sizes[N], &O[1])
    for i = N; i > 0; i-- {
        // preload first (previous) column
        file.Seek(offsets[i-1], 0)
        loadFromFile(file, sizes[i-1], &O[0])

        if O[1][k] != O[0][k] {
            x[i-1] = 1
            k -= w[i-1]
        }
    }

    // print best set
    for i = 0; i < N; i++ {
        if i == N-1 {
            fmt.Printf("%d", x[i])
        } else {
            fmt.Printf("%d ", x[i])
        }
    }
    fmt.Printf("\n")

    file.Close()
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

    switch {
    case alg == "estimate":
        fmt.Println("DP estimated memory usage, MB:",
                    (int(K+1) * int(n+1) * 4 + int(n)) / 1024 / 1024)
    case alg == "dp":
        solveDynamicProgramming(K, v, w)
    case alg == "bnb":
        solveBranchAndBound(K, v, w)
    default:
        solveBranchAndBound(K, v, w)
    }
}

func main() {
    alg := "auto"
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    solveFile(os.Args[1], alg)
}
