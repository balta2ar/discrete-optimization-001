#!/bin/sh

#
# $1 -- input problem file
# $2 -- output solution file
#

echo -n "read $1
optimize
write solution $2
quit
" | ~/mnt/big_ext4/bin/cs/scip/scip-3.0.1.linux.x86_64.gnu.opt.spx

