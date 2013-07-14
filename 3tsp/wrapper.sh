#!/bin/bash

#
# $1 -- input data file
# [$2] -- algorithm
# [$3] -- number of colors (for csp)
#

#
# 1 arg  == find best cached solution; if exists, print; otherwise
#           calculate, cache solution
# 2 args == calc solution and determine ncolors automatically, cache solution
# 3 args == calc solution with specified ncolors, cache solution
#




/usr/bin/time -v go run solver.go $1 $2 > sol
cat sol
python2 ./plot.py $1 < sol
#python2 ./plot2.py $1 sol > plot.dot #| dot -Tpng > snap.png





exit 0









SOLUTIONDIR="solution"
INPUT=$(basename $1)

if [ ! -d "$SOLUTIONDIR" ]; then
    mkdir $SOLUTIONDIR
fi

if [ $# -eq 1 ]; then
    BEST=$(ls $SOLUTIONDIR/$INPUT* 2>/dev/null | sort | head -1)
    if [ -n "$BEST" ]; then
        #echo "Showing $BEST"
        cat $BEST
        exit 0
    fi
fi

TEMP="output.tmp"
go run solver.go $1 $2 $3 > $TEMP

CODE=$?
if [ $CODE -eq 0 ]; then
    # show the solution and save it
    cat $TEMP
    NCOLORS=$(cat $TEMP | head -1 | awk '{print $1}')
    SOLUTION="$SOLUTIONDIR/$INPUT.$NCOLORS"
    cat $TEMP > $SOLUTION
else
    # failed to find the solution
    echo "Could not find the solution"
    exit 1
fi
