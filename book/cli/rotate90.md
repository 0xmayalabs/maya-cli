## Rotate 90

To prove that an image is correctly rotated by 90 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" labsmaya/maya:latest prove rotate90 \
--original-image=./sample/original.png \
--final-image=./sample/rotated90.png \
--proof-dir=proofs
```

To verify that an image is correctly rotated by 90 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" labsmaya/maya:latest verify rotate90 \
--final-image=./sample/rotated90.png \
--proof-dir=proofs
```