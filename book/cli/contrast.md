## Contrast

To prove that an image is correctly contrasted by a brightness factor `factor`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" labsmaya/maya:latest prove constrast \
--original-image=./sample/original.png \
--final-image=./sample/contrasted.png \
--factor=2 \
--proof-dir=proofs
```

To verify that an image is correctly contrasted by a brightness factor `factor`, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" labsmaya/maya:latest verify contrast \
--final-image=./sample/contrasted.png \
--factor=2 \
--proof-dir=proofs
```
