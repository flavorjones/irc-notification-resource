
# Checklist for cutting a new release.

- [ ] Update `CHANGELOG.md`
- [ ] Update `README.md` if necessary
- [ ] Commit and push.
- [ ] Create a git tag and push it
- [ ] `make all` to create a docker image
- [ ] Tag the docker image, e.g. `docker tag flavorjones/irc-notification-resource:latest flavorjones/irc-notification-resource:v1.1.0`
- [ ] `make docker-push`
- [ ] Copy README to dockerhub overview
- [ ] Create a github release with CHANGELOG snippet as description
