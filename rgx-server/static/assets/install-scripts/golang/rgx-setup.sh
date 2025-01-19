#!/usr/bin/env sh
set -e

mv go go-${RGX_PACKAGE_VERSION}
export GO_INSTALL_DIR=${RGX_PACKAGE_SCRIPTDIR}/go-${RGX_PACKAGE_VERSION}
echo go-${RGX_PACKAGE_VERSION} is now present in ${GO_INSTALL_DIR}

export PATH=${RGX_PACKAGS_DIR}/golang/go-${RGX_PACKAGE_VERSION}/bin:\$PATH

echo Writing .golang rc file...
echo export GOLANG_HOME=${GO_INSTALL_DIR}> ${RGX_RCFILE_DIR}/.golang-$RGX_PACKAGE_MAJORVERSION-rc
echo export PATH=\$GOLANG_HOME/bin:$PATH>> ${RGX_RCFILE_DIR}/.golang-$RGX_PACKAGE_MAJORVERSION-rc
echo eecho \"Go ${RGX_PACKAGE_MAJORVERSION} added to PATH\">> ${RGX_RCFILE_DIR}/.golang-$RGX_PACKAGE_MAJORVERSION-rv
echo To use Go, source ${RGX_RCFILE_DIR}/.golang-${RGX_PACKAGE_MAJORVERSION}-rc