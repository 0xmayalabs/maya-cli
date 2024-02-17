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
08:28:21 INF compiling circuit
08:28:21 INF parsed circuit inputs nbPublic=375 nbSecret=0
08:28:21 INF building constraint builder nbConstraints=75
Circuit compilation time: 0.001732209
08:28:21 DBG constraint system solver done nbConstraints=75 took=0.410666
08:28:21 DBG prover done backend=groth16 curve=bn254 nbConstraints=75 took=2.501583
Time taken to prove:  0.003019167
```

## To verify a crop transformation
```shell
project-maya verify crop --cropped-image=./sample/cropped.png  --original-image=./sample/original.png --proof-dir=proofs/
```

You would get output like this:
```shell
Image has width 10 and height 10
Image has width 5 and height 5
08:29:52 DBG verifier done backend=groth16 curve=bn254 took=2.562125
Proof verified 🎉
```