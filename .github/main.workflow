workflow "New workflow" {
  on = "push"
  resolves = [
    # "Docker Build",
    # "Docker Tag",
    # "Docker Push",
    "goreleaser",
  ]
}

# action "Filters" {
#   uses = "actions/bin/filter@9d4ef995a71b0771f438dd7438851858f4a55d0c"
#   args = "tag"
# }

# action "Docker Login" {
#   uses = "actions/docker/login@master"
#   needs = ["Filters"]
#   secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
# }

# action "Docker Build" {
#   uses = "actions/docker/cli@aea64bb1b97c42fa69b90523667fef56b90d7cff"
#   needs = ["Docker Login"]
#   args = "build -t assume ."
# }

# action "Docker Tag" {
#   uses = "actions/docker/tag@aea64bb1b97c42fa69b90523667fef56b90d7cff"
#   needs = ["Docker Build"]
#   args = "assume s3than/assume"
# }

# action "Docker Push" {
#   uses = "actions/docker/cli@aea64bb1b97c42fa69b90523667fef56b90d7cff"
#   needs = ["Docker Tag"]
#   args = "push s3than/assume"
# }

action "goreleaser" {
  uses = "docker://goreleaser/goreleaser"
  # needs = ["Filters"]
  args = "release"
  secrets = ["GITHUB_TOKEN"]
}
