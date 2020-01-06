#!/bin/bash

#**
#* update_version.sh
#*
#* Copyright (C) 2019 by Liu YunFeng.
#*
#*        Create : 2019-12-05 15:04:48
#* Last Modified : 2019-12-05 15:04:48
#
#
# TODO: read current version and tips next version
#       input a new version if choose no
# 
#**

declare -i VERSION_MAJOR=0
declare -i VERSION_MINOR=0
declare -i VERSION_PATCH=0
VERSION_FILE=pkg/version/version.go
VERSION_MAKEFILE=Makefile

function generate_version_file(){
	cat > ${VERSION_FILE} <<-EOF
	package version

	import "fmt"

	const (
	Major = ${VERSION_MAJOR}
	Minor = ${VERSION_MINOR}
	Patch = ${VERSION_PATCH}
	)

	func Version() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}
EOF

if [ "$(command -v gofmt)" != "" ]; then
	gofmt -w ${VERSION_FILE}
fi
}

function list_version_tag() {
	if [ "$(command -v git)" != "" ]; then
		git tag
	fi
}

function read_version(){
	read -p "Major version: " -n2 VERSION_MAJOR
	read -p "Minor version: " -n2 VERSION_MINOR
	read -p "Patch version: " -n2 VERSION_PATCH

}

function update_makefile(){
	eval sed -i '/^MAJOR=*/c\MAJOR=${VERSION_MAJOR}' ${VERSION_MAKEFILE}
	eval sed -i '/^MINOR=*/c\MINOR=${VERSION_MINOR}' ${VERSION_MAKEFILE}
	eval sed -i '/^PATCH=*/c\PATCH=${VERSION_PATCH}' ${VERSION_MAKEFILE}
}

list_version_tag
read_version
update_makefile
generate_version_file

echo
echo Version updated, commit the modied code before release new version ${VERSION_MAJOR}.${VERSION_MINOR}.${VERSION_PATCH}
echo
