#!/usr/bin/python
# -*- coding: utf-8 -*-

import math

def length(customer1, customer2):
    return math.sqrt((customer1[1] - customer2[1])**2 + (customer1[2] - customer2[2])**2)

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
    for i in range(1, customerCount+1):
        line = lines[i]
        parts = line.split()
        customers.append((int(parts[0]), float(parts[1]),float(parts[2])))

    # build a trivial solution
    # assign customers to vehicles starting by the largest customer demands

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

