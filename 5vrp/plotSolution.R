#require(ggplot2)
library(RColorBrewer)

args <- commandArgs(trailingOnly=T)
problem <- args[1]
solution <- args[2]
output <- 'solution.png'
if (length(args) > 2) {
    output <- args[3]
}

# problem <- 'data/vrp_5_4_1'
# solution <- 'sol'
# output <- 'solution.png'
png(output, width=800, height=800)

#filename = 'data/vrp_30_3_1'
#filename = 'data/vrp_200_17_2'
#filename = 'problem_6.dat'
par(mar=c(2, 2, 1, 1), lwd=3)

conn <- file(problem, 'r')
header <- readLines(conn, 1)
info <- as.integer(read.table(textConnection(header))[1,])
V <- info[2]
colors <- colorRampPalette(brewer.pal(8, 'Dark2'))(V)
# print(V)
# print(colors)

d <- na.omit(read.table(conn, sep=' ', header=F, skip=0, blank.lines.skip=T))
x <- d$V2
y <- d$V3
plot(x, y)
legendLabels = (1:V)-1

first <-function(x) x[1:(length(x)-1)]
last <- function(x) x[2:(length(x))]
lshift <- function(x) c(x[2:length(x)], x[1])

i <- 1
conn <- file(solution, 'r')
l <- readLines(conn, 1)
l <- readLines(conn, 1)
while (length(l) == 1) {
    o <- as.integer(read.table(textConnection(l))[1,]) + 1
    legendLabels[i] <- paste(legendLabels[i], '/', length(o)-2, sep='')
    if (length(o) > 2) {
        a <- first(x[o])
        b <- first(y[o])
        d <- last(x[o])
        e <- last(y[o])
        arrows(a, b, d, e,
               length=0.25, angle=20, col=colors[i])
    }
    l <- readLines(conn, 1)
    i <- i + 1
}
legend(min(x)-3, max(y)+3, legendLabels, col=colors,
       lty=c(1,1), lwd=c(5, 5), cex=2, bty='n')
text(x, y, labels=0:(length(x)-1), pos=2, cex=2)
close(conn)
dev.off()

# tt   <- read.table( filename, sep=" ", header=FALSE )
# tt   <- tt[-1,]
# depo <- tt[1,]
# 
# pl <- ggplot( tt, aes(x=V2, y=V3, size=V1) ) + 
#   geom_point( color='navyblue' ) +
#   geom_point( size=9, color='darkgreen', data=depo ) +
#   theme_bw()
# 
# ggsave(plot=pl, file=output, width=6, height=6 )
