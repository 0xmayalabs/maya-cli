# Performance

We evaluate the efficiency and performance of generating zero-knowledge (ZK) proofs for image transformations through the
development of custom ZK circuits.

These circuits are developed using [gnark](https://github.com/Consensys/gnark), an open-source framework to design ZK circuits. 
Our objective is to benchmark the performance of our custom circuits against general purpose zkVMs like
as [RiscZero](https://github.com/risc0/risc0) and [SP1](https://github.com/succinctlabs/sp1).

It's important to note that during this performance evaluation, the original image remains private, 
and only the transformed image is revealed.