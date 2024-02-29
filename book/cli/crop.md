## Crop

To prove that an image is cropped correctly, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest prove crop \
--original-image=./sample/original.png \
--cropped-image=./sample/cropped.png \
--height-start-new=0 \
--width-start-new=0 \
--proof-dir=proofs
```

To verify that an image is cropped correctly, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest verify crop \
--cropped-image=./sample/cropped.png \
--proof-dir=proofs
```