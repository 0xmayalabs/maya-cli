# Performance

We are currently comparing the performance of generating zk proof of transformations by implementing it
in risc-zero, SP1 and our hand-rolled custom [gnark](https://github.com/Consensys/gnark) circuits.

## Maya circuits

Maya circuits are optimized ZK circuits written in [gnark](https://github.com/Consensys/gnark) to prove
image transformations. We are currently in the process of benchmarking proving code on [SP1](https://github.com/succinctlabs/sp1) 
and [risczero](https://github.com/risc0/risc0).

Note that the performance is assessed while keeping the original image private,
with only the final, altered image being revealed.