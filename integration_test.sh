#! /usr/bin/env bash

if [[ ! -e artifacts/out ]] ; then
  echo "ERROR: artifacts/out doesn't exist, try running 'make'"
  exit 1
fi

tempdir=$(mktemp -d "${TMPDIR:-/tmp/}$(basename $0).XXXXXXXXXXXX")

cat > $tempdir/input.json <<EOF
{
  "source": {
    "server": "chat.freenode.net",
    "port": 7070,
    "channel": "#random",
    "user": "randobot1337",
    "password": "secretsecret"
  },
  "params": {
    "message": "This is from \${BUILD_ID}, a.k.a. \${BUILD_NAME}, see \${BUILD_URL}",
    "dry_run": true
  }
}
EOF

cat > $tempdir/expected.json <<EOF
{
  "version": {
    "ref": "none"
  },
  "metadata": [
    {
      "name": "host",
      "value": "chat.freenode.net:7070"
    },
    {
      "name": "channel",
      "value": "#random"
    },
    {
      "name": "user",
      "value": "randobot1337"
    },
    {
      "name": "message",
      "value": "This is from id-123, a.k.a. name-asdf, see https://ci.example.com/teams/team-name-asdf/pipelines/pipeline-name-asdf/jobs/job-name-asdf/builds/name-asdf"
    },
    {
      "name": "dry_run",
      "value": true
    }
  ]
}
EOF

export BUILD_ID=id-123
export BUILD_NAME=name-asdf
export BUILD_JOB_NAME=job-name-asdf
export BUILD_PIPELINE_NAME=pipeline-name-asdf
export BUILD_TEAM_NAME=team-name-asdf
export ATC_EXTERNAL_URL=https://ci.example.com

cat $tempdir/input.json | ./artifacts/out > $tempdir/output.json

expected_checksum=$(sha256sum $tempdir/expected.json | awk '{print $1}')
output_checksum=$(sha256sum $tempdir/output.json | awk '{print $1}')

if [[ $expected_checksum != $output_checksum ]] ; then
  echo "FAIL: output was not what was expected"
  echo "→ expected:"
  cat $tempdir/expected.json
  echo
  echo "→ actual:"
  cat $tempdir/output.json
  exit 1
fi

echo "SUCCESS: output matched expectations"
exit 0
