#! /usr/bin/env bash
set -o errexit \
    -o pipefail

if [[ -n "${DEBUG}" ]]; then
  set -o xtrace \
      -o verbose
fi

if [[ ! -e artifacts/out ]] ; then
  echo "ERROR: artifacts/out doesn't exist, try running 'make'"
  exit 1
fi

tempdir=$(mktemp -d "${TMPDIR:-/tmp/}$(basename $0).XXXXXXXXXXXX")

resource_version=$(fgrep "const VERSION" pkg/irc/irc.go  | cut -d'"' -f2)

for test_name in $(ls test/integration/*.input.json | sed 's/.input.json//' | uniq) ; do
  echo "testing ${test_name}.input.json (in $tempdir)"
  cat ${test_name}.input.json | sed "s/\${resource_version}/${resource_version}/" > $tempdir/input.json
  cat ${test_name}.expected.json | sed "s/\${resource_version}/${resource_version}/" > $tempdir/expected.json

  export BUILD_ID=id-123
  export BUILD_NAME=name-asdf
  export BUILD_JOB_NAME=job-name-asdf
  export BUILD_PIPELINE_NAME=pipeline-name-asdf
  export BUILD_TEAM_NAME=team-name-asdf
  export ATC_EXTERNAL_URL=https://ci.example.com

  cat $tempdir/input.json | ./artifacts/out $(pwd) > $tempdir/output.json

  expected_checksum=$(sha256sum $tempdir/expected.json | awk '{print $1}')
  output_checksum=$(sha256sum $tempdir/output.json | awk '{print $1}')

  if [[ $expected_checksum != $output_checksum ]] ; then
    echo "FAIL: output was not what was expected"
    echo "→ diff:"
    diff -u $tempdir/expected.json $tempdir/output.json
    exit 1
  fi
done

echo "SUCCESS: output matched expectations"
exit 0
