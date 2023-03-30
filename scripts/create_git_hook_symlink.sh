#!/bin/bash

FILE_PATH="$(dirname "$(realpath $0)")"

ln $FILE_PATH/git_hooks/commit_msg $FILE_PATH/../.git/hooks/commit_msg