#!/usr/bin/env bash

###############################################################################
#
# This script regenerates the source files that embed the sysproxy-cmd executable.
#
###############################################################################

set -euo pipefail
function die() {
  echo "$*"
  exit 1
}

if [ -z "$BNS_CERT" ] || [ -z "$BNS_CERT_PASS" ]
then
  die "$0: Please set BNS_CERT and BNS_CERT_PASS to the bns_cert.p12 signing key and the password for that key"
fi

BINPATH=../sysproxy-cmd/binaries

# Check for a recent version of osslsigncode that can handle 32-bit Windows
# binaries.
osslsigncode_version=$(osslsigncode --version 2>&1 | grep "using:" | cut -d " " -f 2 | cut -d "," -f 1)
if [[ "$osslsigncode_version" < "1.7.1" ]]
then
  die "$0: Please upgrade osslsigncode to at least version 1.7.1"
fi

osslsigncode sign -pkcs12 "$BNS_CERT" -pass "$BNS_CERT_PASS" -in $BINPATH/windows/sysproxy_386.exe -out $BINPATH/windows/sysproxy_386.exe || die "Could not sign windows 386"
osslsigncode sign -pkcs12 "$BNS_CERT" -pass "$BNS_CERT_PASS" -in $BINPATH/windows/sysproxy_amd64.exe -out $BINPATH/windows/sysproxy_amd64.exe || die "Could not sign windows 386"
cp $BINPATH/windows/sysproxy_386.exe binaries/windows
cp $BINPATH/windows/sysproxy_amd64.exe binaries/windows

codesign --options runtime --strict --timestamp --force --deep -r="designated => anchor trusted and identifier com.getlantern.lantern" -s "Developer ID Application: Innovate Labs LLC (4FYC28AXA2)" -v $BINPATH/darwin
cp $BINPATH/darwin binaries
