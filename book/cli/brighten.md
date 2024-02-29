## Brighten

To prove that an image is correctly brightened by a brightness factor `factor`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya:latest prove brighten \
--original-image=./sample/original.png \
--final-image=./sample/brightened.png \
--factor=2 \
--proof-dir=proofs
```

To verify that an image is correctly brightened by a brightness factor `factor`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya:latest verify brighten \
--final-image=./sample/brightened.png \
--factor=2 \
--proof-dir=proofs
```