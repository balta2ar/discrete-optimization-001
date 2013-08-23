require(ggplot2)
library(RColorBrewer)

args <- commandArgs(trailingOnly=T)
filename <- args[1]
output <- 'view.png'
if (length(args) > 1) {
    output <- args[2]
    if (output == 'id') {
        output = paste(filename, '.png', sep='')
    }
}

#filename = 'problem_4.dat'
conn <- file(filename, 'r')
header <- readLines(conn, 1)
V <- as.integer(read.table(textConnection(header))[1,])
# V <- info[2]

tt <- read.table( conn, sep=" ", header=FALSE )
close(conn)
# tt   <- tt[-1,]
names(tt) <- c('x','y')

# depo <- tt[1,]

N <- 25
km         <- kmeans(tt, N, iter.max=1000, nstart=100)
tt$cluster <- factor(km$cluster)
# colors     <- colorRampPalette(brewer.pal(8, 'Dark2'))(length(tt$cluster))
pl         <- ggplot(tt, aes(x=x, y=y)) + geom_point(size=1.0, color=tt$cluster) + theme_bw()
ggsave(plot=pl, file=output, width=10, height=10)

if (F) {
    zz         <- tt[ , c(2,3) ]
    km         <- kmeans( zz, V, nstart=100 )
    tt$cluster <- factor( km$cluster )

    pl <- ggplot( tt, aes( x=x, y=y, size=demand ) ) + 
          geom_point( aes( color=cluster ) ) +
          geom_point( size=9, color='darkgreen', data=depo ) +
          theme_bw()
}


