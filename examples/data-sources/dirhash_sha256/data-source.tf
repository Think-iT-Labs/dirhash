data "dirhash_sha256" "example" {
  directory = "/path/to/directory"
  ignore = [
    "glob_pattern_1/*",
    "glob_pattern_2/*"
  ]
}
