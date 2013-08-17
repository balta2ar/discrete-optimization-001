#!/bin/bash

#
# $1 -- input data file
# $2 -- solution file
#

#/usr/bin/time -v go run solver.go $1 $2 > sol
cp $1 problem
# time go run solver.go $1 $2 > sol
time python2 ./solverPy.py $1 $2 > sol
# time pypy ./solverPy.py $1 $2 > sol
cat sol
Rscript ./plotSolution.R $1 sol 1>&2
# Rscript ./plotDataKmeans.R $1 1>&2
#python2 ./plot.py $1 < sol
#python2 ./plot2.py $1 sol > plot.dot #| dot -Tpng > snap.png
