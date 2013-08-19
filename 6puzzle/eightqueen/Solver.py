#
# Eight Queens Problem Solver
# Solver module
# Solves Eight Queen Problem
# (c) 2008 baltazar
#

import time

import Board
import Global

class Solver(object):
    def __init__(self, size):
        self.size = size
        self.board = Board.Board(size)
        self.steps = 0
        self.cycles = time.clock()
        self.infinity = size * size

    def solution(self):
        return self.board.toString()

    def showBoard(self):
        # self.board.show()
        print(self.board.toString())

    def showStatistics(self):
        print '> Resursion steps:', self.steps
        print '> CPU cycles:', self.cycles

    def getEmptyColumn(self):
        if Global.heuristicForVar == 'Brute':
            for i in range(0, self.size):
                if self.board.isColumnEmpty(i):
                    return i
            return infinity
        elif Global.heuristicForVar == 'MRV':
            minimum = self.infinity
            output = self.infinity
            for i in range(0, self.size):
                rowsFree = self.board.rowsFree(i)
                if rowsFree < minimum and rowsFree > 0:
                    minimum = rowsFree
                    output = i

            return output
        elif Global.heuristicForVar == 'MCV':
            maximum = 0
            output = self.infinity
            for i in range(0, self.size):
                constrainingFactor = self.calculateCF(i)
                if self.board.rowsFree(i) > 0 and constrainingFactor > maximum:
                    maximum = constrainingFactor
                    output = i

            return output

        return self.infinity

    def calculateCF(self, col):
        rowsFree = self.board.rowsFree(col)
        if rowsFree == 0: return 0

        output = 0
        for i in range(0, self.size):
            output += self.board.countRangeAffectedCells(col, i)

        return output * 100 / rowsFree

    def solve(self):
        self.steps += 1

        if self.board.isSolved():
            self.cycles = time.clock() - self.cycles
            return True

        nextCol = self.getEmptyColumn()
        if nextCol == self.infinity: return False

        if Global.heuristicForVal == 'Brute':
            for i in range(0, self.size):
                if self.board.queenAllowed(nextCol, i):
                    self.board.placeQueen(nextCol, i)
                    if self.solve():
                        return True
                    else:
                        self.board.removeQueen(nextCol, i)
        elif Global.heuristicForVal == 'LCV':
            tried = [False for i in range(0, self.size)]
            t = self.board.getNextLCV(nextCol, tried)
            while t < self.infinity:
                self.board.placeQueen(nextCol, t)
                if self.solve():
                    return True
                else:
                    self.board.removeQueen(nextCol, t)
                tried[t] = True
                t = self.board.getNextLCV(nextCol, tried)

        return False

    def start(self):
        self.cycles = time.clock()
        self.solve()
