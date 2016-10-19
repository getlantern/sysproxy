###############################################################################
#
# This script regenerates the source files that embed the sysproxy-cmd executable.
#
###############################################################################

function die() {
  echo $*
  exit 1
}

if [ -z "$BNS_CERT" ] || [ -z "$BNS_CERT_PASS" ]
then
  die "$0: Please set BNS_CERT and BNS_CERT_PASS to the bns_cert.p12 signing key and the password for that key"
fi

BINPATH=../sysproxy-cmd/binaries

osslsigncode sign -pkcs12 "$BNS_CERT" -pass "$BNS_CERT_PASS" -in $BINPATH/windows/sysproxy_386.exe -out $BINPATH/windows/sysproxy_386.exe || die "Could not sign windows 386"
# osslsigncode sign -pkcs12 "$BNS_CERT" -pass "$BNS_CERT_PASS" -in $BINPATH/windows/sysproxy_amd64.exe -out $BINPATH/windows/sysproxy_amd64.exe || die "Could not sign windows amd64"
go-bindata -nomemcopy -nocompress -pkg sysproxy -prefix $BINPATH/windows -o sysproxy_bytes_windows.go $BINPATH/windows

go-bindata -nomemcopy -nocompress -pkg sysproxy -prefix $BINPATH/linux_386 -o sysproxy_bytes_linux_386.go $BINPATH/linux_386
go-bindata -nomemcopy -nocompress -pkg sysproxy -prefix $BINPATH/linux_amd64 -o sysproxy_bytes_linux_amd64.go $BINPATH/linux_amd64
#go-bindata -nomemcopy -nocompress -pkg sysproxy -prefix $BINPATH/linux_arm -o sysproxy_bytes_linux_arm.go $BINPATH/linux_arm

codesign -s "Developer ID Application: Brave New Software Project, Inc" -f $BINPATH/darwin/sysproxy || die "Could not sign macintosh"
go-bindata -nomemcopy -nocompress -pkg sysproxy -prefix $BINPATH/darwin -o sysproxy_bytes_darwin.go $BINPATH/darwin
