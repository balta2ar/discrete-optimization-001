package main

import "fmt"
import "os"

func max(a int32, b int32) (r int32) {
    if a > b {
        return a
    } else {
        return b
    }
}

// TODO: PriorityQueue

type Node struct {
    index int32
    value int32
    weight int32
    bound float32
}

func (node *Node) estimate(K int32, N int32, all []*Node) float32 {
    var j, k int32
    var totweight int32
    var result float32

    if node.weight >= K {
        return 0
    }

    result = float32(node.value)
    totweight = node.weight
    j = node.index + 1

    for j < N && totweight + all[j].weight <= K {
        totweight += all[j].weight
        result += float32(all[j].value)
        j++
    }

    k = j
    if k < N {
        result += float32((K - totweight) * all[k].value / all[k].weight)
    }

    return result
}

// TODO: BnB

func solveBranchAndBound(K int32, v []int32, w []int32) {
}

func solveDynamicProgramming(K int32, v []int32, w []int32) {
    var N int32 = int32(len(v))

    var O = make([][]int32, K+1)
    //fmt.Println(len(O))
    var x = make([]int32, N)
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

func solveFile(filename string) {
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

    //fmt.Println(n, K)
    //fmt.Println(v, w)
    solveDynamicProgramming(K, v, w)
    //solveBranchAndBound(K, v, w)
}

func main() {
    //fmt.Println("Solving file", os.Args[1])
    solveFile(os.Args[1])
}
