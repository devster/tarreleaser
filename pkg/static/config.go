package static

const ExampleConfig = `# This is an example .tarreleaser.yaml file with some sane defaults.
#dist: "./dist"
archive:
  #compression_level: 6
  wrap_in_directory: "myapp"
  includes:
    - "src/**/*"

publish:

`
