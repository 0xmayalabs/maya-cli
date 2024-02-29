# Install

To run the commands, you can either use docker or build from source.

## Docker

Install [docker](https://docs.docker.com/engine/install/) and run these [commands](./cli/runmaya.md) using docker.

## Build from source

1. Build and install the binary
    ```shell
    go install .
    ```
2. Test if it is installed properly
    ```
   maya-cli --help
   ```
