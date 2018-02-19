#! /usr/bin/env bash

export BUILD_ID=id-123
export BUILD_NAME=name-asdf
export BUILD_JOB_NAME=job-name-asdf
export BUILD_PIPELINE_NAME=pipeline-name-asdf
export BUILD_TEAM_NAME=team-name-asdf
export ATC_EXTERNAL_URL=https://ci.example.com

cat cmd/out/test_input.json | ./artifacts/out
