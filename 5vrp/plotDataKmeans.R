require(ggplot2)

args <- commandArgs(trailingOnly=T)
filename <- args[1]
output <- 'view.png'
if (length(args) > 1) {
    output <- args[2]
}

#filename = 'problem_4.dat'
tt   <- read.table( filename, sep=" ", header=FALSE )
tt   <- tt[-1,]
names( tt ) <- c( 'demand','x','y' )

depo <- tt[1,]

zz         <- tt[ , c(2,3) ]
km         <- kmeans( zz, 10, nstart=10 )
tt$cluster <- factor( km$cluster )

pl <- ggplot( tt, aes( x=x, y=y, size=demand ) ) + 
      geom_point( aes( color=cluster ) ) +
      geom_point( size=9, color='darkgreen', data=depo ) +
      theme_bw()

ggsave( plot=pl, file=output, width=6, height=6 )
