BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.buildTime=${BUILD}"
MAKE_HOME=${PWD}
TEST_HOME=${MAKE_HOME}/integration_tests

.PHONY: install test-multiple log-test log-test-reset

install:
	cd cmd/pkgr; go install ${LDFLAGS}

test-multiple:	
	cd ${TEST_HOME}

	rm -rf master/test-library/*
	rm -rf mixed-source/test-library/*
	rm -rf pull-source-for-one/test-library/*
	rm -rf simple/test-library/*
	rm -rf simple-suggests/test-library/*

	-cd ${TEST_HOME}/master; pkgr install

	-cd ${TEST_HOME}/mixed-source; pkgr install

	-cd ${TEST_HOME}/pull-source-for-one; pkgr install

	-cd ${TEST_HOME}/simple;	pkgr install

	#-cd ${TEST_HOME}/simple-suggests; pkgr install

log-test: install
	cd ${TEST_HOME}/logging-config/install-log; pkgr install
	cd ${TEST_HOME}/logging-config/default; pkgr install

log-test-reset:
	cd ${TEST_HOME}/logging-config/install-log; rm -rf logs/*
	cd ${TEST_HOME}/logging-config/default; rm -rf logs/*