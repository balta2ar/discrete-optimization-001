#!/usr/bin/env python

import sys


def read_graph(stream):
    n = int(stream.readline())
    for ln in stream:
        yield map(float, ln.split())


def read_route(stream):
    obj = stream.readline()
    return map(int, stream.readline().split())


def plot_graph(points, walk):
    print "graph G {\n  node [shape=point, fixedsize=true];"
    for (i, (x, y)) in enumerate(points, 1):
        print '  %d [label="", pos="%d,%d!"];' % (i, x, y)
    prev = walk[-1]
    for i in walk:
        print "  %d -- %d;" % (prev + 1, i + 1)
        prev = i
    print "}"


def main(fname, result):
    with file(fname) as stream:
        points = list(read_graph(stream))
    with file(result) as stream:
        route = read_route(stream)
    plot_graph(points, route)

if __name__ == "__main__":
    main(sys.argv[1], sys.argv[2])

#  #!/usr/bin/env python
#
# import sys
# #import alg
#
#
# def read_graph(stream):
#     #n = int(stream.readline())
#     for ln in stream:
#         yield map(int, ln.split())
#
#
# def read_solution():
#     h = sys.stdin.readline()  # skip header
#     if len(h) == 0:
#         return []
#     return map(int, sys.stdin.readline().split())
#
#
# def plot_graph(points, walk):
#     print "graph G {\n  node [shape=point, fixedsize=true];"
#     for (i, (x, y)) in enumerate(points, 1):
#         print '  %d [label="", pos="%d,%d!"];' % (i, x, y)
#     prev = walk[-1]
#     for i in walk:
#         print "  %d -- %d;" % (prev + 1, i + 1)
#         prev = i
#     print "}"
#
#
# def main(fname):
#     with file(fname) as stream:
#         points = list(read_graph(stream))
#         plot_graph(points, read_solution)
#
#
# if __name__ == "__main__":
#     main(sys.argv[1])
