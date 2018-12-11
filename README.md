# IRC Notification Resource for [Concourse](https://concourse.ci)

Sends messages to an IRC channel from a [Concourse CI](https://concourse.ci) pipeline.


## Resource configuration

These parameters go into the `source` fields of the resource type.


### Required

* `server`: The IRC server fully qualified domain name.
* `port`: The TCP port on which to connect.
* `channel`: The name of the channel to be notified (should include leading `#`, e.g., `#go-nuts`).
* `user`: The username used for authentication.
* `password`: The password used for authentication.


### Optional

* `usetls`: Use TLS (a.k.a. SSL) to encrypt your connection to the IRC server. [default: __`true`__]
* `join`: Join the channel before sending the message (and leave afterwards). This is necessary if the channel mode includes `+n`. [default: __`false`__]


## Behaviour

### `check`, `in`

This resource only supports the `put` phase of a job plan, so these are no-ops.


### `out`

Connects to the IRC server, authenticates, and sends the given message to the named channel via a `PRIVMSG` command.


#### Parameters

* `message`: The text of the message to be sent.

Any Concourse [metadata][] in the `message` will be evaluated prior to sending the tweet.

Note also that the pseudo-metadata `BUILD_URL` will expand to:

> `${ATC_EXTERNAL_URL}/teams/${BUILD_TEAM_NAME}/pipelines/${BUILD_PIPELINE_NAME}/jobs/${BUILD_JOB_NAME}/builds/${BUILD_NAME}`

  [metadata]: http://concourse.ci/implementing-resources.html#resource-metadata


## Example usage

``` yml
resource_types:
- name: irc-notification
  type: docker-image
  source:
    repository: flavorjones/irc-notification-resource

resources:
- name: random-channel
  type: irc-notification
  source:
    server: chat.freenode.net
    port: 7070
    channel: "#random"
    user: randobot1337
    password: # your password here

jobs:
- name: post-that-message-to-irc
  plan:
  - put: random-channel
    params:
      message: >
        This is a message about build ${BUILD_ID}, view it at ${BUILD_URL}
```


## Contributing

Pull requests are welcome, as are Github issues opened to discuss bugs or desired features.


### Development and Running the tests

Requires

* `go` >= 1.11
* `make`

``` sh
make test
```

Or if you prefer to use Docker, using the `Dockerfile`, which runs the tests as part of the image build:

```
make docker
```


## License

Distributed under the MIT license, see the `LICENSE` file.
