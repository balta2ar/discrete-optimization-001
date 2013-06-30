#!/usr/bin/python
# -*- coding: utf-8 -*-

def solveIt(n):
    # Modify this code to run your puzzle solving algorithm
    
    # define the domains of all the variables (0..n-1)
    domains = [range(0,n)]*n
    
    # start a trivial depth first search for a solution
    sol = tryall([],domains)
    
    # prepare the solution in the specified output format
    # if no solution is found, put 0s
    outputData = str(n) + '\n'
    if sol == None:
        print 'no solution found.'
        outputData += ' '.join(map(str, [0]*n))+'\n'
    else: 
        outputData += ' '.join(map(str, sol))+'\n'
        
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
            sol = tryall(assignment[:]+[v],domains[1:])
            if sol != None:
                return sol


# checks if an assignment is feasible
def checkIt(sol):
    n = len(sol)
    count = {}
    for i in range(0,n):
        count[i] = 0
    for i in range(0,n):
        count[sol[i]] += 1
    for i in range(0,n):
        if sol[i] != count[i]:
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
        print('This test requires an instance size.  Please select the size of problem to solve. (i.e. python magicSeriesSolver.py 5)')

