# IRC Notification Resource for [Concourse](https://concourse.ci)

Sends messages to an IRC channel from a [Concourse CI](https://concourse.ci) pipeline.


## Resource configuration

These parameters go into the `source` fields of the resource type.


### Required

* `server`: the IRC server domain name
* `port`: the TCP port on which to connect
* `channel`: the name of the channel to be notified (should include leading `#`, e.g., `#go-nuts`)
* `user`: the username used for authentication
* `password`: the password used for authentication


### Optional

* `use_tls`: use TLS (a.k.a. SSL) to encrypt your connection (default: __`true`__)
* `nick`: the publicly visible nickname for the connection (default: same as `user` parameter)


## Behaviour

### `check`, `in`

This resource only supports the `put` phase of a job plan, so these
are effectively no-ops.


### `out`

Connects to the IRC server, authenticates, and sends the given message
to the named channel via a `PRIVMSG` command.


#### Parameters

* `message`: The text of the message to be sent.

Any Concourse [metadata][] in the `message` will be evaluated prior to
sending the tweet.

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


## License

Distributed under the MIT license, see the `LICENSE` file.
