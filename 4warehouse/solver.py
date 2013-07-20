#!/usr/bin/python
# -*- coding: utf-8 -*-


#import os
import sys
from subprocess import Popen, PIPE
from time import gmtime, strftime


PIP_NAME = 'problem.pip'
SOL_NAME = 'problem.sol'
INDENT = ' ' * 9


def pe(*args):
    now = strftime("%Y-%m-%d %H:%M:%S", gmtime())
    msg = now + ' ' + ' '.join(map(str, args)) + '\n'
    sys.stderr.write(msg)
    sys.stderr.flush()


def generatePip(warehouses, customerSizes, customerCosts):
    def assignment(c, w):
        return 'a_{0}_{1}'.format(c, w)

    N = len(warehouses)
    M = len(customerSizes)
    with open(PIP_NAME, 'w') as f:
        f.write('Minimize\n\n')
        f.write('\\ 1. Total setup cost\n')
        f.write('    obj:\n')
        for w in range(N):
            for c in range(M):
                f.write('{0}{1} {2} +\n'.format(INDENT, warehouses[w][1], assignment(c, w)))
        f.write('\n')
        f.write('\\ 2. Total travel cost from client to assigned warehouse\n')
        for w in range(N):
            for c in range(M):
                plus = '' if (w == N-1) and (c == M-1) else ' +'
                f.write('{0}{1} {2}{3}\n'.format(INDENT, customerCosts[c][w], assignment(c, w), plus))
        f.write('\n')

        f.write('Subject to\n')
        f.write('\\ Total clients\' demand for warehouse <= warehouse capacity\n')
        for w in range(N):
            for c in range(M):
                less = ' <= {0}\n'.format(warehouses[w][0]) if c == M-1 else ' +'
                f.write('{0}{1} {2}{3}\n'.format(INDENT, customerSizes[c], assignment(c, w), less))
        #f.write('\n')

        f.write('\\ Each client is assigned to only one warehouse\n')
        f.write('\\ i.e. sum(Ac)|w=1..N == 1\n')
        for c in range(M):
            for w in range(N):
                plus = ' == 1\n' if w == N-1 else ' +'
                f.write('{0}{1}{2}\n'.format(INDENT, assignment(c, w), plus))
        #f.write('\n')

        f.write('Bounds\n')
        f.write('\n')

        f.write('Binary\n')
        for w in range(N):
            for c in range(M):
                f.write('{0}{1}\n'.format(INDENT, assignment(c, w)))
            f.write('\n')
        #f.write('\n')

        f.write('End\n')


def runSolver():
    process = Popen(['./runSolver.sh', PIP_NAME, SOL_NAME], stdout=PIPE)
    (stdout, stderr) = process.communicate()


def initAssignments(N, M):
    return [0] * M


def parseSolution(N, M):
    value = 0
    clientWarehouse = initAssignments(N, M)

    with open(SOL_NAME) as f:
        f.readline() # skip 'solution found' message
        value = f.readline()[len('objective value:'):].strip()
        while True:
            line = f.readline()
            if len(line) == 0:
                break
            if not line.startswith('a_'): # assignment variable prefix
                continue
            name, val, _ = line.strip().split()
            _, c, w = name.split('_')
            clientWarehouse[int(c)] = int(w)

    return value, clientWarehouse


def formatSolution(obj, assignments):
    first = '{0} 0'.format(obj, 0)
    second = ' '.join(map(str, assignments))
    print(first)
    print(second)
    return '{0}\n{1}'.format(first, second)


def solveWLP(warehouses, customerSizes, customerCosts):
    N = len(warehouses)
    M = len(customerSizes)
    pe('Generating problem description')
    generatePip(warehouses, customerSizes, customerCosts)
    pe('Running solver')
    runSolver()
    obj, solution = parseSolution(N, M)
    pe('Solution is ready')
    return formatSolution(obj, solution)


def solveIt(inputData):
    # Modify this code to run your optimization algorithm

    # parse the input
    lines = inputData.split('\n')

    parts = lines[0].split()
    warehouseCount = int(parts[0])
    customerCount = int(parts[1])

    warehouses = []
    for i in range(1, warehouseCount+1):
        line = lines[i]
        parts = line.split()
        warehouses.append((int(parts[0]), float(parts[1])))

    customerSizes = []
    customerCosts = []

    lineIndex = warehouseCount+1
    for i in range(0, customerCount):
        customerSize = int(lines[lineIndex+2*i])
        customerCost = map(float, lines[lineIndex+2*i+1].split())
        customerSizes.append(customerSize)
        customerCosts.append(customerCost)

    return solveWLP(warehouses, customerSizes, customerCosts)

    # build a trivial solution
    # pack the warehouses one by one until all the customers are served

    solution = [-1] * customerCount
    capacityRemaining = [w[0] for w in warehouses]

    warehouseIndex = 0
    for c in range(0, customerCount):
        if capacityRemaining[warehouseIndex] >= customerSizes[c]:
            solution[c] = warehouseIndex
            capacityRemaining[warehouseIndex] -= customerSizes[c]
        else:
            warehouseIndex += 1
            assert capacityRemaining[warehouseIndex] >= customerSizes[c]
            solution[c] = warehouseIndex
            capacityRemaining[warehouseIndex] -= customerSizes[c]

    used = [0]*warehouseCount
    for wa in solution:
        used[wa] = 1

    # calculate the cost of the solution
    obj = sum([warehouses[x][1]*used[x] for x in range(0,warehouseCount)])
    for c in range(0, customerCount):
        obj += customerCosts[c][solution[c]]

    # prepare the solution in the specified output format
    outputData = str(obj) + ' ' + str(0) + '\n'
    outputData += ' '.join(map(str, solution))

    print(outputData)
    #return outputData


import sys

if __name__ == '__main__':
    if len(sys.argv) > 1:
        fileLocation = sys.argv[1].strip()
        inputDataFile = open(fileLocation, 'r')
        inputData = ''.join(inputDataFile.readlines())
        inputDataFile.close()
        solveIt(inputData)
        #print 'Solving:', fileLocation
        #print solveIt(inputData)
    else:
        print 'This test requires an input file.  Please select one from the data directory. (i.e. python solver.py ./data/wl_16_1)'

