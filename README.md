# Maya CLI

The Maya CLI is a 100% open source, contributor-friendly command line tool written in [go](https://go.dev/doc/install) to generate and verify zero-knowledge proof of image transformations.

**[Install](https://docs.mayalabs.tech/install.md)**
| [Docs](https://docs.mayalabs.tech)
| [Telegram](https://t.me/+hM1lNjgLFRdjMGE1)

Okay

## Benchmark

We have run benchmark tests on [Macbook Pro M1](https://www.apple.com/in/shop/buy-mac/macbook-pro/16-inch-macbook-pro), 16GB RAM as well as 
on [AWS EC2 r6i.8xlarge](https://aws.amazon.com/ec2/instance-types/r6i/) with 32 cores of CPU and 256 GB RAM.

To run benchmark tests on a new machine (say `m7g.8xlarge`), run:
```shell
go test -v ./cmd --results-dir=../book/perf/m7g.8xlarge -run ^TestBenchmark
```
