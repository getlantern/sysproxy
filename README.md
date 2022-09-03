# sysproxy

[sysproxy](https://github.com/getlantern/sysproxy) is a simple Go library to
toggle the system proxy on and off for Windows and MacOS. It will
extract a helper tool and use it to actually change the system proxy settings.

```go
sysproxy.EnsureHelperToolPresent(fullPath, prompt, iconFullPath)
sysproxy.On(proxyAddr string)
sysproxy.Off()
```

See 'example/main.go' for detailed usage.

### Embedding sysproxy-cmd

sysproxy uses binaries from the
[sysproxy-cmd](https://github.com/getlantern/sysproxy-cmd) project and from
[sysproxy-cmd-darwin](https://github.com/getlantern/sysproxy-cmd-darwin).

To embed the binaries for different platforms, use the `sysproxy2go.sh` script.
This script takes care of code signing the Windows and MacOS executables.

This script signs the Windows executable, which requires that
[osslsigncode](http://sourceforge.net/projects/osslsigncode/) utility be
installed. On OS X with homebrew, you can do this with
`brew install osslsigncode`.

You will also need to set the environment variables BNS_CERT and BNS_CERT_PASS
to point to [bns-cert.p12](https://github.com/getlantern/too-many-secrets/blob/master/bns_cert.p12)
and its [password](https://github.com/getlantern/too-many-secrets/blob/master/build-installers/env-vars.txt#L3)
so that the script can sign the Windows executable.

This script also signs the MacOS executable, which requires you to install a MacOS signing certificate.
