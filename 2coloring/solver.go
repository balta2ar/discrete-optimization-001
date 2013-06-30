package main

import "sort"
import "fmt"
import "os"
//import "io"
//import "bytes"
//import "container/heap"
//import "encoding/binary"
//import "compress/gzip"

// functions which Go developers should have implemented but happened
// to be too lazy and religious to do so

func max(a int32, b int32) (r int32) {
    if a > b {
        return a
    } else {
        return b
    }
}

type Edge struct {
    u int32 // first vertex id
    v int32 // second vertex id
}

type Vertex struct {
    color int32
    E []int32 // list of connected edges
}

type Edges []Edge
type Vertices []Vertex

type Graph struct {
    E Edges
    V Vertices
}

// v -- global index of the vertex
// e -- local index of the edge in vertex edge list
// return global index of the other vertex
func (self *Graph) otherVertex(v int32, e int32) int32 {
    V := self.V[v]
    E := self.E[V.E[e]]
    if E.v == v {
        return E.u
    } else {
        return E.v
    }
}

type ByInt32 []int32
func (self ByInt32) Len() int { return len(self) }
func (self ByInt32) Less(i, j int) bool { return self[i] < self[j] }
func (self ByInt32) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

// type Node struct {
//     index int32   // index in the input data
//     value int32
//     weight int32
//     bound float32 // this is used as priority
//     selected byte
//     sel []byte
// }
// 
// // Priority queue -------------------------------------------------------------
// 
// type Items []Node
// 
// func (self Items) Len() int { return len(self) }
// func (self Items) Less(i, j int) bool { return self[i].bound < self[j].bound }
// func (self Items) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
// func (self *Items) Push(x interface{}) { *self = append(*self, x.(Node)) }
// func (self *Items) Pop() (popped interface{}) {
//     popped = (*self)[len(*self)-1]
//     *self = (*self)[:len(*self)-1]
//     return
// }
// 
// // Sorting --------------------------------------------------------------------
// 
// type ByValuePerWeight Items
// func (self ByValuePerWeight) Len() int { return len(self) }
// func (self ByValuePerWeight) Less(i, j int) bool {
//     a := float32(self[i].value) / float32(self[i].weight)
//     b := float32(self[j].value) / float32(self[j].weight)
//     return a > b
// }
// func (self ByValuePerWeight) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
// 
// type ByIndex Items
// func (self ByIndex) Len() int { return len(self) }
// func (self ByIndex) Less(i, j int) bool { return self[i].index < self[j].index }
// func (self ByIndex) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
// 
// // Branch and Bound -----------------------------------------------------------
// 
// func (node *Node) estimate(K int32, N int32, items Items) float32 {
//     var j, k int32
//     var totweight int32
//     var result float32
// 
//     if node.weight >= K {
//         return 0
//     }
// 
//     result = float32(node.value)
//     totweight = node.weight
//     j = node.index + 1
// 
//     for j < N && totweight + items[j].weight <= K {
//         totweight += items[j].weight
//         result += float32(items[j].value)
//         j++
//     }
// 
//     k = j
//     if k < N {
//         result += float32((K - totweight) * items[k].value / items[k].weight)
//     }
// 
//     return result
// }
// 
// // see
// // http://books.google.ru/books?id=QrvsNy9paOYC&pg=PA235&lpg=PA235&dq=knapsack+problem+branch+and+bound+C%2B%2B&source=bl&ots=e6ok2kODMN&sig=Yh5__d3iAFa5rEkaCoBJ2JAWybk&hl=en&sa=X&ei=k1EDULDrHIfKqgHqtYyxDA&redir_esc=y#v=onepage&q&f=true
// 
// func knapsackBranchAndBound(K int32, items Items, maxvalue *int32) []byte {
//     var N int32 = int32(len(items))
//     var u, v Node
//     //var x = make([]byte, N) // currently selected items
//     var bestset = make([]byte, N) // best selected items
//     pq := &Items{}
// 
//     heap.Init(pq)
//     *maxvalue = 0
// 
//     // initialize root
//     u = Node{0, 0, 0, 0, 0, make([]byte, N)}
//     // index = -1, start with fake root node
//     v = Node{-1, 0, 0, 0, 0, make([]byte, N)}
//     v.bound = v.estimate(K, N, items)
//     heap.Push(pq, v)
// 
//     for pq.Len() != 0 {
//         v = heap.Pop(pq).(Node)
//         if v.bound > float32(*maxvalue) {
//             // make child that includes the item
//             u = Node{v.index+1,
//                      v.value + items[v.index+1].value,
//                      v.weight + items[v.index+1].weight,
//                      0,
//                      0,
//                      make([]byte, N)}
// 
//             copy(u.sel, v.sel)
//             u.sel[u.index] = 1
// 
//             if u.weight <= K && u.value > *maxvalue {
//                 *maxvalue = u.value
//                 copy(bestset, u.sel)
//             }
//             u.bound = u.estimate(K, N, items)
//             if u.bound > float32(*maxvalue) {
//                 heap.Push(pq, u)
//             }
// 
//             // make child that does not include the item
//             u = Node{v.index+1,
//                      v.value,
//                      v.weight,
//                      0,
//                      0,
//                      make([]byte, N)}
//             u.bound = u.estimate(K, N, items)
//             copy(u.sel, v.sel)
//             u.sel[u.index] = 0
// 
//             if u.bound > float32(*maxvalue) {
//                 heap.Push(pq, u)
//             }
//         }
//     }
// 
//     return bestset
// }
// 
// func solveBranchAndBound(K int32, v []int32, w []int32) {
//     N := len(v)
//     items := make([]Node, N)
//     for i := 0; i < N; i++ {
//         items[i] = Node{int32(i), v[i], w[i], -1, 0, make([]byte, N)}
//     }
//     sort.Sort(ByValuePerWeight(items))
// 
//     var maxvalue int32 = -1
//     bestset := knapsackBranchAndBound(K, items, &maxvalue)
//     fmt.Println(maxvalue, 1) // not always optimal (1), actually
// 
//     // restore indexes
//     for i := 0; i < N; i++ {
//         items[i].selected = bestset[i]
//     }
//     sort.Sort(ByIndex(items))
//     for i := 0; i < N; i++ {
//         if items[i].selected == 1 {
//             fmt.Printf("1")
//         } else {
//             fmt.Printf("0")
//         }
//         if i != N-1 { // last?
//             fmt.Printf(" ")
//         }
//     }
//     fmt.Printf("\n")
// }
// 
// // Dynamic Programming --------------------------------------------------------
// 
// func dumpToFile(file *os.File, data []int32) int64 {
//     // write data into bytes Buffer
//     var buf bytes.Buffer
//     binary.Write(&buf, binary.LittleEndian, data)
// 
//     // prepare packed buffer
//     var packedBuf bytes.Buffer
//     z := gzip.NewWriter(&packedBuf)
// 
//     // write unpacked to packed through gzip
//     buf.WriteTo(z)
//     z.Flush()
//     z.Close()
// 
//     size := int64(packedBuf.Len())
//     packedBuf.WriteTo(file)
//     return size
// }
// 
// func loadFromFile(file *os.File, size int, data *[]int32) {
//     var packedIn bytes.Buffer
//     packedIn.ReadFrom(io.LimitReader(file, int64(size)))
// 
//     unz, _ := gzip.NewReader(&packedIn)
//     var unpacked bytes.Buffer
//     unpacked.ReadFrom(unz)
// 
//     //var dataIn = make([]int32, N)
//     binary.Read(&unpacked, binary.LittleEndian, data)
// }
// 
// func solveDynamicProgramming(K int32, v []int32, w []int32) {
//     var N int32 = int32(len(v))
// 
//     file, _ := os.Create("dptable.bin")
//     var offsets = make([]int64, N+1) // store offsets of dumped columns in the file
//     var sizes = make([]int, N+1)
//     var position int64
//     var lastPackedSize int64
// 
//     var O = make([][]int32, 2)
//     var x = make([]byte, N)
//     // create O lookup table
// 
//     var k int32
//     var j int32
//     var i int32
// 
//     for i = 0; i <= 1; i++ {
//         O[i] = make([]int32, K+1)
//         O[i][0] = 0
//     }
// 
//     // reset item-in-use table
//     for j = 0; j < N; j++ {
//         x[j] = 0
//     }
//     O[0][0] = 0
//     O[1][0] = 0
//     // O(k,j) denotes the optimal solution to the
//     // knapsack problem with capacity k and
//     // items [1..j]
// 
//     lastPackedSize = dumpToFile(file, O[0])
//     offsets[0] = position
//     sizes[0] = int(lastPackedSize)
//     position += lastPackedSize
// 
//     // for all items
//     for j = 1; j <= N; j++ {
//         // for all capacities
//         for k = 1; k <= K; k++ {
//             if w[j-1] <= k {
//                 O[1][k] = max(O[0][k], v[j-1] + O[0][k-w[j-1]])
//             } else {
//                 O[1][k] = O[0][k]
//             }
//         }
// 
//         // dump to disk and save offset
//         lastPackedSize = dumpToFile(file, O[1])
//         offsets[j] = position
//         sizes[j] = int(lastPackedSize)
//         position += lastPackedSize
// 
//         for i := 0; i <= int(K); i++ {
//             O[0][i] = O[1][i]
//         }
//     }
// 
//     file.Sync()
//     fmt.Println(O[1][K], 1)
// 
//     // restore best set of items
//     k = K
//     file.Seek(offsets[N], 0)
//     loadFromFile(file, sizes[N], &O[1])
//     for i = N; i > 0; i-- {
//         // preload first (previous) column
//         file.Seek(offsets[i-1], 0)
//         loadFromFile(file, sizes[i-1], &O[0])
// 
//         if O[1][k] != O[0][k] {
//             x[i-1] = 1
//             k -= w[i-1]
//         }
//     }
// 
//     // print best set
//     for i = 0; i < N; i++ {
//         if i == N-1 {
//             fmt.Printf("%d", x[i])
//         } else {
//             fmt.Printf("%d ", x[i])
//         }
//     }
//     fmt.Printf("\n")
// 
//     file.Close()
// }

func degree(g *Graph) int32 {
    var maxDegree int32 = 0
    for i := 0; i < len(g.V); i++ {
        maxDegree = max(maxDegree, int32(len(g.V[i].E)))
    }
    return maxDegree
}

func (g *Graph) chromaticNumber() int32 {
    var maxChNum int32 = 0
    for i := 0; i < len(g.V); i++ {
        maxChNum = max(maxChNum, int32(g.V[i].color))
    }
    return maxChNum
}

func (g *Graph) vertexNeighborColors(i int) []int32 {
    neibColors := make([]int32, 0)
    // get colors of all neighbors
    for j := 0; j < len(g.V[i].E); j++ {
        neibVertex := g.V[g.otherVertex(int32(i), int32(j))]
        neibColors = append(neibColors, neibVertex.color)
    }
    return neibColors
}

func minUnusedColor(colors *[]int32) int32 {
    sort.Sort(ByInt32(*colors))
    //fmt.Println(*colors)

    if len(*colors) == 1 {
        return (*colors)[0] + 1
    }

    for i := 0; i < len(*colors) - 1; i++ {
        if (*colors)[i+1] - (*colors)[i] > 1 {
            return (*colors)[i] + 1
        }
    }

    return (*colors)[len(*colors) - 1] + 1
}

func (g *Graph) assignVertexColor(i int) {
    neibColors := g.vertexNeighborColors(i)
    // find min unused color
    min_color := minUnusedColor(&neibColors)
    g.V[i].color = min_color
    //fmt.Println(min_color)
}

func (g *Graph) printColors() {
    for i := 0; i < len(g.V); i++ {
        if (i != len(g.V) - 1) {
            fmt.Printf("%d ", g.V[i].color - 1)
        } else {
            fmt.Printf("%d", g.V[i].color - 1)
        }
    }
    fmt.Printf("\n")
}

func (g *Graph) solveGreedySimple() {
    //NE := len(g.E)
    NV := len(g.V)
    //D := degree(&g)

    //fmt.Println(D)

    for i := 0; i < NV; i++ {
        g.assignVertexColor(i)
    }

    fmt.Println(g.chromaticNumber(), 0)
    g.printColors()
}

func solveFile(filename string, alg string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Cannot open file:", filename, err)
        return
    }
    defer file.Close()

    var NV, NE int32
    var i, v, u int32

    fmt.Fscanf(file, "%d %d", &NV, &NE)

    //v := make([]int32, n)
    E := make([]Edge, NE)
    V := make([]Vertex, NV)

    for i = 0; i < NV; i++ {
        V[i] = Vertex{0, make([]int32, 0)}
    }

    for i = 0; i < NE; i++ {
        fmt.Fscanf(file, "%d %d", &v, &u)
        E[i] = Edge{v, u}
        V[v].E = append(V[v].E, i)
        V[u].E = append(V[u].E, i)
    }

    g := Graph{E, V}

    g.solveGreedySimple()

    //fmt.Println(E)
    //fmt.Println(V)

    // switch {
    // case alg == "estimate":
    //     fmt.Println("DP estimated memory usage, MB:",
    //                 (int(K+1) * int(n+1) * 4 + int(n)) / 1024 / 1024)
    // case alg == "dp":
    //     solveDynamicProgramming(K, v, w)
    // case alg == "bnb":
    //     solveBranchAndBound(K, v, w)
    // default:
    //     solveBranchAndBound(K, v, w)
    // }
}

func main() {
    alg := "auto"
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    solveFile(os.Args[1], alg)

    /*
    c1 := []int32{1, 2, 3}
    fmt.Println(minUnusedColor(&c1))
    c2 := []int32{0, 2, 3}
    fmt.Println(minUnusedColor(&c2))
    c3 := []int32{1, 2, 3, 4, 6, 7, 8, 9}
    fmt.Println(minUnusedColor(&c3))
    */
}
