discrete-optimization-001
=========================

## Coursera [Discrete Optimization](https://class.coursera.org/optimization-001/class/index) course programming assignments source code

Below is a list of optimization techniques/topics and languages used for each assignment.
Note that all assignments come with default trivial solution in python -- I don't count
python unless I implemented something non-trivial in this assignment. Also,
I started almost every assignment with implementing greedy or random solver,
I don't count them either.

#### Knapsack

- Go (DP, BnB solver)
- Dynamic Programming (DP)
- Branch and Bound (BnB)

#### Graph Coloring (GC)

- Bash (wrapper)
- Go (CP solver)
- Constraint Programming (CP)
  - Minimum Remaining Variable (MRV)
  - Least Constraining Value (LCV)
  - Arc Consistency (AC3)

#### Traveling Salesman Problem (TSP)

- Bash (wrapper)
- Go (LS solver)
- R (visualization)
- Python
  - visualization (igraph)
  - MIP problem generator, parser
  - scipy.kmeans
- Local Search (LS)
  - Simulated Annealing (SA)
  - Metropolis
  - 2-opt
  - [Late Acceptance Hill Climbing](http://www.cs.stir.ac.uk/research/publications/techreps/pdf/TR192.pdf)
- Mixed Integer Programming (MIP)
  - external solver: [SCIP](http://scip.zib.de)
  - problem format: PIP
  - problem formulations: Miller-Tucker-Zemlin, subtour elimination

#### Warehouse Location Problem (WLP)

- Bash (wrapper)
- Python (MIP problem generator, parser)
- Mixed Integer Programming (MIP)
  - external solver: [SCIP](http://scip.zib.de)
  - problem format: PIP
  - problem formulations: SimpleModel, LectureModel

#### Vehicle Routing Problem (VRP)

- Bash (wrapper)
- Go (LS solver, unit test)
- R (visualization)
- Python (MIP problem generator, parser)
- Local Search (LS)
  - Simulated Annealing (SA)
  - Metropolis
  - neighbour generation moves: 1. move customer from one route to another 2. swap two customers
- Mixed Integer Programming (MIP)
  - external solver: [SCIP](http://scip.zib.de)
  - problem format: PIP
  - problem formulations: AssignCustomersModel (similar to WLP), OrderCustomersModel (similar to TSP)

#### Puzzle Challenge (PC)

- Python (nqueens CP solver back from university times)
- Octave (magic square)
- Constraint Programming (CP)
  - Minimum Remaining Variable (MRV)
  - Least Constraining Value (LCV)
