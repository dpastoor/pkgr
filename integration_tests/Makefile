BUILD=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.buildTime=${BUILD}"
# MAKE_HOME=${PWD}
# TEST_HOME=${MAKE_HOME}/integration_tests
BUILD_HOME=${PWD}/../cmd/pkgr
TEST_HOME=${PWD}

.PHONY: install test-setup test-mixed-source-reset test-simple-reset test-threads-reset test-logging-config-reset test-outdated-pkgs-reset test-repo-order-reset test-repo-local-reset test-repo-local-and-remote-reset test-rollback-reset test-simple-suggests-reset test-rollback-reset test-rollback-disabled-reset test-create-library-reset test-strict-mode-reset test-outdated-pkgs-no-update-reset test-misc-reset test-bad-customizations-reset test-description-repo-bug-reset test-load-reset test-load-fail-reset test-recommended-reset test-tarball-install-reset test-tarball-rollback-reset test-tilde-expansion-reset test-binaries-mac-reset test-binaries-linux-reset test-path-bug-reset test-env-vars-reset env-unset

install:
	cd ${BUILD_HOME}; go install ${LDFLAGS}

env-unset:
	export PKGR_RPATH=""

test-setup: test-mixed-source-reset test-simple-reset test-threads-reset test-logging-config-reset test-outdated-pkgs-reset test-outdated-pkgs-stress-reset test-repo-order-reset test-repo-local-reset test-repo-local-and-remote-reset test-simple-suggests-reset test-rollback-reset test-rollback-disabled-reset test-create-library-reset test-strict-mode-reset test-outdated-pkgs-no-update-reset test-misc-reset test-bad-customizations-reset test-description-repo-bug-reset test-load-reset test-load-fail-reset test-recommended-reset test-tarball-install-reset test-tarball-rollback-reset test-tilde-expansion-reset test-binaries-mac-reset test-binaries-linux-reset test-path-bug-reset test-env-vars-reset install

test-mixed-source-reset:
	cd ${TEST_HOME}/mixed-source; rm -rf test-library/
	cd ${TEST_HOME}/mixed-source; rm -rf testsets/
	cd ${TEST_HOME}/mixed-source; mkdir test-library
	cd ${TEST_HOME}/mixed-source; mkdir -p testsets/pkgcache

test-simple-reset:
	cd ${TEST_HOME}/simple; rm -rf test-library/
	cd ${TEST_HOME}/simple; mkdir test-library/

test-threads-reset:
	cd ${TEST_HOME}/threads; rm -rf test-library/
	cd ${TEST_HOME}/threads; mkdir test-library/

test-logging-config-reset:
	cd ${TEST_HOME}/logging-config/install-log; rm -rf logs/
	cd ${TEST_HOME}/logging-config/install-log; rm -rf test-library/
	cd ${TEST_HOME}/logging-config/install-log; mkdir logs
	cd ${TEST_HOME}/logging-config/install-log; mkdir test-library
	cd ${TEST_HOME}/logging-config/default; rm -rf logs/
	cd ${TEST_HOME}/logging-config/default; rm -rf test-library/
	cd ${TEST_HOME}/logging-config/default; mkdir logs
	cd ${TEST_HOME}/logging-config/default; mkdir test-library
	cd ${TEST_HOME}/logging-config/overwrite-setting; rm -rf logs/
	cd ${TEST_HOME}/logging-config/overwrite-setting; rm -rf test-library/
	cd ${TEST_HOME}/logging-config/overwrite-setting; mkdir logs
	cd ${TEST_HOME}/logging-config/overwrite-setting; mkdir test-library
	cd ${TEST_HOME}/logging-config/overwrite-setting; echo "This text should be deleted" > logs/all.log
	cd ${TEST_HOME}/logging-config/overwrite-setting; echo "This text should be deleted" > logs/install.log

test-outdated-pkgs-reset:
	rm -rf ${TEST_HOME}/outdated-pkgs/test-library/
	mkdir ${TEST_HOME}/outdated-pkgs/test-library
	cp -r ${TEST_HOME}/outdated-pkgs/outdated-library/* ${TEST_HOME}/outdated-pkgs/test-library/

test-description-repo-bug-reset:
	rm -rf ${TEST_HOME}/description-repo-bug/test-library/
	mkdir ${TEST_HOME}/description-repo-bug/test-library
	cp -r ${TEST_HOME}/description-repo-bug/outdated-library/* ${TEST_HOME}/description-repo-bug/test-library/

test-outdated-pkgs-no-update-reset:
	rm -rf ${TEST_HOME}/outdated-pkgs-no-update/test-library/
	mkdir ${TEST_HOME}/outdated-pkgs-no-update/test-library
	cp -r ${TEST_HOME}/outdated-pkgs-no-update/outdated-library/* ${TEST_HOME}/outdated-pkgs-no-update/test-library/

test-outdated-pkgs-stress-reset:
	rm -rf ${TEST_HOME}/outdated-pkgs-stress/test-library/
	mkdir ${TEST_HOME}/outdated-pkgs-stress/test-library
	cp -r ${TEST_HOME}/outdated-pkgs-stress/outdated-library/* ${TEST_HOME}/outdated-pkgs-stress/test-library/

test-repo-order-reset:
	rm -rf ${TEST_HOME}/repo-order/test-library/
	mkdir ${TEST_HOME}/repo-order/test-library/

test-repo-local-reset:
	rm -rf ${TEST_HOME}/repo-local/test-library/
	mkdir ${TEST_HOME}/repo-local/test-library/

test-repo-local-and-remote-reset:
	rm -rf ${TEST_HOME}/repo-local-and-remote/test-library/
	mkdir ${TEST_HOME}/repo-local-and-remote/test-library/

test-simple-suggests-reset:
	rm -rf ${TEST_HOME}/simple-suggests/test-library/
	rm -rf ${TEST_HOME}/simple-suggests/pkgcache/
	mkdir ${TEST_HOME}/simple-suggests/test-library
	mkdir ${TEST_HOME}/simple-suggests/pkgcache/

test-rollback-reset:
	cd ${TEST_HOME}/rollback; rm -rf test-library/
	cd ${TEST_HOME}/rollback; mkdir test-library/
	cd ${TEST_HOME}/rollback; cp -r ${TEST_HOME}/rollback/preinstalled-library/* ${TEST_HOME}/rollback/test-library

test-rollback-disabled-reset:
	cd ${TEST_HOME}/rollback-disabled; rm -rf test-library/
	cd ${TEST_HOME}/rollback-disabled; mkdir test-library/
	cd ${TEST_HOME}/rollback-disabled; cp -r ${TEST_HOME}/rollback-disabled/preinstalled-library/* ${TEST_HOME}/rollback-disabled/test-library

test-create-library-reset:
	cd ${TEST_HOME}/create-library; rm -rf test-library/

test-strict-mode-reset:
	cd ${TEST_HOME}/strict-mode; rm -rf test-library/

test-misc-reset:
	cd ${TEST_HOME}/misc; rm -rf test-library/
	cd ${TEST_HOME}/misc; mkdir test-library/
	cd ${TEST_HOME}/misc; rm -rf test-library2/
	cd ${TEST_HOME}/misc; mkdir test-library2/
	cd ${TEST_HOME}/misc; rm -rf localtmp

test-recommended-reset:
	cd ${TEST_HOME}/recommended; rm -rf test-library/
	cd ${TEST_HOME}/recommended; mkdir test-library/

test-bad-customizations-reset:
	rm -rf ${TEST_HOME}/bad-customization/test-library/
	mkdir ${TEST_HOME}/bad-customization/test-library

test-load-reset:
	cd ${TEST_HOME}/load; rm -rf test-library/
	cd ${TEST_HOME}/load; mkdir test-library/
	cd ${TEST_HOME}/load; cp -r ${TEST_HOME}/load/preinstalled-library/* ${TEST_HOME}/load/test-library

test-load-fail-reset:
	cd ${TEST_HOME}/load-fail; rm -rf test-library/
	cd ${TEST_HOME}/load-fail; mkdir test-library/
	cd ${TEST_HOME}/load-fail; cp -r ${TEST_HOME}/load-fail/preinstalled-library/* ${TEST_HOME}/load-fail/test-library

test-tarball-install-reset:
	rm -rf ${TEST_HOME}/tarball-install/test-library/
	rm -rf ${TEST_HOME}/tarball-install/localcache/
	mkdir ${TEST_HOME}/tarball-install/test-library

test-tarball-rollback-reset:
	rm -rf ${TEST_HOME}/tarball-rollback/test-library/
	rm -rf ${TEST_HOME}/tarball-rollback/localcache/
	mkdir ${TEST_HOME}/tarball-rollback/test-library/
	cd ${TEST_HOME}/tarball-rollback; cp -r ${TEST_HOME}/tarball-rollback/preinstalled-library/* ${TEST_HOME}/tarball-rollback/test-library

test-tilde-expansion-reset:
	rm -rf ${TEST_HOME}/tilde-expansion/test-library/
	rm -rf ${TEST_HOME}/tilde-expansion/localcache/
	rm -rf ${TEST_HOME}/tilde-expansion/logs
	mkdir ${TEST_HOME}/tilde-expansion/logs

test-binaries-mac-reset:
	rm -rf ${TEST_HOME}/binaries-mac/test-library
	mkdir ${TEST_HOME}/binaries-mac/test-library

test-binaries-linux-reset:
	rm -rf ${TEST_HOME}/binaries-linux/test-library
	mkdir ${TEST_HOME}/binaries-linux/test-library

test-path-bug-reset:
	rm -rf ${TEST_HOME}/path-bug/test-library
	mkdir ${TEST_HOME}/path-bug/test-library

test-env-vars-reset:
	rm -rf ${TEST_HOME}/env-vars/test-library
	mkdir ${TEST_HOME}/env-vars/test-library
