## Crop

To run sample code, follow these steps:
1. Clone the [maya-cli](https://github.com/0xmayalabs/maya-cli) repository
    ```shell
    git clone https://github.com/0xmayalabs/maya-cli.git
    ```
2. Cd into the directory
    ```shell
    cd maya-cli
    ```
3. To prove that an image is cropped correctly, run:
   ```shell
   docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest prove crop \
   --original-image=./sample/original.png \
   --cropped-image=./sample/cropped.png \
   --height-start-new=0 \
   --width-start-new=0 \
   --proof-dir=proofs
   ```
4. To verify that an image is cropped correctly, run:
```shell
docker run --rm -v "$(pwd):/opt/maya" 0xmayalabs/maya-cli:latest verify crop \
--cropped-image=./sample/cropped.png \
--proof-dir=proofs
```

Please note that the repository contains sample images that you can use to get started quickly, 
but you don't need to clone the repository to run the `prove` or `verify` commands.



