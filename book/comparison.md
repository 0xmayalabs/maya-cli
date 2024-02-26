# Comparison

We are currently comparing the performance of generating zk proof of transformations by implementing it
in risc-zero, SP1 and our hand-rolled custom gnark circuits.

## Maya circuits

## Crop transformation

| Original size | Final size | Circuit compilation | Crop time   | Proving       |
|---------------|------------|---------------|-------------|---------------|
| 1000x1000     | 10x10      | 0.944125959s | 45.484458ms | 16.449535208s |
| 1000x1000     | 100x100    | 0.909873625s | 40.301291ms | 20.899984083s |
| 1000x1000     | 250x250    | 1.10235225s | 54.143208ms | 32.518697042s |
| 1000x1000     | 500x500    |1.4310838750000001s | 95.847667ms | 61.816334333s |
| 1000x1000     | 750x750    | 2.125492958s | 154.770583ms | 108.671270792s |
