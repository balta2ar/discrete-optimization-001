require(ggplot2)

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
info <- as.integer(read.table(textConnection(header))[1,])
V <- info[2]

tt   <- read.table( conn, sep=" ", header=FALSE )
# tt   <- tt[-1,]
names( tt ) <- c( 'demand','x','y' )

depo <- tt[1,]

zz         <- tt[ , c(2,3) ]
km         <- kmeans( zz, V, nstart=100 )
tt$cluster <- factor( km$cluster )

pl <- ggplot( tt, aes( x=x, y=y, size=demand ) ) + 
      geom_point( aes( color=cluster ) ) +
      geom_point( size=9, color='darkgreen', data=depo ) +
      theme_bw()

ggsave( plot=pl, file=output, width=6, height=6 )

close(conn)
