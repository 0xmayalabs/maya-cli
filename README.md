## Maya CLI

This repo contains the CLI commands to generate and verify zero-knowledge proof of image transformations.

See [docs](https://docs.mayalabs.tech) for more information.

To run benchmark tests:
```shell
go test -v ./cmd --results-dir=../book/perf/<machine> -run ^TestBenchmark
```
