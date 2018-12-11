
# Checklist for cutting a new release.

- [ ] Update `CHANGELOG.md`
- [ ] Update `README.md` if necessary
- [ ] Commit and push.
- [ ] Create a git tag and push it
- [ ] `make all docker-push`
- [ ] Push the docker image with a specific tag name, e.g. `docker push flavorjones/irc-notification-resource:v1.1.0`
- [ ] Copy README to dockerhub overview
