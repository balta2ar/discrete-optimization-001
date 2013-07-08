import sys
import igraph


def main():
    if len(sys.argv) < 2:
        print "usage %s graph_file [xsize] [ysize] < solution_file" % sys.argv[0]
        sys.exit(1)

    fn = sys.argv[1]
    sx = sys.argv[2] if len(sys.argv) >= 3 else 400
    sy = sys.argv[3] if len(sys.argv) >= 4 else sx

    fin = open(fn)
    nodeCount = int(fin.readline())
    layout = [map(float, line.split()) for line in fin if line]
    fin.close()

    sys.stdin.readline()  # skip header
    order = map(int, sys.stdin.readline().split())
    edges = zip(order, order[1:] + order[:1])

    if sorted(order) != range(nodeCount):
        print "something wrong with solution!"

    g = igraph.Graph(edges=edges, directed=True)
    style = {}
    style["margin"] = 25
    style["layout"] = layout
    style["vertex_size"] = 25
    style["bbox"] = (sx, sy)
    style["vertex_label"] = map(str, range(nodeCount))
    style["vertex_label_dist"] = 0
    style["vertex_color"] = "white"
    igraph.plot(g, fn + ".png", **style)


if __name__ == "__main__":
    main()
