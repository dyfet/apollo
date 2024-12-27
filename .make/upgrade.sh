#!/bin/bash
# Copyright (C) 2024 Tycho Softworks.
#
# This file is free software; as a special exception the author gives
# unlimited permission to copy and/or distribute it, with or without
# modifications, as long as this notice is preserved.
#
# This program is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY, to the extent permitted by law; without even the
# implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

outdated_deps=$(go list -u -m -mod=mod all | grep '

\[.*\]

')

while IFS= read -r dep; do
    module_path=$(echo "$dep" | awk '{print $1}')
    current_version=$(echo "$dep" | awk '{print $2}')
    latest_version=$(echo "$dep" | awk '{print $NF}' | tr -d '[]')
    if [ "$current_version" != "$latest_version" ]; then
        echo "Updating $module_path from $current_version to $latest_version"
        go get -u "$module_path@$latest_version" >/dev/null
    fi
done <<< "$outdated_deps"
go mod tidy
go mod vendor









#!/bin/bash
deps=$(go list -mod=mod -m all | grep -v github.com/golang | grep -v golang.org | awk '{print $1}')
for dep in $deps; do
    current_version=$(go list -mod=mod -m $dep | awk '{print $2}')
    latest_version=$(go list -mod=mod -u -m $dep | awk '{print $2}')
    echo $dep $current_version $latest_version

    if [ "$current_version" != "$latest_version" ]; then
        echo "Updating $dep from $current_version to $latest_version"
        go get -u $dep
    fi
done
go mod tidy

