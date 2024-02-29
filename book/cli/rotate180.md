## Rotate 180

To prove that an image is correctly rotated by 180 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya:latest prove rotate180 \
--original-image=./sample/original.png \
--final-image=./sample/rotated180.png \
--proof-dir=proofs
```

To verify that an image is correctly rotated by 180 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya:latest verify rotate180 \
--final-image=./sample/rotated180.png \
--proof-dir=proofs
```
