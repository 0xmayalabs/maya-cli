# Run maya

You can run maya commands easily using docker.

Maya CLI provides commands for both `proving` an image transformation and `verifying` it.

You can also look at the available transformations for `proving` by running:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove --help
```

Similarly, to look at the available transformations for `verifying`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify --help
```

## Crop

To prove that an image is cropped correctly, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove crop \
--original-image=./sample/original.png \
--cropped-image=./sample/cropped.png \
--height-start-new=0 \
--width-start-new=0 \
--proof-dir=proofs
```

To verify that an image is cropped correctly, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify crop \
--cropped-image=./sample/cropped.png \
--proof-dir=proofs
```