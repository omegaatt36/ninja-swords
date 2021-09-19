#!/bin/bash

set -ex  # Exit on error; debugging enabled.
set -o pipefail  # Fail a pipe if any sub-command fails.

# not makes sure the command passed to it does not exit with a return code of 0.
not() {
  # This is required instead of the earlier (! $COMMAND) because subshells and
  # pipefail don't work the same on Darwin as in Linux.
  ! "$@"
}

die() {
  echo "$@" >&2
  exit 1
}

check_status() {
  # Check to make sure it's safe to modify the user's git repo.
  local out=$(git status --porcelain)
  if [ ! -z "$out" ]; then
    echo "status not clean"
    echo $out
    exit 1
  fi
}

check_status


# Undo any edits made by this script.
cleanup() {
  git reset --hard HEAD
}
trap cleanup EXIT

PATH="${GOPATH}/bin:${GOROOT}/bin:${PATH}"

if [[ "$1" = "-install" ]]; then
  # Check for module support
  if go help mod >& /dev/null; then
    pushd ./test/tools
    # Install the pinned versions as defined in module tools.
    go install \
      golang.org/x/tools/cmd/goimports \
      github.com/client9/misspell/cmd/misspell \
      github.com/gogo/protobuf/protoc-gen-gogoslick
    go install \
      honnef.co/go/tools/cmd/staticcheck@2022.1.3
    go install \
      github.com/mgechev/revive@v1.2.3
    go install \
      github.com/daixiang0/gci@latest
    go install \
      github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0
    popd
  else
    echo "we don't support old go get anymore"
    exit 1
  fi
  exit 0
elif [[ "$#" -ne 0 ]]; then
  die "Unknown argument(s): $*"
fi

# - gofmt, goimports, golint (with exceptions for generated code), go vet.
# gofmt -s -d -l . 2>&1 || exit 1
# goimports -l . 2>&1 | not grep -vE "(_mock|\.pb)\.go"
# revive -exclude pkg/... -formatter friendly -config test/tools/revive.toml  ./...
# go vet -all ./...

# misspell -error */**

make fmt check &&
    check_status || \
    (git status; git --no-pager diff; exit 1)

set +x

echo SUCCESS
