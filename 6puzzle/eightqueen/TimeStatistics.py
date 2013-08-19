#
# Eight Queen Problem Solver
# Time statistics module
# (c) 2008 baltazar
#

import datetime

class TimeStatistics(object):
    def __init__(self):
        self.reset()

    def reset(self):
        self.stat = []

    def mark(self):
        self.stat.append(datetime.datetime.today())

    def printAvegare(self):
        if len(self.stat) == 0: return
        average = self.stat[0]
        for i in range(1, len(self.stat)):
            average += self.stat[i]
        average /= len(self.stat)
        print 'Average duration:', average

    def printLongest(self):
        if len(self.stat) == 0: return
        longest = self.stat[0]
        for i in range(1, len(self.stat)):
            if self.stat[i] > longest:
                longest = self.stat[i]
        print 'Longest duration:', longest

    def printShortest(self):
        if len(self.stat) == 0: return
        shortest = self.stat[0]
        for i in range(1, len(self.stat)):
            if self.stat[i] < shortest:
                shortest = self.stat[i]
        print 'Shortest duration:', shortest

    def printTotal(self):
        if len(self.stat) < 2: return
        print 'Total duration:', self.stat[-1] - self.stat[0]
