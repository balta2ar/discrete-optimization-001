#!/usr/bin/python
# -*- coding: utf-8 -*-

import os
from subprocess import Popen, PIPE


def solveIt(inputData):

    # Writes the inputData to a temporay file

    tmpFileName = 'tmp.data'
    tmpFile = open(tmpFileName, 'w')
    tmpFile.write(inputData)
    tmpFile.close()

    # Runs the command: java Solver -file=tmp.data

    process = Popen(['java', 'Solver', '-file=' + tmpFileName],
                    stdout=PIPE)
    (stdout, stderr) = process.communicate()

    # removes the temporay file

    os.remove(tmpFileName)

    return stdout.strip()


import sys

if __name__ == '__main__':
    if len(sys.argv) > 1:
        fileLocation = sys.argv[1].strip()
        inputDataFile = open(fileLocation, 'r')
        inputData = ''.join(inputDataFile.readlines())
        inputDataFile.close()
        print solveIt(inputData)
    else:
        print 'This test requires an input file.  Please select one from the data directory. (i.e. python solver.py ./data/ks_4_0)'

