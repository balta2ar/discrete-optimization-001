#!/usr/bin/python
# -*- coding: utf-8 -*-

import sys
import math

import itertools
from subprocess import Popen, PIPE


INDENT = '    '
ORDER_PIP = 'order.pip'
ORDER_SOL = 'order.sol'


def runSolver(infile, outfile):
    process = Popen(['./runSolver.sh', infile, outfile], stdout=PIPE)
    (stdout, stderr) = process.communicate()


def perr(*args):
    sys.stderr.write(' '.join(map(str, args)) + '\n')


def length(point1, point2):
    return math.sqrt((point1[0] - point2[0])**2 + (point1[1] - point2[1])**2)


def allSubtours(items, a, b):
    for r in range(a, b):
        for tour in itertools.combinations(items, r):
            yield list(tour)
            # yield list(reversed(tour))


class OrderCustomersModel(object):
    '''
    Finally, I re-use my TSP MIP to optimize the order each of those
    customers is visited for each centroid.

    This model is a TSP MIP formulation.

    Constants:
        C -- # of clients / customers
        t_i_j -- transportation cost from customer i to j

    Variables:
        a_i_j -- customer i is followed by j (there is i -> j edge assigned)
        u_i -- helper variables for MTZ (see below)

    Objective f:
        // minimize 1. Total travel cost
        min sum(t_i_j * a_i_j)

    Subject to:
        // There is only one incoming and one outgoing edge
        sum(a_i_j)|i in C == 1 (forall j)
        sum(a_i_j)|j in C == 1 (forall i)

        // Miller-Tucker-Zemlin (MTZ) formulation for subtour elimination
        u_0 == 0
        1 <= u_i <= C - 1, (forall i != 0)
        u_i - u_j + 1 <= C * (1 - a_i_j) (forall i != 0, j != 0)

        # // Subtour formulation
        # sum(a_i_j)|i in S, j in S <= |S| - 1 (S is a subset of C, S != C, |S| > 1)
    '''
    def __init__(self, customers):
        self.customers = customers
        self.C = len(customers)

    def generatePip(self, outname):
        # 10 cheapest warehouses for problem #6
        # fixed = set([264, 109, 484, 462, 145, 482, 414, 401, 115, 95])
        #fixed = set([95])
        # lastFixed = sorted(list(fixed))[-1]

        def dist(i, j):
            ix, iy = self.customers[i][0], self.customers[i][1]
            jx, jy = self.customers[j][0], self.customers[j][1]
            dx, dy = ix - jx, iy - jy
            s = math.sqrt(dx * dx + dy * dy)
            # perr('i j', i, j, 'ix iy', ix, iy, 'jx jy', jx, jy,
            #      'dx dy', dx, dy, 's', s)
            return s

        def ass(i, j):
            return 'a_{0}_{1}'.format(i, j)

        def u(i):
            return 'u_{0}'.format(i)

        # def wopen(w):
        #     if w in fixed:
        #         return 'o_{0}'.format(w)
        #     else:
        #         return ''

        with open(outname, 'w') as f:
            for line in self.__class__.__doc__.splitlines():
                f.write('\ {0}\n'.format(line))

            f.write('\nMinimize\n\n')
            f.write('    obj:\n')
            f.write('\\ 1. Total travel cost\n')
            for i in range(self.C):
                for j in range(self.C):
                    plus = '' if (i == self.C-1) and (j == self.C-1) else ' +'
                    f.write('{0}{1} {2}{3}\n'.format(INDENT, dist(i, j), ass(i, j), plus))
                f.write('\n')

            f.write('Subject to\n')
            f.write('\\ 1. There is only one incoming and one outgoing edge\n')
            f.write('\\ i.e. sum(a_i_j)|i in C == 1 (forall j)\n')
            f.write('\\ i.e. sum(a_i_j)|j in C == 1 (forall i)\n')
            for i in range(self.C):
                for j in range(self.C):
                    plus = ' +'
                    f.write('{0}{1}{2}\n'.format(INDENT, ass(i, j), plus))
                f.write('{0}0 == 1\n'.format(INDENT))
                f.write('\n')

                for j in range(self.C):
                    plus = ' +'
                    f.write('{0}{1}{2}\n'.format(INDENT, ass(j, i), plus))
                f.write('{0}0 == 1\n'.format(INDENT))
                f.write('\n')

            # f.write('\\ 2. Subtour formulation\n')
            # f.write('\\ i.e. sum(a_i_j)|i in S, j in S <= |S| - 1 (S is a subset of C, S != C, |S| > 1)\n')
            # allIndices = list(range(self.C))
            # for tour in allSubtours(allIndices, 1, self.C / 2 + 1):
            #     for u in tour:
            #         f.write('{0}'.format(INDENT))
            #         for v in tour:
            #             f.write('{0} + '.format(ass(u, v)))
            #         f.write('\n')
            #     f.write('{0}0 <= {1}\n\n'.format(INDENT, len(tour) - 1))
            # f.write('\n')

                    #plus = '' if i == len(tour) - 2
            # print('N subtours', len(subtours))

            f.write('\\ 2. Miller-Tucker-Zemlin (MTZ) formulation for subtour elimination\n')
            f.write('\\ u_0 == 0\n')
            f.write('\\ 1 <= u_i <= C - 1, (forall i != 0) (see Bounds)\n')
            f.write('\\ u_i - u_j + 1 <= C * (1 - a_i_j) (forall i != 0, j != 0)\n')
            f.write('{0}{1}\n\n'.format(INDENT, 'u_0 = 0'))
            for i in range(1, self.C):
                for j in range(1, self.C):
                    f.write('{0}{1} - {2} + {3} {4} <= {5}\n'
                            ''.format(INDENT, u(i), u(j), self.C - 1, ass(i, j), self.C - 2))
                f.write('\n')

            f.write('Bounds\n')
            f.write('\\ 1 <= u_i <= C - 1, (forall i != 0)\n')
            for i in range(1, self.C):
                f.write('{0}1 <= {1} <= {2}\n'
                        ''.format(INDENT, u(i), self.C - 1))
            f.write('\n')

            f.write('Binary\n')
            f.write('\\ 1. Customer i is followed by j (there is i -> j edge assigned)\n')
            for i in range(self.C):
                for j in range(self.C):
                    f.write('{0}{1}\n'.format(INDENT, ass(i, j)))
                f.write('\n')

            f.write('End\n')

    def initAssignments(self):
        return [0] * self.C

    def parseSolution(self, inname):
        value = 0.0
        clientVehicle = self.initAssignments()

        with open(inname) as f:
            first = f.readline() # skip 'solution found' message
            if 'infeasible' in first:
                perr('Infeasible solution: SCIP refused to solve it')
                return None
            value = float(f.readline()[len('objective value:'):].strip())
            while True:
                line = f.readline()
                if len(line) == 0:
                    break
                if not line.startswith('a_'):  # assign variable
                    continue
                name, _, _ = line.strip().split()
                _, i, j = name.split('_')
                i, j = int(i), int(j)
                clientVehicle[i] = j

        current = 0
        order = [0]
        while True:
            nxt = clientVehicle[current]
            # perr('nxt', nxt)
            if nxt == 0:
                break
            order.append(nxt)
            current = nxt

        self.objectiveValue = value
        self.clientAssignments = order

        return order


def solveIt(inputData):
    # Modify this code to run your optimization algorithm

    # parse the input
    lines = inputData.split('\n')

    nodeCount = int(lines[0])

    points = []
    for i in range(1, nodeCount+1):
        line = lines[i]
        parts = line.split()
        points.append((float(parts[0]), float(parts[1])))

    orderModel = OrderCustomersModel(points)
    orderModel.generatePip(ORDER_PIP)
    # runSolver(ORDER_PIP, ORDER_SOL)
    # order = orderModel.parseSolution(ORDER_SOL)
    # perr('order', order)

    return ''

    # build a trivial solution
    # visit the nodes in the order they appear in the file
    solution = range(0, nodeCount)

    # calculate the length of the tour
    obj = length(points[solution[-1]], points[solution[0]])
    for index in range(0, nodeCount-1):
        obj += length(points[solution[index]], points[solution[index+1]])

    # prepare the solution in the specified output format
    outputData = str(obj) + ' ' + str(0) + '\n'
    outputData += ' '.join(map(str, solution))

    return outputData


if __name__ == '__main__':
    if len(sys.argv) > 1:
        fileLocation = sys.argv[1].strip()
        inputDataFile = open(fileLocation, 'r')
        inputData = ''.join(inputDataFile.readlines())
        inputDataFile.close()
        print solveIt(inputData)
    else:
        print 'This test requires an input file.  Please select one from the data directory. (i.e. python solver.py ./data/tsp_51_1)'
