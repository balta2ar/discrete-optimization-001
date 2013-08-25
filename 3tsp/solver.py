#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
import sys
from subprocess import Popen, PIPE
import pickle
import math

# from numpy import array
# from scipy.cluster.vq import kmeans, vq


SUBPROBLEM = 'cluster/problem/subproblem.{0}.txt'
SUBSOLUTION = 'cluster/solution/subproblem.{0}.sol'
CLUSTER_CONTEXT = 'cluster/cluster.context.bin'


def perr(*args):
    sys.stderr.write(' '.join(map(str, args)) + '\n')


def indices(x, xs):
    return [i for i in range(len(xs)) if xs[i] == x]


def saveContext(name, customers, centroids, idx, clusterIndices):
    data = customers, centroids, idx, clusterIndices
    pickle.dump(data, open(name, 'w'))


def loadContext(name):
    return pickle.load(open(name))


def solveTSPInParts(customers, N, V):
    # perr('Kmeans of {0} clusters'.format(V))
    # coords = array([c for c in customers])
    # centroids, _ = kmeans(coords, V)
    # idx, _ = vq(coords, centroids)
    # clusterIndices = [indices(i, idx) for i in range(V)]

    # # print(idx)
    # # print(clusterIndices)
    # for i, v in enumerate(clusterIndices):
    #     perr('{0}: {1}'.format(i, len(v)))

    # # split into subproblems
    # perr('Splitting into {0} clusters'.format(V))
    # for i, v in enumerate(clusterIndices):
    #     clusterCoords = [pair for j, pair in enumerate(customers) if j in v]
    #     points = '\n'.join(['{0} {1}'.format(x, y) for x, y in clusterCoords])
    #     problem = '{0}\n{1}'.format(len(v), points)
    #     filename = SUBPROBLEM.format(i)
    #     open(filename, 'w').write(problem)
    #     perr('Saved cluster {0} to file {1}'.format(i, filename))

    # saveContext(CLUSTER_CONTEXT, customers, centroids.tolist(), idx.tolist(), clusterIndices)
    # return ''

    def makedist(xs):
        def dist(i, j):
            ax, ay = xs[i]
            bx, by = xs[j]
            dx, dy = ax - bx, ay - by
            return math.sqrt(dx ** 2 + dy ** 2)
        return dist

    customersDist = makedist(customers)

    def cost(solution):
        c = 0.0
        for i in range(len(solution)):
            j = (i + 1) % len(solution)
            c += customersDist(solution[i], solution[j])
            # perr(c)
        # perr(c)
        return c

    customers1, centroids, idx, clusterIndices = loadContext(CLUSTER_CONTEXT)
    # saveContext(CLUSTER_CONTEXT + '.list', customers, centroids.tolist(), idx.tolist(), clusterIndices)
    # return ''

    # centroidsDist = makedist(centroids)

    # perr(customers[:5])
    # perr(customers1[:5])
    # return ''

    # return ''

    # print(list(idx))
    solution = []
    reorderedClusterIndices = []
    for i, v in enumerate(clusterIndices):
        filename = SUBSOLUTION.format(i)
        order = map(int, open(filename).read().splitlines()[1].split())
        reordered = [v[x] for x in order]
        reorderedClusterIndices.append(reordered)
        solution.extend(reordered)
        # print(v)
        # print(order)
        # print(reordered)
        # break
    # sol = '{0} 0\n{1}'.format(cost(solution), ' '.join(map(str, solution)))
    # return sol

    # for i, v in enumerate(centroids):
    #     d = centroidsDist(0, i)
    #     perr('dist({0}, {1}) = {2}'.format(0, i, d))

    def nearestCluster(clusterId, clusters):
        '''Find nearest cluster to clusterId'''
        minPair, minCost = None, None
        for otherCluster in range(len(clusters)):
            if otherCluster == clusterId:
                continue
            minPair, minCost = None, None
            for i in range(len(clusters[clusterId]) - 1):
                for j in range(len(clusters[otherCluster]) - 1):
                    fromA = clusters[clusterId][i]
                    fromB = clusters[clusterId][i + 1]
                    toA = clusters[otherCluster][j]
                    toB = clusters[otherCluster][j + 1]
                    delta = (- customersDist(fromA, fromB)
                             - customersDist(toA, toB)
                             + customersDist(fromA, toA)
                             + customersDist(fromB, toB))
                    if (minCost is None) or (delta < minCost):
                        minCost = delta
                        minPair = clusterId, otherCluster, i, j
                        # minPair = clusterId, otherCluster, i, j, delta, fromA, toA, customers[fromA], customers[fromB], customers[toA], customers[toB]
            # print('min pair for cluster {0}: {1}'.format(otherCluster, minPair))
        return minPair

    def mergeClusters(clusters, clusterFrom, clusterTo, i, j):
        '''Merge two selected clusters and return new list of clusters'''
        merged = []
        rFrom, rTo = clusters[clusterFrom], clusters[clusterTo]
        merged.extend(rFrom[:i + 1])
        merged.extend(reversed(rTo[:j + 1]))
        merged.extend(reversed(rTo[j + 1:]))
        merged.extend(rFrom[i + 1:])
        newClusters = [merged]
        for k, v in enumerate(clusters):
            if (k != clusterFrom) and (k != clusterTo):
                newClusters.append(v)
        return newClusters

    # mergingClusters = copy.deepcopy(reorderedClusterIndices)
    mergingClusters = reorderedClusterIndices
    for i in range(len(mergingClusters) - 1):
        perr('Merge {0}, len {1}'.format(i, len(mergingClusters)))
        nearest = nearestCluster(0, mergingClusters)
        mergingClusters = mergeClusters(mergingClusters, *nearest)
        # perr(type(mergingClusters))
        # perr(len(mergingClusters))

    # perr(mergingClusters)
    solution = mergingClusters[0]

    # perr(minPair)
    # clusterFrom, clusterTo, i, j, delta, fromA, toA, fromA, fromB, toA, toB = minPair

    # ca = clusterIndices[clusterFrom]
    # cb = clusterIndices[clusterTo]
    # perr(customers[ca[0]])
    # perr(customers[ca[len(ca) - 1]])
    # perr(customers[cb[0]])
    # perr(customers[cb[len(cb) - 1]])

    # solution.extend(rFrom[:i + 1])
    # solution.extend(rTo[j + 1:])
    # solution.extend(reversed(rTo[:j + 1]))
    # solution.extend(reversed(rFrom[i + 1:]))

    # for i, v in enumerate(reorderedClusterIndices):
    #     if i not in bannedClusters:
    #         solution.extend(v)

    sol = '{0} 0\n{1}'.format(cost(solution), ' '.join(map(str, solution)))
    return sol


def solveIt(inputData):
    # return open('sol').read()

    lines = inputData.split('\n')

    N = int(lines[0])
    customers = []
    for i in range(1, N+1):
        pair = map(float, lines[i].split())
        customers.append((pair[0], pair[1]))

    V = 4
    return solveTSPInParts(customers, N, V)

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
