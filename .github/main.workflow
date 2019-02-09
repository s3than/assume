workflow "New workflow" {
  on = "push"
  resolves = ["GitHub Action for Docker-1"]
}

action "Filters for GitHub Actions" {
  uses = "actions/bin/filter@9d4ef995a71b0771f438dd7438851858f4a55d0c"
  args = "tag"
}

action "GitHub Action for Docker" {
  uses = "actions/docker/cli@aea64bb1b97c42fa69b90523667fef56b90d7cff"
  needs = ["Filters for GitHub Actions"]
  args = "build -t s3than/assume:${GITHUB_REF} ."
}

action "Docker Login" {
  uses = "actions/docker/login@master"
  needs = ["GitHub Action for Docker"]
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "GitHub Action for Docker-1" {
  uses = "actions/docker/cli@aea64bb1b97c42fa69b90523667fef56b90d7cff"
  needs = ["Docker Login"]
  args = "push s3than/assume:${GITHUB_REF}"
}
