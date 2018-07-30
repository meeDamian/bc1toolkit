bc1isup
=======

This is a minimal & focused unix-style tool that allows you to quickly check the status and basic info about one or more of the Bitcoin nodes.  

### Example usage:

```bash
# check all Bitcoin addresses from a file
cat addresses.txt | bc1isup

# check if there's a mainnet node running on default port on current machine
bc1isup --mainnet localhost

# check multiple addresses for running Bitcoin nodes. Use Tor for .onion addresses only
bc1isup localhost --tor-mode=native tfvfqbkl4e53uzk2.onion:8333 example.com 192.168.1.201:18333

# check all addresses from a file for running mainnet or testnet nodes and aggregate results into one flat JSON array 
cat addresses.txt | bc1isup | jq '.[]' | jq -s
```

#### Note

`bc1isup` works great with `jq`:

```bash
sudo apt install jq
brew install jq
```

`bc1isup` returns proper exit codes:
```bash
bc1isup localhost:8333
[]

echo $?
1

bc1isup localhost:8555
[{"address":"localhost:8555","useragent":"/Satoshi:0.16.99/","protocol":70015,"lastblock":534397,"testnet":false}]

echo $?
0

bc1isup -M localhost:8333
[{"address":"localhost","error":"can't connect to peer: dial tcp 127.0.0.1:8333: connect: connection refused"}]

echo $?
1

bc1isup -T --output=none localhost:8333
 
echo $?
1
```

### Build

First, install Go, then:

```bash
go get -u github.com/meedamian/bc1toolkit/bc1isupcd
cd $GOPATH/src/github.com/meedamian/bc1toolkit/bc1isupcd
make install
```
