## Brighten

To run sample code for crop transformation, follow these steps:
1. Clone the [maya-cli](https://github.com/0xmayalabs/maya-cli) repository
    ```shell
    git clone https://github.com/0xmayalabs/maya-cli.git
    ```
2. Cd into the directory
    ```shell
    cd maya-cli
    ```
3. To prove that an image is correctly brightened by a brightness factor `factor`, run:
    ```shell
    docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest prove brighten \
    --original-image=./sample/original.png \
    --final-image=./sample/brightened.png \
    --brightening-factor=2 \
    --proof-dir=proofs
    ```
4. To verify that an image is correctly brightened by a brightness factor `factor`, run:
    ```shell
    docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest verify brighten \
    --final-image=./sample/brightened.png \
    --proof-dir=proofs
    ```

Please note that the repository contains sample images that you can use to get started quickly,
but you don't need to clone the repository to run the `prove` or `verify` commands.