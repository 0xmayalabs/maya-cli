# Run maya

You can run maya commands easily using docker.

Maya CLI provides commands for both `proving` an image transformation and `verifying` it.

You can also look at the available transformations for `proving` by running:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya:latest prove --help
```

Similarly, to look at the available transformations for `verifying`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya:latest verify --help
```
