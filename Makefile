PRJ=nfs-hotpot-regist
PROG_REGIST=${PRJ}
PROG_BOOTKIT=nfs-boot-kit
PROG_DBUS=org.nfs.HotpotRegist1
PREFIX=/usr
VAR=/var/lib
PWD=$(shell pwd)
GOCODE=/usr/share/gocode
GOPATH_DIR=gopath
CURRENT_DIR=$(notdir $(shell pwd))
ARCH=$(shell arch)
export GO111MODULE=off

all: build

prepare:
	@if [ ! -d ${GOPATH_DIR}/src/${PRJ} ]; then \
		mkdir -p ${GOPATH_DIR}/src/${PRJ}; \
		ln -sf ${PWD}/pkg ${GOPATH_DIR}/src/${PRJ}; \
		ln -sf ${PWD}/cmd ${GOPATH_DIR}/src/${PRJ}/; \
	fi

$(info, $(GOPATH))
$(warning, $(GOPATH))
$(error, $(GOPATH))


build: prepare
	mkdir -p ${PWD}/out
	@env GOPATH=${PWD}/${GOPATH_DIR}:${GOCODE}:${GOPATH} ls -al ${PWD}/${GOPATH_DIR}/src/nfs-hotpot-regist/*
	@env GOPATH=${PWD}/${GOPATH_DIR}:${GOCODE}:${GOPATH} ls -al ${PWD}/${GOPATH_DIR}/src/nfs-hotpot-regist/cmd/*

	@env GOPATH=${PWD}/${GOPATH_DIR}:${GOCODE}:${GOPATH}  \
	CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2"  CGO_LDFLAGS="-Wl,-z,relro,-z,now" \
	go build -o ${PWD}/${PROG_REGIST} $(GO_BUILD_FLAGS) ${PRJ}/cmd

install-nfs:
	@mkdir -p ${DESTDIR}${PREFIX}/share/dbus-1/system.d/
	@cp -f ${PWD}/configs/dbus/${PROG_DBUS}.conf  ${DESTDIR}${PREFIX}/share/dbus-1/system.d/

	@mkdir -p ${DESTDIR}${PREFIX}/share/dbus-1/system-services/
	@cp -f ${PWD}/configs/dbus/${PROG_DBUS}.service  ${DESTDIR}${PREFIX}/share/dbus-1/system-services/

	@mkdir -p ${DESTDIR}${PREFIX}/sbin
	@cp -f ${PWD}/${PROG_REGIST} ${DESTDIR}${PREFIX}/sbin

install: install-nfs

uninstall-nfs:
	@rm -f ${DESTDIR}${PREFIX}/sbin/${PROG_REGIST}
	@rm -f ${DESTDIR}${PREFIX}/share/dbus-1/system.d/${PROG_DBUS}.conf
	@rm -f ${DESTDIR}${PREFIX}/share/dbus-1/system-services/${PROG_DBUS}.service

uninstall: uninstall-nfs

clean:
	@rm -rf ${GOPATH_DIR}
	@rm -rf ${PWD}/${PROG_REGIST}
	@rm -rf ${PWD}/out

rebuild: clean build
