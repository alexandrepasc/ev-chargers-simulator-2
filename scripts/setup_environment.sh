#!/bin/bash

FILE_PATH="$(dirname "$(realpath $0)")"

echo "Start golangci-lint installation"
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
echo "Finished golangci-lint installation"

echo ""

echo "Copy commit-msg hook to .git folder"
ln -s $FILE_PATH/git_hooks/commit-msg $FILE_PATH/../.git/hooks/commit-msg
echo "Finished copiying commit-msg to .git folder"

echo ""

echo "Copy pre-commit hook to .git folder"
ln -s $FILE_PATH/git_hooks/pre-commit $FILE_PATH/../.git/hooks/pre-commit
echo "Finished copiying pre-commit to .git folder"