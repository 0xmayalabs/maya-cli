## Flip vertical

To prove that an image is correctly flipped vertically, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya:latest prove flip-vertical \
--original-image=./sample/original.png \
--final-image=./sample/flipped_vertical.png \
--proof-dir=proofs
```

To verify that an image is correctly flipped vertically, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya:latest verify flip-vertical \
--final-image=./sample/flipped_vertical.png \
--proof-dir=proofs
```