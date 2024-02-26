# Comparison

We are currently comparing the performance of generating zk proof of transformations by implementing it
in risc-zero, SP1 and our hand-rolled custom gnark circuits.

## Maya circuits

## Crop transformation

| Original size | Final size | Circuit compilation | Crop time    | Proving        |
|---------------|------------|---------------------|--------------|----------------|
| 1000x1000     | 10x10      | 0.944125959s        | 45.484458ms  | 16.449535208s  |
| 1000x1000     | 100x100    | 0.909873625s        | 40.301291ms  | 20.899984083s  |
| 1000x1000     | 250x250    | 1.10235225s         | 54.143208ms  | 32.518697042s  |
| 1000x1000     | 500x500    | 1.4310838750000001s | 95.847667ms  | 61.816334333s  |
| 1000x1000     | 750x750    | 2.125492958s        | 154.770583ms | 108.671270792s |

## Rotate90 transformation

| Original size | Circuit compilation | Rotate90 time (ms) | Proving (s)   | Proof size (bytes) |
|---------------|---------------------|--------------------|---------------|--------------------|
| 10x10         | 0.000375375s        | 275.417µs          | 0.032882208   | 164                |
| 100x100       | 0.027196292         | 2.292125           | 1.833014333   | 164                |
| 250x250       | 0.171775334         | 17.180625          | 11.3016845    | 164                |
| 500x500       | 0.656973583         | 60.774875          | 51.5077695    | 164                |
| 1000x1000     | 1.593718708         | 157.525208         | 107.182107334 | 164                |
