#!/bin/bash

#
# $1 -- input data file
# [$2] -- algorithm
#

# /usr/bin/time -v go run solver.go $1 $2 > sol
# python2 ./solver.py $1 > sol
pypy ./solver.py $1 > sol

# python2 ./solverPy.py $1 > sol
cat sol
# python2 ./plot.py $1 < sol
#python2 ./plot2.py $1 sol > plot.dot #| dot -Tpng > snap.png
