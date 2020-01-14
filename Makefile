.DEFAULT:all
.PHONY:clean

MAJOR=1
MINOR=0
PATCH=7
VER=v$(MAJOR).$(MINOR).$(PATCH)

PKG=zktnotify
LINUX-AMD64=$(PKG)-$(VER)-linux-amd64
WINDOWS-AMD64=$(PKG)-$(VER)-windows-amd64
SQLITE-LINUX-AMD64=$(PKG)-$(VER)-linux-amd64-sqlite

UPX=$(shell which upx)
SRC=$(shell find . -name "*.go" -o -name [Mm]akefile)

LDFLAGS=-ldflags "-s -w"

origin:$(SRC)
	go build -o $(PKG)

all: origin build
build: linux-amd64 windows-amd64

linux-amd64:$(SRC)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(LINUX-AMD64)
ifneq ("$(UPX)","")
	$(UPX) -9 $(LINUX-AMD64)
endif

windows-amd64:$(SRC)
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(WINDOWS-AMD64)
ifneq ("$(UPX)","")
	$(UPX) -9 $(WINDOWS-AMD64)
endif

sqlite:$(SRC)
	go build -tags sqlite3 -o $(SQLITE-LINUX-AMD64)
ifneq ("$(UPX)","")
	$(UPX) -9 $(SQLITE-LINUX-AMD64)
endif

release:all
	-@./$(PKG) release
	-@echo release finished...

upgrade:
	-@./$(PKG) upgrade
	-@echo upgrade finished...

update-version:
	-@./update_version.sh

fixme:
	-@grep -rnw "FIXME" cmd pkg router models

todo:
	-@grep -rnw "TODO" cmd pkg router models

clean:
	@-go clean
	@-rm -rf $(PKG)
	@-rm -rf $(LINUX-AMD64)
	@-rm -rf $(WINDOWS-AMD64)
	@-rm -rf $(SQLITE-LINUX-AMD64)
