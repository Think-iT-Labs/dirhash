# Dirhash
Calculating the checksum of a directory made easy.

## Dirhash CLI
A CLI that calculates the checksum of a directory

### Usage
```sh
dirhash sha256 [--ignore=comma,delimited,list,of,patterns] directory_path
```


### How to build
 1. First download dependencies
 ```sh
 go mod download
 ```
 2. Build the CLI
 ```sh
 make build-cli
 ```