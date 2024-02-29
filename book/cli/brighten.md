## Brighten

To prove that an image is correctly brightened by a brightness factor `factor`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest prove brighten \
--original-image=./sample/original.png \
--final-image=./sample/brightened.png \
--brightening-factor=2 \
--proof-dir=proofs
```

To verify that an image is correctly brightened by a brightness factor `factor`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest verify brighten \
--final-image=./sample/brightened.png \
--proof-dir=proofs
```