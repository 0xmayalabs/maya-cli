## Maya ZK Benchmarks

This repo contains cli commands to generate and verify zero-knowledge proof of image transformation.
This enables a verifier to verify if an image was actually transformed correctly from an original image.

See [docs](https://docs.mayalabs.tech) for more information.

To run benchmark tests:
```shell
go test -v ./cmd --results-dir=../book/perf/<machine> -run ^TestBenchmark
```
