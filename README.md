## Project Maya

This repo contains cli commands to generate and verify zero-knowledge proof of image crop transformation.
This enables a verifier to verify if an image was actually cropped from an original image.

## Steps to run
1. Install [golang](https://go.dev/doc/install)
2. Run `go install .`
3. Run `mkdir proofs`

## To prove a crop transformation

```shell
project-maya prove crop --cropped-image=./sample/cropped.png  --original-image=./sample/original.png --proof-dir=proofs/
```

You would get output like this:
```shell
Image has width 10 and height 10
Image has width 5 and height 5
13:49:15 INF compiling circuit
13:49:15 INF parsed circuit inputs nbPublic=75 nbSecret=300
13:49:15 INF building constraint builder nbConstraints=75
Circuit compilation time: 0.000574583
13:49:15 DBG constraint system solver done nbConstraints=75 took=0.081958
13:49:15 DBG prover done backend=groth16 curve=bn254 nbConstraints=75 took=2.720417
Time taken to prove:  0.002965834
```

## To verify a crop transformation
```shell
project-maya verify crop --cropped-image=./sample/cropped.png --proof-dir=proofs/
```

You would get output like this:
```shell
Image has width 5 and height 5
13:49:37 DBG verifier done backend=groth16 curve=bn254 took=1.765916
Proof verified 🎉
```