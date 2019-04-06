package static

const ExampleConfig = `# This is an example .tarreleaser.yml file; please edit accordingly to your needs.
#dist: "dist/"
archive:
  name: "latest.tar.gz"
#  compression_level: 6 # Default to -1 (golang default compression) [1-9]
  wrap_in_directory: "{{.Timestamp}}"
  includes:
    - "./**/*"
  excludes:
    - ".git"
  # Insert a release info file into the archive
  info_file:
    name: "release.txt"
#    content: |
#      Date: {{ .Date }}
#      Tag: {{ .Tag }}
#      Commit: {{ .FullCommit }}

#publish:
#  s3:
#    folder: "my-app/{{.Branch}}"
#    bucket: "my-bucket"
#    region: "eu-west-1"
`

const DefaultReleaseFileContent = `Date: {{ .Date }}
Tag: {{ .Tag }}
Commit: {{ .FullCommit }}
Commit info: {{ .Commit.Message }} - {{ .Commit.Author }}
`
