# Dirhash
Calculating the checksum of a directory made easy.

## Dirhash CLI
A CLI that calculates the checksum of a directory

### Usage
```sh
dirhash sha256 [--ignore=comma,delimited,list,of,patterns] directory_path
```

## Dirhash Provider

Dirhash was created to calculate checksum of a directory, 
a useful provider if you try to make Terraform react and make changes 
based on whether changes have been made inside a directory or not.

### Usage

As noted above, the Dirhash provider can calculate the checksum of a directory, 
the data_source `dirhash_sha256` can be used to retrieve the SHA256 of a directory.

For example:

```terraform
provider "dirhash" {}
```
```terraform
data "dirhash_sha256" "example" {
  directory = "/path/to/directory"
  ignore = [
    "glob_pattern_1/*",
    "glob_pattern_2/*"
  ]
}

output "directory_sha256_checksum" {
  description = "the Output SHA256 Checksum for the directory"
  value       = data.dirhash_sha256.example.checksum
}
```

## How to build
1. First download dependencies
```sh
go mod download
```
2. Build the CLI
```sh
make build-cli
```
Your CLI binary will be under `bin/dirhash`
3. Build the provider
```
make build-provider
```
Your Terraform provider binary will be under `bin/dirhash-terraform-provider`
