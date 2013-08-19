#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
import sys
import math
from subprocess import Popen, PIPE
import itertools
import random
import datetime
import pickle

import pylab
from numpy import array
from scipy.cluster.vq import kmeans, vq
import matplotlib.pyplot as plt


def length(customer1, customer2):
    return math.sqrt((customer1[1] - customer2[1])**2 + (customer1[2] - customer2[2])**2)


def perr(*args):
    return
    sys.stderr.write(' '.join(map(str, args)) + '\n')


def perr2(*args):
    sys.stderr.write(' '.join(map(str, args)) + '\n')


def perr2nonl(*args):
    sys.stderr.write(' '.join(map(str, args)))


INDENT = '    '
KMEANS_PIC = 'kmeans.png'
ASSIGN_PIP = 'assign.pip'
ASSIGN_SOL = 'assign.sol'
ASSIGN_PIC = 'assign.png'
ORDER_PIP = 'order.pip'
ORDER_SOL = 'order.sol'
OPTIMIZE_PIC = 'optimize.png'
KMEANS_SET_BIN = 'kmeans.set.bin'
CLUSTER_SET_BIN = 'cluster.set.bin'


class AssignCustomersModel(object):
    '''
    Following this, I run a relatively simple MIP model to assign customers to
    the generated centroids without violating any capacity constraints. The
    objective is the sum of the distances from the assigned centroid.

    This model is very similiar to Warehouse Location Problem. Just a little
    simplier.

    Constants:
        V -- # of vehicles
        C -- # of clients / customers
        VC -- vehicle capacity (same for all vehicles)
        d_c -- customer demand
        t_c_v -- transportation cost from customer c to vehicle v

    Variables:
        a_c_v -- customer c assigned to vehicle v

    Objective f:
        // minimize 1. transportation cost
        min sum(t_c_v * a_c_v)

    Subject to:
        // customer assigned to only one vehicle
        sum(a_c_v)|v in V == 1 (forall c in C)

        // total assigned demand <= vehicle capacity
        sum(d_c * a_c_v)|c = 1..C <= VC
    '''
    def __init__(self, vehicles, capacity, customers):
        self.vehicles = vehicles
        self.customers = customers
        self.capacity = capacity
        self.C = len(customers)
        self.V = len(vehicles)

    def generatePip(self, outname):
        # 10 cheapest warehouses for problem #6
        # fixed = set([264, 109, 484, 462, 145, 482, 414, 401, 115, 95])
        #fixed = set([95])
        # lastFixed = sorted(list(fixed))[-1]

        def dist(c, v):
            vx, vy = self.vehicles[v][0], self.vehicles[v][1]
            cx, cy = self.customers[c][1], self.customers[c][2]
            dx, dy = cx - vx, cy - vy
            return math.sqrt(dx * dx + dy * dy)

        def ass(c, v):
            return 'a_{0}_{1}'.format(c, v)

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
            f.write('\\ 1. Total travel cost from client to assigned vehicle\n')
            for v in range(self.V):
                for c in range(self.C):
                    plus = '' if (v == self.V-1) and (c == self.C-1) else ' +'
                    f.write('{0}{1} {2}{3}\n'.format(INDENT, dist(c, v), ass(c, v), plus))
                f.write('\n')

            f.write('Subject to\n')
            f.write('\\ 1. Each client is assigned to only one vehicle\n')
            f.write('\\ i.e. sum(Ac)|v=1..V == 1\n')
            for c in range(self.C):
                for v in range(self.V):
                    plus = ' +'
                    f.write('{0}{1}{2}\n'.format(INDENT, ass(c, v), plus))
                f.write('{0}0 == 1\n'.format(INDENT))
                f.write('\n')
            # f.write('\n')

            f.write('\\ 2. Total clients\' demand for vehicle <= vehicle capacity VC\n')
            for v in range(self.V):
                for c in range(self.C):
                    less = ' <= {0}\n'.format(self.capacity) if c == self.C-1 else ' +'
                    f.write('{0}{1} {2}{3}\n'.format(INDENT, self.customers[c][0], ass(c, v), less))
            #f.write('\n')

            f.write('Bounds\n')
            #f.write('Binary\n')
            f.write('\n')

            f.write('Binary\n')
            f.write('\\ 1. assignment a_c_v -- client c assigned to vehicle v\n')
            #f.write('Bounds\n')
            for v in range(self.V):
                for c in range(self.C):
                    f.write('{0}{1}\n'.format(INDENT, ass(c, v)))
                    #f.write('{0}0 <= {1} <= 1\n'.format(INDENT, assignment(c, w)))
                f.write('\n')

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
                # perr('Infeasible solution: SCIP refused to solve it')
                return None
            value = float(f.readline()[len('objective value:'):].strip())
            while True:
                line = f.readline()
                if len(line) == 0:
                    break
                if not line.startswith('a_'): # assignment variable prefix
                    continue
                name, val, _ = line.strip().split()
                _, c, v = name.split('_')
                c, v = int(c), int(v)
                clientVehicle[c] = v

        self.objectiveValue = value
        self.clientAssignments = clientVehicle

        return clientVehicle

    # def cost(self, c, w):
    #     # setup cost + transportation cost
    #     return self.warehouses[w][1] + self.customerCosts[c][w]

    # def calcObjectiveValue(self):
    #     objValue = 0.0
    #     warehouseUsed = [False] * self.N
    #     for ci in range(self.M):
    #         wi = self.clientAssignments[ci]
    #         # client transportation cost
    #         objValue += self.customerCosts[ci][wi]
    #         if not warehouseUsed[wi]:
    #             warehouseUsed[wi] = True
    #             # warehouse setup cost
    #             objValue += self.warehouses[wi][1]
    #     return objValue

    # def formatSolution(self):
    #     first = '{0} {1}'.format(self.objectiveValue, 0)
    #     second = ' '.join(map(str, self.clientAssignments))
    #     print(first)
    #     print(second)
    #     return '{0}\n{1}'.format(first, second)


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

        # // Miller-Tucker-Zemlin (MTZ) formulation for subtour elimination
        # u_0 == 0
        # 1 <= u_i <= C - 1, (forall i != 0)
        # u_i - u_j + 1 <= C * (1 - a_i_j) (forall i != 0, j != 0)

        // Subtour formulation
        sum(a_i_j)|i in S, j in S <= |S| - 1 (S is a subset of C, S != C, |S| > 1)
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
            ix, iy = self.customers[i][1], self.customers[i][2]
            jx, jy = self.customers[j][1], self.customers[j][2]
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

            f.write('\\ 2. Subtour formulation\n')
            f.write('\\ i.e. sum(a_i_j)|i in S, j in S <= |S| - 1 (S is a subset of C, S != C, |S| > 1)\n')
            allIndices = list(range(self.C))
            for tour in allSubtours(allIndices, 1, self.C / 2 + 1):
                for u in tour:
                    f.write('{0}'.format(INDENT))
                    for v in tour:
                        f.write('{0} + '.format(ass(u, v)))
                    f.write('\n')
                f.write('{0}0 <= {1}\n\n'.format(INDENT, len(tour) - 1))
            f.write('\n')

                    #plus = '' if i == len(tour) - 2
            # print('N subtours', len(subtours))

            # f.write('\\ 2. Miller-Tucker-Zemlin (MTZ) formulation for subtour elimination\n')
            # f.write('\\ u_0 == 0\n')
            # f.write('\\ 1 <= u_i <= C - 1, (forall i != 0) (see Bounds)\n')
            # f.write('\\ u_i - u_j + 1 <= C * (1 - a_i_j) (forall i != 0, j != 0)\n')
            # f.write('{0}{1}\n\n'.format(INDENT, 'u_0 = 0'))
            # for i in range(1, self.C):
            #     for j in range(1, self.C):
            #         f.write('{0}{1} - {2} + {3} {4} <= {5}\n'
            #                 ''.format(INDENT, u(i), u(j), self.C - 1, ass(i, j), self.C - 2))
            #     f.write('\n')

            f.write('Bounds\n')
            # f.write('\\ 1 <= u_i <= C - 1, (forall i != 0)\n')
            # for i in range(1, self.C):
            #     f.write('{0}1 <= {1} <= {2}\n'
            #             ''.format(INDENT, u(i), self.C - 1))
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


def runSolver(infile, outfile):
    process = Popen(['./runSolver.sh', infile, outfile], stdout=PIPE)
    (stdout, stderr) = process.communicate()


def indices(x, xs):
    return [i for i in range(len(xs)) if xs[i] == x]


def plotAssignment(idx, coords, centroids, warehouseCoord, output):
    V = max(idx) + 1
    cmap = plt.cm.get_cmap('Dark2')
    customerColors = [cmap(1.0 * idx[i] / V) for i in range(len(idx))]
    centroidColors = [cmap(1.0 * i / V) for i in range(V)]
    xy = coords

    pylab.scatter([warehouseCoord[0]], [warehouseCoord[1]])
    pylab.scatter(centroids[:,0], centroids[:,1], marker='o', s=500, linewidths=2, c='none')
    pylab.scatter(centroids[:,0], centroids[:,1], marker='x', s=500, linewidths=2, c=centroidColors)
    pylab.scatter(xy[:,0], xy[:,1], s=100, c=customerColors)

    for cluster in range(V):
        customerIndices = indices(cluster, idx)
        clusterCustomersCoords = [warehouseCoord] + [list(coords[i]) for i in customerIndices]

        N = len(clusterCustomersCoords)
        for i in range(N):
            j = (i+1) % N
            x = clusterCustomersCoords[i][0]
            y = clusterCustomersCoords[i][1]
            dx = clusterCustomersCoords[j][0] - clusterCustomersCoords[i][0]
            dy = clusterCustomersCoords[j][1] - clusterCustomersCoords[i][1]
            pylab.arrow(x, y, dx, dy, color=centroidColors[cluster], fc="k",
                        head_width=1.0, head_length=2.5, length_includes_head=True)

    pylab.savefig(output)
    pylab.close()


def makedist(customers):
    def dist(i, j):
        ix, iy = customers[i][1], customers[i][2]
        jx, jy = customers[j][1], customers[j][2]
        dx, dy = ix - jx, iy - jy
        return math.sqrt(dx * dx + dy * dy)
    return dist


def objectiveValue(customers, paths):
    dist = makedist(customers)
    result = 0.0
    for path in paths:
        # perr('obj value path {0}'.format(path))
        for i in range(len(path) - 1):
            result += dist(path[i], path[i + 1])
    return result


def validatePath(customers, path, vehicleCapacity):
    dist = makedist(customers)
    demand, cost = 0, 0.0
    for j, v in enumerate(path):
        demand += customers[v][0]
    for i in range(len(path) - 1):
        cost += dist(path[i], path[i + 1])
    perr('path {0}: demand {1} / {2}, cost {3}'
         ''.format(path, demand, vehicleCapacity, cost))
    return demand <= vehicleCapacity


def validate(customers, paths, vehicleCapacity):
    dist = makedist(customers)
    perr('obj', objectiveValue(customers, paths))
    numErrors = 0
    for i, path in enumerate(paths):
        perr('validate path {0}'.format(path))
        demand, cost = 0, 0.0
        for j, v in enumerate(path):
            demand += customers[v][0]
        for j in range(len(path) - 1):
            # perr('cost {0}->{1} = {2}'
            #      ''.format(path[j], path[j + 1], dist(path[j], path[j + 1])))
            cost += dist(path[j], path[j + 1])
        perr('path {0}: demand {1} / {2}, cost {3}'
             ''.format(i, demand, vehicleCapacity, cost))
        if demand > vehicleCapacity:
            numErrors += 1

    if numErrors == 0:
        return True
    perr('Invalid solution, {0} violations'.format(numErrors))
    return False


def solutionFromIndexes(idx, customers, V):
    solution = []
    for cluster in range(V):
        path = [0] + [i + 1 for i in indices(cluster, idx)] + [0]
        solution.append(path)
    return solution


def formatSolution(customers, solution):
    obj = objectiveValue(customers, solution)
    lines = [' '.join(map(str, order)) for order in solution]
    return '''{0} 0
{1}'''.format(obj, '\n'.join(lines))


def makeSolver(inputData, kmeansSet, clusterSet):
    # parse the input
    lines = inputData.split('\n')

    parts = lines[0].split()
    customerCount = int(parts[0])
    V = int(parts[1])
    vehicleCapacity = int(parts[2])
    # depotIndex = 0

    customers = []
    for i in range(1, customerCount + 1):
        line = lines[i]
        parts = line.split()
        customers.append((int(parts[0]), float(parts[1]), float(parts[2])))

    def solver():
        return solveThreeStepsMIP(customers, V, vehicleCapacity, kmeansSet, clusterSet)

    return solver


def allSubtours(items, a, b):
    for r in range(a, b):
        for tour in itertools.combinations(items, r):
            yield list(tour)
            # yield list(reversed(tour))


def splitOnClusters(idx):
    V = max(idx) + 1
    result = []
    for c in range(V):
        result.append(tuple(sorted(indices(c, idx))))
    return tuple(sorted(result))


def solveThreeStepsMIP(customers, V, vehicleCapacity, kmeansSet, clusterSet):
    #
    # Find clusters using kmeans
    #

    # perr2(len(kmeansSet))

    i = 0
    while True:
        onlyCustomers = customers[1:]
        coords = array([[c[1], c[2]] for c in onlyCustomers])
        centroids, variance = kmeans(coords, V)
        idx, dist = vq(coords, centroids)
        warehouseCoord = [customers[0][1], customers[0][2]]

        tupleidx = splitOnClusters(tuple(idx))
        # perr2('>>>', splitOnClusters(tupleidx))
        # perr2(tupleidx)
        if tupleidx not in kmeansSet:
            # perr2('Added tupleidx set {0}'.format(tupleidx))
            kmeansSet.add(tupleidx)
            break

        #perr2('({0}) Idx set is present {1}'.format(i, tupleidx))
        perr2nonl('\rK-means generation round {0} (set size {1})'
                  ''.format(i, len(kmeansSet)))
        i += 1
    perr2nonl('\n')

    perr()
    perr('Running K-means')
    perr('centroids', centroids)
    perr('idx', list(idx))
    plotAssignment(idx, coords, centroids, warehouseCoord, KMEANS_PIC)
    perr('Validating solution from K-means')
    solution = solutionFromIndexes(idx, customers, V)
    validate(customers, solution, vehicleCapacity)
    perr(formatSolution(customers, solution))

    #
    # Reassign customers so that capacity is not violated
    #

    perr()
    perr('Running AssignCustomersModel')
    assignModel = AssignCustomersModel(centroids, vehicleCapacity, onlyCustomers)
    assignModel.generatePip(ASSIGN_PIP)
    runSolver(ASSIGN_PIP, ASSIGN_SOL)
    assign = assignModel.parseSolution(ASSIGN_SOL)
    if assign is None:
        # infeasible solution
        return None, None
    perr('assign', assign)
    plotAssignment(assign, coords, centroids, warehouseCoord, ASSIGN_PIC)
    perr('Validating solution from AssignCustomersModel')
    solution = solutionFromIndexes(assign, customers, V)
    validate(customers, solution, vehicleCapacity)
    perr(formatSolution(customers, solution))

    # tupleidx = splitOnClusters(tuple(idx))
    # tupleassign = splitOnClusters(tuple(assign))
    # perr2('>>>', tupleidx, '=>', tupleassign)

    #
    # Find best vehicle travel order
    #

    optimized = []
    for cluster in range(V):
        customerIndices = indices(cluster, assign)

        tupleIndices = tuple(sorted(customerIndices))
        if tupleIndices in clusterSet:
            indexOrder = list(clusterSet[tupleIndices])

        else:
            clusterCustomers = [customers[0]] + [customers[i+1] for i in customerIndices]
            perr()
            # perr2('Solving for cluster {0} / {1}'.format(cluster, V))
            perr2nonl('\rCluster {0} / {1}'.format(cluster + 1, V))
            perr('cluster, indices', cluster, customerIndices)
            perr('clusterCustomers', clusterCustomers)

            localCustomers = [customers[0]] + [customers[i+1] for i in customerIndices] + [customers[0]]
            perr('localCustomers', localCustomers)
            indexOrder = [0] + [i + 1 for i in customerIndices] + [0]
            validatePath(customers, indexOrder, vehicleCapacity)
            perr()

            orderModel = OrderCustomersModel(clusterCustomers)
            orderModel.generatePip(ORDER_PIP)
            runSolver(ORDER_PIP, ORDER_SOL)
            order = orderModel.parseSolution(ORDER_SOL)
            perr('order', order)

            # return ''

            localOrder = [i - 1 for i in order[1:]]
            perr('localOrder', localOrder)
            indexOrder = [0] + [customerIndices[i] + 1 for i in localOrder] + [0]
            perr('indexOrder', indexOrder)

            localCustomers = [customers[i] for i in indexOrder]
            perr('localCustomers', localCustomers)
            validatePath(customers, indexOrder, vehicleCapacity)

            clusterSet.add(tuple(indexOrder))

        optimized.append(indexOrder)
        open('current.sol', 'w').write(formatSolution(customers, optimized))

        # return ''

        # for i in range(len(idx)):
        #     if cluster == :
        # perr('optimized', optimized)

        # plotAssignment(assign, coords, centroids, warehouseCoord, OPTIMIZE_PIC)

        # break
    perr2nonl('\n')

    perr('optimized', optimized)
    validate(customers, optimized, vehicleCapacity)
    return objectiveValue(customers, optimized), formatSolution(customers, optimized)

    return ''

    # vehicleTours = []

    # customerIndexs = set(range(1, customerCount))  # start at 1 to remove depot index

    # for v in range(0, vehicleCount):
    #     # print "Start Vehicle: ",v
    #     vehicleTours.append([])
    #     capacityRemaining = vehicleCapacity
    #     while sum([capacityRemaining >= customers[ci][0] for ci in customerIndexs]) > 0:
    #         used = set()
    #         order = sorted(customerIndexs, key=lambda ci: -customers[ci][0])
    #         for ci in order:
    #             if capacityRemaining >= customers[ci][0]:
    #                 capacityRemaining -= customers[ci][0]
    #                 vehicleTours[v].append(ci)
    #                 # print '   add', ci, capacityRemaining
    #                 used.add(ci)
    #         customerIndexs -= used

    # # checks that the number of customers served is correct
    # assert sum([len(v) for v in vehicleTours]) == customerCount - 1

    # # calculate the cost of the solution; for each vehicle the length of the route
    # obj = 0
    # for v in range(0, vehicleCount):
    #     vehicleTour = vehicleTours[v]
    #     if len(vehicleTour) > 0:
    #         obj += length(customers[depotIndex],customers[vehicleTour[0]])
    #         for i in range(0, len(vehicleTour) - 1):
    #             obj += length(customers[vehicleTour[i]],customers[vehicleTour[i + 1]])
    #         obj += length(customers[vehicleTour[-1]],customers[depotIndex])

    # # prepare the solution in the specified output format
    # outputData = str(obj) + ' ' + str(0) + '\n'
    # for v in range(0, vehicleCount):
    #     outputData += str(depotIndex) + ' ' + ' '.join(map(str,vehicleTours[v])) + ' ' + str(depotIndex) + '\n'

    # return outputData


def saveSet(s, filename):
    pickle.dump(s, open(filename, 'w'))


def loadSet(filename):
    if os.path.exists(filename):
        return pickle.load(open(filename))
    return set()


def solveBest(inputData, K=100000):
    start = datetime.datetime.now()
    last = start
    bestSolution = None
    bestCost = -1
    kmeansSet = loadSet(KMEANS_SET_BIN)
    clusterSet = loadSet(CLUSTER_SET_BIN)
    solver = makeSolver(inputData, kmeansSet, clusterSet)

    for i in range(K):
        runningLast = str(datetime.datetime.now() - last)
        runningTotal = str(datetime.datetime.now() - start)
        last = datetime.datetime.now()
        perr2('Iteration {0} / {1} (running {2} / {3})'
              ''.format(i + 1, K, runningLast, runningTotal))
        cost, solution = solver()
        if cost is None:
            continue
        if (bestSolution is None) or (cost < bestCost):
            bestCost = cost
            bestSolution = solution
            open('current.sol', 'w').write(solution)
            open('best.sol', 'w').write(solution)
            perr2('New best solution: {0}'.format(cost))

        saveSet(kmeansSet, KMEANS_SET_BIN)
        saveSet(clusterSet, CLUSTER_SET_BIN)
        # for x in kmeansSet:
        #     perr2(x)
    return bestSolution


def solveIt(inputData):
    return solveBest(inputData)


if __name__ == '__main__':
    random.seed()
    if len(sys.argv) > 1:
        fileLocation = sys.argv[1].strip()
        inputDataFile = open(fileLocation, 'r')
        inputData = ''.join(inputDataFile.readlines())
        inputDataFile.close()
        perr('Solving:', fileLocation)
        print solveIt(inputData)
    else:
        print 'This test requires an input file.  Please select one from the data directory. (i.e. python solver.py ./data/vrp_5_4_1)'
