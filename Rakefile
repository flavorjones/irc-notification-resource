require "concourse"
require "hoe"

Hoe.plugin :markdown

Concourse.new("irc-notification-resource", directory: "ci", fly_target: "flavorjones").create_tasks!

Hoe.spec("irc-notification-resource") do
  developer "Mike Dalessio", "mike.dalessio@gmail.com"

  self.urls = {
    "bugs" => "https://github.com/flavorjones/irc-notification-resource/issues",
  }
end

#
#  build a docker image with test dependencies
#
namespace "docker" do
  DOCKER_TAG = "flavorjones/irc-notification-resource:test"
  desc "Build a docker image for testing"
  task "build" do
    sh "docker build -t #{DOCKER_TAG} -f ci/images/Dockerfile ."
  end

  desc "Push a docker image for testing"
  task "push" do
    sh "docker push  #{DOCKER_TAG}"
  end
end

desc "Build and push a docker image for testing"
task "docker" => ["docker:build", "docker:push"]
