require(ggplot2)

args <- commandArgs(trailingOnly=T)
filename <- args[1]
output <- 'view.png'
if (length(args) > 1) {
    output <- args[2]
}

#filename = 'data/vrp_30_3_1'
#filename = 'data/vrp_200_17_2'
#filename = 'problem_6.dat'

tt   <- read.table( filename, sep=" ", header=FALSE )
tt   <- tt[-1,]
depo <- tt[1,]

pl <- ggplot( tt, aes(x=V2, y=V3, size=V1) ) + 
  geom_point( color='navyblue' ) +
  geom_point( size=9, color='darkgreen', data=depo ) +
  theme_bw()

ggsave(plot=pl, file=output, width=6, height=6 )
