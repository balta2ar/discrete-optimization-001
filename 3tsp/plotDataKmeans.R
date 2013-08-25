require(ggplot2)
library(RColorBrewer)

args <- commandArgs(trailingOnly=T)
filename <- args[1]
output <- 'view.png'
solution <- ''
if (length(args) == 3) {
    solution <- args[2]
    output <- args[3]
}
if (length(args) == 2) {
    s <- args[2]
    x <- substr(s, nchar(s) - 3, nchar(s))
    if (x == '.png') {
        output <- s
    } else {
        solution <- s
    }
}
if (output == 'id') {
    output = paste(filename, '.png', sep='')
}

readData <- function(filename, solution) {
    conn <- file(filename, 'r')
    header <- readLines(conn, 1)
    V <- as.integer(read.table(textConnection(header))[1,])

    tt <- read.table(conn, sep=" ", header=FALSE)
    close(conn)

    cost <- 0.0
    if (solution != '') {
        conn <- file(solution, 'r')
        costHeader <- readLines(conn, 1)
        cost <- as.numeric(read.table(textConnection(costHeader)))[1]
        ord <- as.integer(read.table(conn, sep=" ", header=FALSE)) + 1
        close(conn)
        tt <- tt[ord, ]
    }

    tt <- rbind(tt, tt[1, ])
    names(tt) <- c('x','y')
    list(cost, tt)
}

clusterize <- function(d, N, nowork) {
    if (nowork) {
        NA
    } else {
        km         <- kmeans(d, N, iter.max=1000, nstart=100)
        cluster    <- factor(km$cluster)
        cluster
    }
}

plotData <- function(d, cluster, cost) {
    point_size <- 0.5
    # colors     <- colorRampPalette(brewer.pal(8, 'Dark2'))(length(tt$cluster))
    pl <- ggplot(d, aes(x=x, y=y)) +
        geom_text(aes(x=min(x), y=min(y)), label=cost, hjust=0, vjust=0)
    if (length(cluster) == 1) {
        print('no cluster')
        pl = pl +
            geom_point(size=point_size) +
            geom_path(size=0.1)
    } else {
        print('cluster')
        pl = pl +
            geom_point(size=point_size, color=cluster) +
            geom_path(size=0.1, color=cluster)
    }
    pl = pl + coord_fixed() + theme_bw()
    big <- data.frame(cbind(
        x=c(14416.0, 14592.0, 14880.0, 14592.0),
        y=c(9328.0, 9636.0, 9108.0, 9328.0)))
    endsA <- data.frame(cbind(
        x=c(12784.0, 11984.0),
        y=c(10098.0, 11660.0)))
    endsB <- data.frame(cbind(
        x=c(18704.0, 12784.0),
        y=c(8580.0, 4840.0)))
    he = head(d, 4)
    ta = head(tail(d, 4), 4)
    print(he)
    print(ta)
    pl = pl + geom_point(data=big[1:2,], size=1.0, col='red') +
              geom_point(data=big[3:4,], size=1.0, col='blue') +
              # geom_point(data=endsA, size=2.0, col='green') +
              # geom_point(data=endsB, size=2.0, col='magenta') +
              geom_point(data=he, size=1.0, col='darkred') +
              geom_point(data=ta, size=1.0, col='green')
    pl
}

N <- 4
# point_size <- 0.5

lst <- readData(filename, solution)
cost <- lst[[1]]
tt <- lst[[2]]
cluster <- clusterize(tt, N, solution == '')
pl <- plotData(tt, cluster, cost)

# ggsave(plot=pl, file=output, width=10, height=10)
# ggsave(plot=pl, file=output, width=side * ratio, height=side)
ggsave(plot=pl, file=output)
