
# Checklist for cutting a new release.

- [ ] Update `CHANGELOG.md`
- [ ] Update `README.md` if necessary
- [ ] Bump version in `pkg/irc/irc.go`
- [ ] Commit and push.
- [ ] Create a git tag and push it
- [ ] `make all` to create a docker image
- [ ] Tag the docker image, e.g. `docker tag flavorjones/irc-notification-resource:latest flavorjones/irc-notification-resource:v1.1.0`
- [ ] `make docker-push`
- [ ] Copy README to [dockerhub overview](https://cloud.docker.com/repository/docker/flavorjones/irc-notification-resource/general)
- [ ] Create a [github release](https://github.com/flavorjones/irc-notification-resource/releases) with CHANGELOG snippet as description
- [ ] Check that the resource works by kicking off the [`test-notification` job](https://ci.nokogiri.org/teams/flavorjones/pipelines/irc-notification-resource/jobs/test-notification/builds/3)
