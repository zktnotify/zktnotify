#
# Makefile
#
# Copyright (C) 2019 by Liu YunFeng.
#
#        Create : 2019-11-25 15:43:21
# Last Modified : 2019-11-25 15:43:21
#

RED_COLOR = "\033[1;31m"
GRN_COLOR = "\033[1;32m"
YEL_COLOR = "\033[1;33m"
BLU_COLOR = "\033[1;34m"
GRA_COLOR = "\033[1;35m"
SKY_COLOR = "\033[1;36m"
WHT_COLOR = "\033[1;37m"
RST_COLOR = "\033[0m"

DIALOG = echo
QUIET_REMOVE = $(QUIET_RM)rm -f

QUIET_RM  = @printf '%b %b\n' $(RED_COLOR)REMOVE$(RST_COLOR) $(GRN_COLOR)$@$(RST_COLOR) 1>&2;

pkg=zktnotify
path=github.com/zktnotify/zktnotify/pkg/version

src=$(shell find . -name "*.go" -o -name [Mm]akefile)

version=v1.0.1
branch=$(shell cat .git/HEAD | awk -F"/" '{print $$3}' )
buildTime=$(shell date "+%Y-%m-%d %H:%M:%S")
commitID=$(shell git rev-parse --short HEAD || echo unsupported)
buildExtTag="-X $(path).version=$(version) -X '$(path).buildTime=$(buildTime)' -X $(path).commitID=$(commitID) -X '$(path).branch=$(branch)'"

default:help

build:$(pkg)
$(pkg):$(src)
	@$(DIALOG) $(GRN_COLOR)BUILDING $(WHT_COLOR)$(pkg) $(RST_COLOR)
	@GOOS=windows GOARCH=amd64 go build -ldflags $(buildExtTag) -o "$@".windows.amd64

all:upload
deploy:upload
upload:$(pkg)
	@$(DIALOG) $(GRN_COLOR)UPLOADING $(WHT_COLOR)$(pkg) $(RST_COLOR)

help:
	-@$(DIALOG) $(GRN_COLOR)build/$(pkg)$(RST_COLOR) - $(WHT_COLOR)build the project$(RST_COLOR)
	-@$(DIALOG) $(GRN_COLOR)all/deploy/upload$(RST_COLOR) - $(WHT_COLOR)upload program to github.com$(RST_COLOR)

.PHONY:clean
clean:
	$(QUIET_REMOVE) $(pkg)
