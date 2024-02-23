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

## Flip vertical

To prove that an image is correctly flipped vertically, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove flip-vertical \
--original-image=./sample/original.png \
--final-image=./sample/flipped_vertical.png \
--proof-dir=proofs
```

To verify that an image is correctly flipped vertically, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify flip-vertical \
--final-image=./sample/flipped_vertical.png \
--proof-dir=proofs
```

## Flip Horizontal

To prove that an image is correctly flipped horizontally, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove flip-horizontal \
--original-image=./sample/original.png \
--final-image=./sample/flipped_horizontal.png \
--proof-dir=proofs
```

To verify that an image is correctly flipped horizontally, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify flip-horizontal \
--final-image=./sample/flipped_horizontal.png \
--proof-dir=proofs
```

## Rotate 90

To prove that an image is correctly rotated by 90 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove rotate90 \
--original-image=./sample/original.png \
--final-image=./sample/rotated90.png \
--proof-dir=proofs
```

To verify that an image is correctly rotated by 90 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify rotate90 \
--final-image=./sample/rotated90.png \
--proof-dir=proofs
```

## Rotate 180

To prove that an image is correctly rotated by 180 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove rotate180 \
--original-image=./sample/original.png \
--final-image=./sample/rotated180.png \
--proof-dir=proofs
```

To verify that an image is correctly rotated by 180 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify rotate180 \
--final-image=./sample/rotated180.png \
--proof-dir=proofs
```

## Rotate 270

To prove that an image is correctly rotated by 270 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove rotate270 \
--original-image=./sample/original.png \
--final-image=./sample/rotated270.png \
--proof-dir=proofs
```

To verify that an image is correctly rotated by 270 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify rotate270 \
--final-image=./sample/rotated270.png \
--proof-dir=proofs
```

## Brighten

To prove that an image is correctly brightened by a brightness factor `f`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove brighten \
--original-image=./sample/original.png \
--final-image=./sample/rotated180.png \
--factor=2 \
--proof-dir=proofs
```

To verify that an image is correctly brightened by a brightness factor `f`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify brighten \
--final-image=./sample/rotated180.png \
--factor=2 \
--proof-dir=proofs
```

## Contrast

To prove that an image is correctly contrasted by a brightness factor `f`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 prove constrast \
--original-image=./sample/original.png \
--final-image=./sample/contrasted.png \
--factor=2 \
--proof-dir=proofs
```

To verify that an image is correctly contrasted by a brightness factor `f`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" mayalabs/maya:v0.0.1 verify contrast \
--final-image=./sample/contrasted.png \
--factor=2 \
--proof-dir=proofs
```
