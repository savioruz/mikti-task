#!/usr/bin/env sh

export CUR="github.com/savioruz/mikti-task/tree/week-4"
export NEW="github.com/savioruz/mikti-task"
go mod edit -module ${NEW}
find . -type f -name '*.go' -exec perl -pi -e 's/$ENV{CUR}/$ENV{NEW}/g' {} \;