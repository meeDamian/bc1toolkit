bc1isup
=======

A minimal & focused unix-style tool to quickly check the status and get basic info about one or more Bitcoin nodes.  


### Usage:

```
$ bc1isup --help

Usage:
  bc1isup [OPTIONS] (domain|IP)[:port] ...

Checks addresses for running Bitcoin nodes. When addresses are both piped-in and provided at command line, piped ones are first.

Each address provided, outputs its own line with corresponding node status info (customizable with --output=?).
Exit code of 0 is returned only if each address provided had at least one running node.

Tor "auto" behaviour: try using Tor, if not available, fall back to clearnet.

Application Options:
  -v, --version                             Show version and exit
  -V, --verbose                             Enable verbose logging. Specify twice to increase verbosity
      --tor-mode=[always|auto|native|never] When to use Tor. "native" - .onion addresses only. "auto" - see above for details. (default: auto)
      --tor=                                "host:port" to Tor's SOCKS proxy (default: localhost:9050 or localhost:9150)

bc1isup:
  -T, --testnet                             Check for testnet node
  -M, --mainnet                             Check for mainnet node
  -o, --output=[json|simple|none]           Choose line format: 'json' for JSON array. 'simple' for a single "up" or "down". 'none' for no output, and only
                                            exit code (default: json)

Help Options:
  -h, --help                                Show this help message

```

### Examples:

```bash
# check all Bitcoin addresses from a new-line separated file
cat addresses.txt | bc1isup

# check if there's a mainnet node running on default port on current machine
bc1isup --mainnet localhost

# check multiple addresses for running Bitcoin nodes. Use Tor for .onion addresses only
bc1isup localhost --tor-mode=native tfvfqbkl4e53uzk2.onion:8333 example.com 192.168.1.201:18333

# check all addresses from a file for running mainnet or testnet nodes and aggregate results into one flat JSON array 
cat addresses.txt | bc1isup | jq '.[]' | jq -s
```

**Note:** `bc1isup` works great with `jq`: `sudo apt install jq`, `brew install jq`


#### Exit codes

`0` is returned when every address provided returned at least one result. `1` is returned in any other case.

```bash
$ bc1isup localhost:8333
[]

$ echo $?
1

$ bc1isup localhost:8555
[{"address":"localhost:8555","useragent":"/Satoshi:0.16.99/","protocol":70015,"lastblock":534397,"testnet":false}]

$ echo $?
0

$ bc1isup -M localhost:8333
[{"address":"localhost","error":"can't connect to peer: dial tcp 127.0.0.1:8333: connect: connection refused"}]

$ echo $?
1

$ bc1isup -T --output=none localhost:8333
 
$ echo $?
1
```

### Build

1. Install Go
2. Run this:

```bash
go get -u github.com/meedamian/bc1toolkit/bc1isupcd
cd $GOPATH/src/github.com/meedamian/bc1toolkit/bc1isupcd
make install
```
