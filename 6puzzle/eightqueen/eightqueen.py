#
# Eight Queens Problem Solver
# Finds first satisfiable position for eight queens on a
# chessboard so that no one queen may attact any other queen
# (c) 2008 baltazar
#

import sys

import TimeStatistics
import Solver

SIZE = 8
if len(sys.argv) > 1: SIZE = int(sys.argv[1])

def solve(size):
    solver = Solver.Solver(size)
    if solver.solve():
        return solver.solution()
        # solver.showStatistics()
    else:
        return 'No solution found'

def main(size):
    timeStat = TimeStatistics.TimeStatistics()
    timeStat.mark()
    solver = Solver.Solver(size)
    if solver.solve():
        solver.showBoard()
        # solver.showStatistics()
    else:
        print 'No solution found'
    timeStat.mark()

    #print 'Statistics'
    # print 'Size', size
    # timeStat.printTotal()

#for i in range(8, 31):
# for i in range(8, 30):
#     main(i)
# main(SIZE)
