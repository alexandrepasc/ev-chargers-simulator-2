#!/bin/bash

FILE_PATH="$(dirname "$(realpath $0)")"

ln -s $FILE_PATH/git_hooks/commit-msg $FILE_PATH/../.git/hooks/commit-msg