#!/usr/bin/python
# -*- coding: utf-8 -*-

def solveIt(n):
    # Modify this code to run your puzzle solving algorithm
    
    # to be consistent with other example code, we will
    # represent the decision variables as an array (not a matrix)
    # there is no need to do this in your solver
    
    # define the domains of all the variables (1..n*n)
    domains = [range(1,n*n+1)]*(n*n)
    
    # start a trivial depth first search for a solution
    sol = tryall([],domains)
    
    # prepare the solution in the specified output format
    # if no solution is found, put 0s
    outputData = str(n) + '\n'
    if sol == None:
        print 'no solution found.'
        for i in range(0,n):
            outputData += ' '.join(map(str, [0]*n))+'\n'
    else: 
        for i in range(0,n):
            outputData += ' '.join(map(str, sol[i*n:(i+1)*n]))+'\n'
    
    return outputData


# this is a depth first search of all assignments
def tryall(assignment, domains):
    # base-case: if the domains list is empty, all values are assigned
    # check if it is a solution, return None if it is not
    if len(domains) == 0:
        if checkIt(assignment):
            return assignment
        else:
            return None
    
    # recursive-case: try each value in the next domain
    # if we find a solution return it. otherwise, try the next value
    else:
        for v in domains[0]:
            if not v in assignment:
                sol = tryall(assignment[:]+[v],domains[1:])
                if sol != None:
                    return sol


# checks if an assignment is feasible
# because sol is an array (not a matrix), checks are more cryptic
import math
def checkIt(sol):
    n = int(math.sqrt(len(sol)))
    m = n*(n*n+1)/2
    
    #for i in range(0,n):
    #    print sol[i*n:(i+1)*n]
    
    items = set(sol)
    if len(items) != len(sol):
        #print len(items),len(sol) 
        return False
    
    for i in range(0,n):
        #print 'row',i,sol[i*n:(i+1)*n]
        if sum(sol[i*n:(i+1)*n]) != m:
            return False
        #print 'column',i,sol[i:len(sol):n]
        if sum(sol[i:len(sol):n]) != m:
            return False
        if i < n-1:
            if sol[i*n+i] > sol[(i+1)*n+(i+1)]:
                return False 
    
    #print 'diag 1',i,[sol[i*n+i] for i in range(0,n)]
    if sum([sol[i*n+i] for i in range(0,n)]) != m:
        return False
    #print 'diag 2',i,[sol[i*n+(n-i-1)] for i in range(0,n)]
    if sum([sol[i*n+(n-i-1)] for i in range(0,n)]) != m:
        return False
    
    return True


import sys

if __name__ == "__main__":
    if len(sys.argv) > 1:
        try:
            n = int(sys.argv[1].strip())
        except:
            print sys.argv[1].strip(), 'is not an integer'
        print 'Solving Size:', n
        print(solveIt(n))
    
    else:
        print('This test requires an instance size.  Please select the size of problem to solve. (i.e. python magicSquareSolver.py 3)')

