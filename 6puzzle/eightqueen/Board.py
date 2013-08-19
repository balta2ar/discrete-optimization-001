#
# Eight Queens Problem Solver
# Board module
# Contains board class
# (c) 2008 baltazar
#

import math

class Board(object):
    def __init__(self, size):
        # print 'Board size:', size

        self.infinity = size * size
        self.charQueen = '@'
        self.charNoQueen = '.'

        self.size = size
        self.cols = [size for i in range(0, size)]
        self.rows = [size for i in range(0, size)]
        self.queens = 0

    def toString(self):
        return '''{0}
{1}'''.format(self.size, ' '.join(map(str, self.cols)))

    def show(self):
        for i in range(0, self.size):
            for j in range(0, self.size):
                if self.cols[j] == i: print self.charQueen,
                else:                 print self.charNoQueen,
            print

    def placeQueen(self, col, row):
        self.cols[col] = row
        self.rows[row] = col
        self.queens += 1

    def removeQueen(self, col, row):
        self.cols[col] = self.size
        self.rows[row] = self.size
        self.queens -= 1

    def isSolved(self):
        if self.queens == self.size: return True
        else:                        return False

    def isColumnEmpty(self, col):
        if self.cols[col] == self.size: return True
        else:                           return False

    def isRowEmpty(self, row):
        if self.rows[row] == self.size: return True
        else:                           return False

    def isPosWithinBoard(self, col, row):
        if col < 0 or col >= self.size or row < 0 or row >= self.size: return False
        return True

    def cellContainsQueen(self, col, row):
        if not self.isPosWithinBoard(col, row): return False
        if self.cols[col] == row: return True
        return False

    def queenAllowed(self, col, row):
        if not self.isPosWithinBoard(col, row): return False
        if not self.isColumnEmpty(col):         return False
        if not self.isRowEmpty(row):            return False

        for i in range(1, self.size):
            if self.cellContainsQueen(col - i, row - i): return False
            if self.cellContainsQueen(col - i, row + i): return False
            if self.cellContainsQueen(col + i, row - i): return False
            if self.cellContainsQueen(col + i, row + i): return False

        return True

    def rowsFree(self, col):
        free = 0
        for i in range(0, self.size):
            if self.queenAllowed(col, i):
                free += 1

        return free

    def getNextLCV(self, col, tried):
        maximum = self.infinity
        output = self.infinity

        for i in range(0, self.size):
            if not tried[i] and self.queenAllowed(col, i):
                affected = self.countAllAffectedCells(col, i)
                if affected < maximum:
                    maximum = affected
                    output = i

        return output

    def isCellAffected(self, col, row):
        if self.cellContainsQueen(col, row): return 0
        if self.queenAllowed(col, row): return 1
        return 0

    def countAllAffectedCells(self, col, row):
        if self.cellContainsQueen(col, row): return 0

        output = 0
        for i in range(1, self.size):
            output += self.isCellAffected(col, row - i)
            output += self.isCellAffected(col, row + i)
            output += self.isCellAffected(col - i, row)
            output += self.isCellAffected(col + i, row)
            output += self.isCellAffected(col - i, row - i)
            output += self.isCellAffected(col - i, row + i)
            output += self.isCellAffected(col + i, row - i)
            output += self.isCellAffected(col + i, row + i)

        return output

    def countRangeAffectedCells(self, col, row):
        if self.cellContainsQueen(col, row): return 0

        output = 0
        for i in range(1, int(math.log(self.size))):
            output += self.isCellAffected(col, row - i)
            output += self.isCellAffected(col, row + i)
            output += self.isCellAffected(col - i, row)
            output += self.isCellAffected(col + i, row)
            output += self.isCellAffected(col - i, row - i)
            output += self.isCellAffected(col - i, row + i)
            output += self.isCellAffected(col + i, row - i)
            output += self.isCellAffected(col + i, row + i)

        return output
