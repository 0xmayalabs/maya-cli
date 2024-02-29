## Rotate 270

To prove that an image is correctly rotated by 270 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest prove rotate270 \
--original-image=./sample/original.png \
--final-image=./sample/rotated270.png \
--proof-dir=proofs
```

To verify that an image is correctly rotated by 270 degrees clockwise, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest verify rotate270 \
--final-image=./sample/rotated270.png \
--proof-dir=proofs
```