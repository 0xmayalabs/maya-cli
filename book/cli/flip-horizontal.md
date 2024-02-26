## Flip Horizontal

To prove that an image is correctly flipped horizontally, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" labsmaya/maya:latest prove flip-horizontal \
--original-image=./sample/original.png \
--final-image=./sample/flipped_horizontal.png \
--proof-dir=proofs
```

To verify that an image is correctly flipped horizontally, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" labsmaya/maya:latest verify flip-horizontal \
--final-image=./sample/flipped_horizontal.png \
--proof-dir=proofs
```