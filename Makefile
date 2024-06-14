
TARGET=docSyncKit
TARGET_DIR=build
TARGET_Bin=${TARGET_DIR}/bin

default:build

build:
	go build -o ${TARGET_Bin}/${TARGET} src/main.go
	cp -rf rsa_keys ${TARGET_DIR}/
	cp -rf conf ${TARGET_DIR}/

clean:
	rm -rf ${TARGET_DIR}

run:
	./${TARGET_Bin}/${TARGET}