#!/usr/bin/python
# -*- coding: utf-8 -*-

import math
from subprocess import Popen, PIPE

import pylab
from numpy import array
from scipy.cluster.vq import kmeans, vq
import matplotlib.cm as cm
import matplotlib.pyplot as plt


def length(customer1, customer2):
    return math.sqrt((customer1[1] - customer2[1])**2 + (customer1[2] - customer2[2])**2)


INDENT = '    '
KMEANS_PIC = 'kmeans.png'
ASSIGN_PIP = 'assign.pip'
ASSIGN_SOL = 'assign.sol'
ASSIGN_PIC = 'assign.png'
ORDER_PIP = 'order.pip'
ORDER_SOL = 'order.sol'


class AssignCustomersModel(object):
    '''
    Following this, I run a relatively simple MIP model to assign customers to
    the generated centroids without violating any capacity constraints. The
    objective is the sum of the distances from the assigned centroid.


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
        sum(d_c * a_c_v)|c = 1..C <= c_v (forall v in V)
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

            f.write('Minimize\n\n')
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
            f.readline() # skip 'solution found' message
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

        #pe('hello', clientWarehouse)
        #pe('hello', pprint.pformat(probabilities))
        #pe('hello', probabilities)

    def cost(self, c, w):
        # setup cost + transportation cost
        return self.warehouses[w][1] + self.customerCosts[c][w]

    def calcObjectiveValue(self):
        objValue = 0.0
        warehouseUsed = [False] * self.N
        for ci in range(self.M):
            wi = self.clientAssignments[ci]
            # client transportation cost
            objValue += self.customerCosts[ci][wi]
            if not warehouseUsed[wi]:
                warehouseUsed[wi] = True
                # warehouse setup cost
                objValue += self.warehouses[wi][1]
        return objValue

    def formatSolution(self):
        first = '{0} {1}'.format(self.objectiveValue, 0)
        second = ' '.join(map(str, self.clientAssignments))
        print(first)
        print(second)
        return '{0}\n{1}'.format(first, second)


class OrderCustomersModel(object):
    ''' Finally, I re-use my TSP MIP to optimize the order each of those
    customers is visited for each centroid.  '''
    pass


def runSolver(infile, outfile):
    process = Popen(['./runSolver.sh', infile, outfile], stdout=PIPE)
    (stdout, stderr) = process.communicate()


def solveIt(inputData):
    # Modify this code to run your optimization algorithm

    # parse the input
    lines = inputData.split('\n')

    parts = lines[0].split()
    customerCount = int(parts[0])
    vehicleCount = int(parts[1])
    vehicleCapacity = int(parts[2])
    depotIndex = 0

    customers = []
    for i in range(1, customerCount + 1):
        line = lines[i]
        parts = line.split()
        customers.append((int(parts[0]), float(parts[1]), float(parts[2])))

    # build a trivial solution
    # assign customers to vehicles starting by the largest customer demands

    coords = array([[c[1], c[2]] for c in customers])
    centroids, variance = kmeans(coords, vehicleCount)
    idx, dist = vq(coords, centroids)

    print(centroids)
    print(list(idx))
    print(centroids[0][0])

    V = vehicleCount
    cmap = plt.cm.get_cmap('Dark2')
    customerColors = [cmap(1.0 * idx[i] / V) for i in range(len(idx))]
    centroidColors = [cmap(1.0 * i / V) for i in range(V)]
    xy = coords

    pylab.scatter(centroids[:,0],centroids[:,1], marker='o', s = 500, linewidths=2, c='none')
    pylab.scatter(centroids[:,0],centroids[:,1], marker='x', s = 500, linewidths=2, c=centroidColors)
    pylab.scatter(xy[:,0],xy[:,1], s=100, c=customerColors)
    pylab.savefig(KMEANS_PIC)

    assignModel = AssignCustomersModel(centroids, vehicleCapacity, customers)
    assignModel.generatePip(ASSIGN_PIP)
    runSolver(ASSIGN_PIP, ASSIGN_SOL)
    assign = assignModel.parseSolution(ASSIGN_SOL)
    print(assign)


    return ''

    vehicleTours = []

    customerIndexs = set(range(1, customerCount))  # start at 1 to remove depot index

    for v in range(0, vehicleCount):
        # print "Start Vehicle: ",v
        vehicleTours.append([])
        capacityRemaining = vehicleCapacity
        while sum([capacityRemaining >= customers[ci][0] for ci in customerIndexs]) > 0:
            used = set()
            order = sorted(customerIndexs, key=lambda ci: -customers[ci][0])
            for ci in order:
                if capacityRemaining >= customers[ci][0]:
                    capacityRemaining -= customers[ci][0]
                    vehicleTours[v].append(ci)
                    # print '   add', ci, capacityRemaining
                    used.add(ci)
            customerIndexs -= used

    # checks that the number of customers served is correct
    assert sum([len(v) for v in vehicleTours]) == customerCount - 1

    # calculate the cost of the solution; for each vehicle the length of the route
    obj = 0
    for v in range(0, vehicleCount):
        vehicleTour = vehicleTours[v]
        if len(vehicleTour) > 0:
            obj += length(customers[depotIndex],customers[vehicleTour[0]])
            for i in range(0, len(vehicleTour) - 1):
                obj += length(customers[vehicleTour[i]],customers[vehicleTour[i + 1]])
            obj += length(customers[vehicleTour[-1]],customers[depotIndex])

    # prepare the solution in the specified output format
    outputData = str(obj) + ' ' + str(0) + '\n'
    for v in range(0, vehicleCount):
        outputData += str(depotIndex) + ' ' + ' '.join(map(str,vehicleTours[v])) + ' ' + str(depotIndex) + '\n'

    return outputData


import sys

if __name__ == '__main__':
    if len(sys.argv) > 1:
        fileLocation = sys.argv[1].strip()
        inputDataFile = open(fileLocation, 'r')
        inputData = ''.join(inputDataFile.readlines())
        inputDataFile.close()
        print 'Solving:', fileLocation
        print solveIt(inputData)
    else:
        print 'This test requires an input file.  Please select one from the data directory. (i.e. python solver.py ./data/vrp_5_4_1)'
