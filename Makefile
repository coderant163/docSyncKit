
TARGET=docSyncKit
TARGET_DIR=build
TARGET_Bin=${TARGET_DIR}/bin

#git info
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
COMMITID=$(shell git rev-parse HEAD)
BUILDTIME=$(shell date +'%Y-%m-%d %H:%M:%S')
LDFLAGS=-ldflags "-X 'github.com/coderant163/docSyncKit/src/cmd.Branch=$(BRANCH)' -X 'github.com/coderant163/docSyncKit/src/cmd.CommitID=$(COMMITID)' -X 'github.com/coderant163/docSyncKit/src/cmd.BuildTime=$(BUILDTIME)'"

default:build

build:
	go build ${LDFLAGS} -o ${TARGET_Bin}/${TARGET} src/main.go
	cp -rf rsa_keys ${TARGET_DIR}/
	cp -rf conf ${TARGET_DIR}/
	cp -rf sbin ${TARGET_DIR}/
	cp README.md ${TARGET_DIR}/

clean:
	rm -rf ${TARGET_DIR}

run:
	./${TARGET_Bin}/${TARGET}