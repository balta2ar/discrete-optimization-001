#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
import sys
from subprocess import Popen, PIPE

from numpy import array
from scipy.cluster.vq import kmeans, vq


def indices(x, xs):
    return [i for i in range(len(xs)) if xs[i] == x]


def solveTSP(customers):
    return ''


def solveIt(inputData):
    lines = inputData.split('\n')

    N = int(lines[0])
    customers = []
    for i in range(1, N+1):
        pair = map(float, lines[i].split())
        customers.append((pair[0], pair[1]))

    V = 25
    coords = array([c for c in customers])
    centroids, _ = kmeans(coords, V)
    idx, _ = vq(coords, centroids)
    clusterIndices = [indices(i, idx) for i in range(V)]

    # print(idx)
    # print(clusterIndices)
    for i, v in enumerate(clusterIndices):
        print('{0}: {1}'.format(i, len(v)))

    return solveTSP(customers)

    return open('5.sol').read()
    # Writes the inputData to a temporay file

    tmpFileName = 'tmp.data'
    tmpFile = open(tmpFileName, 'w')
    tmpFile.write(inputData)
    tmpFile.close()

    # Runs the command: java Solver -file=tmp.data

    #process = Popen(['go', 'run', 'solver.go', tmpFileName], stdout=PIPE)
    process = Popen(['./wrapper.sh', tmpFileName], stdout=PIPE)
    #process = Popen(['./wrapper.sh', 'data/gc_500_1'], stdout=PIPE)
    (stdout, stderr) = process.communicate()

    # removes the temporay file

    os.remove(tmpFileName)

    return stdout.strip()


if __name__ == '__main__':
    if len(sys.argv) > 1:
        fileLocation = sys.argv[1].strip()
        inputDataFile = open(fileLocation, 'r')
        inputData = ''.join(inputDataFile.readlines())
        inputDataFile.close()
        print solveIt(inputData)
    else:
        print 'This test requires an input file.  Please select one from the data directory. (i.e. python solver.py ./data/ks_4_0)'
