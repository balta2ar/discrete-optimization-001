package main

import "strconv"
import "sort"
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

type Edge struct {
    u int32 // first vertex id
    v int32 // second vertex id
}

type Vertex struct {
    index int32 // original index
    color int32 // vexter color
    E []int32   // list of connected edges
}

type Edges []Edge
type Vertices []Vertex

// graph contains of edges and vertices
type Graph struct {
    E Edges
    V Vertices
}

type VarHeuristic int
type ValHeuristic int

const (
    VAR_BRUTE VarHeuristic = iota
    VAR_MRV
    VAR_MCV
)

const (
    VAL_BRUTE ValHeuristic = iota
    VAL_LCV
)

// additional information (besides the graph), required for CSP
type CSPContext struct {
    g *Graph
    //domain [][]int32 // possible values (colors) for each variable (vertex)
    domains []map[int32]bool // possible values (colors) for each variable (vertex)
    numColors int32  // target number of colors (does not change)
    currentUnassignedVertex int // current vertex in recursive solution calls
    varHeuristic VarHeuristic
    valHeuristic ValHeuristic
}

// save vertex order without reordering graph vertices
type VertexOrder struct {
    g *Graph
    order []int32
}

type LCVColorPair [2]int32

type ByLCVColor []LCVColorPair
func (self ByLCVColor) Len() int { return len(self) }
func (self ByLCVColor) Less(i, j int) bool { return self[i][1] < self[j][1] }
func (self ByLCVColor) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

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

type ByIndex VertexOrder
func (self ByIndex) Len() int { return len(self.order) }
func (self ByIndex) Less(i, j int) bool { return self.order[i] < self.order[j] }
func (self ByIndex) Swap(i, j int) { self.order[i], self.order[j] = self.order[j], self.order[i] }

type ByDegree VertexOrder
func (self ByDegree) Len() int { return len(self.order) }
func (self ByDegree) Less(i, j int) bool { return len(self.g.V[self.order[i]].E) < len(self.g.V[self.order[j]].E) }
func (self ByDegree) Swap(i, j int) { self.order[i], self.order[j] = self.order[j], self.order[i] }

// type ByDegree VertexOrder
// func (self ByDegree) Len() int { return len(self) }
// func (self ByDegree) Less(i, j int) bool { return len(self[i].E) < len(self[j].E) }
// func (self ByDegree) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

func (g *Graph) NV() int { return len(g.V) }
func (g *Graph) NE() int { return len(g.E) }

func (g *Graph) degree() int32 {
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

func (g *Graph) vertexNeighborColors(i int32) []int32 {
    neibColors := make([]int32, 0)
    // get colors of all neighbors
    for j := 0; j < len(g.V[i].E); j++ {
        neibVertex := g.V[g.otherVertex(i, int32(j))]
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

func (g *Graph) assignVertexColor(i int32) {
    neibColors := g.vertexNeighborColors(i)
    // find min unused color
    min_color := minUnusedColor(&neibColors)
    g.V[i].color = min_color
    //fmt.Println(min_color)
}

func (g *Graph) printColors() {
    for i := 0; i < len(g.V); i++ {
        if (i != len(g.V) - 1) {
            //fmt.Printf("%d (%d) ", g.V[i].color - 1, g.V[i].index)
            fmt.Printf("%d ", g.V[i].color - 1)
        } else {
            fmt.Printf("%d", g.V[i].color - 1)
        }
    }
    fmt.Printf("\n")
}

func (g *Graph) printSolution() {
    fmt.Println(g.chromaticNumber(), 0)
    g.printColors()
}

// greedy approach
func (g *Graph) solveGreedySimple() {
    //NE := len(g.E)
    NV := len(g.V)
    //D := degree(&g)

    ord := make([]int32, NV)
    for i := 0; i < NV; i++ {
        ord[i] = int32(i)
    }
    vertexOrder := VertexOrder{g, ord}

    sort.Sort(sort.Reverse(ByDegree(vertexOrder)))

    //sort.Sort(sort.Reverse(ByDegree(g.V)))
    //sort.Sort(ByDegree(g.V))
    //sort.Reverse(ByDegree(g.V))

    //fmt.Println(D)

    for i := 0; i < NV; i++ {
        g.assignVertexColor(vertexOrder.order[i])
    }

    //sort.Sort(ByIndex(g.V))

    //fmt.Println(g.chromaticNumber(), 0)
    //g.printColors()
    g.printSolution()
}

//
// CSP
//

func (c *CSPContext) init(nColors int) {
    c.domains = make([]map[int32]bool, c.g.NV())
    for i := 0; i < c.g.NV(); i++ {
        c.domains[i] = make(map[int32]bool)
        for j := 0; j < nColors; j++ {
            c.domains[i][int32(j + 1)] = true
        }
    }
    //fmt.Println(c.domains)
    c.currentUnassignedVertex = 0
}

func (v Vertex) numSameColorNeighbors(g *Graph, color int32) int {
    num := 0
    // check all neighbor vertices
    for i := 0; i < len(v.E); i++ {
        otherVertexIndex := g.otherVertex(v.index, int32(i))
        if g.V[otherVertexIndex].color == color {
            num += 1
        }
    }
    return num
}

// check if vertex is valid
func (v Vertex) valid(g *Graph) bool {
    // unassigned color?
    if v.color == 0 {
        return false
    }

    return v.numSameColorNeighbors(g, v.color) == 0
}

// check if graph is valid
func (g *Graph) valid() bool {
    // check if all vertices are valid
    for i := 0; i < g.NV(); i++ {
        if !g.V[i].valid(g) {
            return false
        }
    }

    return true
}

func (c *CSPContext) forwardCheckVertexColor(vertex int32, color int32) {
        for j := 0; j < len(c.g.V[vertex].E); j++ {
        neibVertexIndex := c.g.otherVertex(vertex, int32(j))
        delete(c.domains[neibVertexIndex], color)
        //neibColors = append(neibColors, neibVertex.color)
    }
}

// select Minimum Remaining Values vertex
func (c *CSPContext) getMRVVertex() int32 {
    // there must be unset vertices
    if c.currentUnassignedVertex >= c.g.NV() {
        panic("Call to getMRVVertex with no unset variables")
        return -1
    }

    var vertex int32 = -1

    // scan all domains, find the smallest one
    // (number of vertices == number of domains)
    for i := 0; i < c.g.NV(); i++ {
        // only check unset vertices
        if c.g.V[i].color != 0 {
            continue
        }

        // assign vertex if not yet assigned
        if vertex == -1 {
            vertex = int32(i)
        }

        if len(c.domains[i]) < len(c.domains[vertex]) {
            vertex = int32(i)
        }
    }

    if vertex == -1 {
        panic("getMRVVertex: Could not find the vertex")
        return -1
    }
    return vertex;
}

// select Least Constraining Value color
func (c *CSPContext) getLCVColor(vertex int32) int32 {
    if len(c.domains[vertex]) == 0 {
        panic(fmt.Sprintf("getLCVColor: No colors for vertex %d\n", vertex))
        return -1
    }

    var lcvColor int32 = int32(-1)
    lcvValue := -1
    for color, _ := range c.domains[vertex] {
        lcv := c.g.V[vertex].numSameColorNeighbors(c.g, color)

        if (lcvColor == -1) || (lcv < lcvValue) {
            lcvColor = color
            lcvValue = lcv
        }
    }

    return lcvColor
}

// return visit order for colors according to LCV heuristic
func (c *CSPContext) getLCVColorOrder(vertex int32) []LCVColorPair {
    ND := len(c.domains[vertex])
    pairs := make([]LCVColorPair, ND)

    i := 0
    for color, _ := range c.domains[vertex] {
        // color
        pairs[i][0] = color
        // cv value
        pairs[i][1] = int32(c.g.V[vertex].numSameColorNeighbors(c.g, color))
        i += 1
    }
    //fmt.Println(pairs)
    sort.Sort(ByLCVColor(pairs))
    return pairs

    // colors := make([]int32, ND)
    // for i := 0; i < ND; i++ {
    //     colors[i] = pairs[i][0]
    // }

    // return colors
}

func (c *CSPContext) solve(indent int) bool {
    // all vars assigned?
    if c.currentUnassignedVertex >= c.g.NV() {
        return c.g.valid()
    }

    // TODO: support MRV (minimum remaining values):
    // 1.+use maps instead of lists for domains (faster deletion)
    // 2.+forward check domain changes to neighbors after assigning
    //    color to the vertex
    // 3.+select MRV vertex (scan all vertices and select min)
    // 4.+try LCV (for values)
    // 5. try constraint propagation (stronger version of forward checking)

    // select var
    vertex := c.getMRVVertex() //c.currentUnassignedVertex
    //fmt.Println(indent, "Selected vertex", vertex)

    // no more values to try?
    if len(c.domains[vertex]) == 0 {
        //fmt.Println(indent, "Selected vertex", vertex, "is empty")
        return false
    }

    c.currentUnassignedVertex += 1

    // save a copy of state of all current domains
    savedDomains := pushDomains(c.domains)

    // now enumerate colors of the vertex
    if c.valHeuristic == VAL_BRUTE {
        for color, _ := range c.domains[vertex] {
            // set another color
            c.g.V[vertex].color = color

            // propagate color. this will change current domains state
            c.forwardCheckVertexColor(int32(vertex), color)

            //fmt.Println(c.g)
            //c.g.printSolution()
            if c.solve(indent + 1) {
                return true
            }

            // restore domains state to previous
            popDomains(&c.domains, &savedDomains)
        }
    } else if c.valHeuristic == VAL_LCV {
        lcvPairs := c.getLCVColorOrder(vertex)
        //fmt.Println(lcvPairs)
        //return false

        for color, _ := range lcvPairs {
            //color, _ := c.domains[vertex][order[i][0]]

            c.g.V[vertex].color = int32(color)
            c.forwardCheckVertexColor(int32(vertex), int32(color))
            if c.solve(indent + 1) {
                return true
            }
            popDomains(&c.domains, &savedDomains)
        }
    }

    /*
    for color := 0; color < int(c.numColors); color++ {
        c.g.V[vertex].color = c.domain[vertex][int32(color)]
        //fmt.Println(c.g)
        //c.g.printSolution()
        if c.solve() {
            return true
        }
    }
    */

    // unselect var
    c.currentUnassignedVertex -= 1
    c.g.V[vertex].color = 0

    return false
}

// contraint-satisfaction approach
func (g *Graph) solveCSP(nColors int32) int {
    //fmt.Println("Solving for", nColors, "colors")

    csp := CSPContext{g, nil, nColors, 0, VAR_BRUTE, VAL_BRUTE}
    csp.init(int(nColors))
    if csp.solve(0) {
        g.printSolution()
        return 0
    }
    return 1
}

func solveFile(filename string, alg string, nColors int32) int {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Println("Cannot open file:", filename, err)
        return 2
    }
    defer file.Close()

    var NV, NE int32
    var i, v, u int32

    fmt.Fscanf(file, "%d %d", &NV, &NE)

    //v := make([]int32, n)
    E := make([]Edge, NE)
    V := make([]Vertex, NV)

    for i = 0; i < NV; i++ {
        V[i] = Vertex{int32(i), 0, make([]int32, 0)}
    }

    for i = 0; i < NE; i++ {
        fmt.Fscanf(file, "%d %d", &v, &u)
        E[i] = Edge{v, u}
        V[v].E = append(V[v].E, i)
        V[u].E = append(V[u].E, i)
    }

    g := Graph{E, V}

    if nColors == -1 {
        nColors = g.degree() + 1
    }

    switch {
    case alg == "estimate":
        fmt.Println("Estimation is not implemented in this assignment")
        //fmt.Println("DP estimated memory usage, MB:",
        //            (int(K+1) * int(n+1) * 4 + int(n)) / 1024 / 1024)
    case alg == "greedy":
        g.solveGreedySimple()
    case alg == "csp":
        return g.solveCSP(nColors)
    default:
        return g.solveCSP(nColors)
    }

    return 0
}

func pushDomains(domains []map[int32]bool) []map[int32]bool {
    newDomains := make([]map[int32]bool, len(domains))
    for i := 0; i < len(domains); i++ {
        newDomains[i] = make(map[int32]bool)
        for k, v := range domains[i] {
            newDomains[i][k] = v
        }
    }
    return newDomains
}

func popDomains(dst *[]map[int32]bool, src *[]map[int32]bool) {
    newDomains := pushDomains(*src)
    *dst = newDomains
}

func test(alg string) {
    N := 2
    d := make([]map[int32]bool, N)

    d[0] = make(map[int32]bool)
    d[1] = make(map[int32]bool)

    d[0][1] = true
    d[0][2] = true
    d[0][3] = true

    d[1][2] = true

    saved := pushDomains(d)

    d[1][1] = true
    d[1][3] = true

    fmt.Println(d)
    fmt.Println(saved)

    popDomains(&d, &saved)

    fmt.Println(d)
    fmt.Println(len(d[0]))
}

func main() {
    alg := "auto"
    nColors := -1
    if len(os.Args) > 2 {
        alg = os.Args[2]
    }
    if len(os.Args) > 3 {
        nColors, _ = strconv.Atoi(os.Args[3])
    }
    os.Exit(solveFile(os.Args[1], alg, int32(nColors)))
    //test(alg)

}
