bc1explore
=======

A minimal & self contained Bitcoin Block explorer w/o dependencies. Inspired by [Dumb Block Explorer] by [@jonasschnelli].    

### Usage:

```bash
go get -u github.com/meeDamian/bc1toolkit/bc1explore
cd $GOPATH/src/github.com/meeDamian/bc1toolkit/bc1explore
bc1explore
```

Open `http://127.0.0.1:8080` in your local web browser.

```bash
$ bc1explore --help
Usage of bc1explore:
  -mainnet string
    	Point to the REST interface of your mainnet Bitcoin node (default "http://127.0.0.1:8332")
  -port int
    	What port should this blockchain explorer work on (default 8080)
  -testnet string
    	Point to the REST interface of your testnet Bitcoin node (default "http://127.0.0.1:18332")
```
### Note

Last version with no dependencies is [available here]. Note that it needs all templates to be copied together with the binary, and might contain bugs fixed in later versions.

[available here]: https://github.com/meeDamian/bc1toolkit/tree/bc1explore-nodeps/bc1explore  

### Demo

Soon™.

In the meantime you can browse [this explorer] by [@jonasschnelli]. 

### Build

1. Install Go
2. Run this:

```bash
go get -u github.com/meedamian/bc1toolkit/bc1explore
cd $GOPATH/src/github.com/meedamian/bc1toolkit/bc1explore
go build main.go
```

[Dumb Block Explorer]: https://github.com/jonasschnelli/dumb-block-explorer
[@jonasschnelli]: https://github.com/jonasschnelli
[this explorer]: https://bitcointools.jonasschnelli.ch


### TODO:

- [ ] Squash bugs that are still out there
- [ ] add footer & breadcrumbs
- [ ] create 2nd version with [packr] to get a portable binary
- [ ] add to `Makefile`
- [ ] write better run instructions
- [ ] add torrc config
- [ ] add nginx config
- [ ] make `age` human-friendly
- [ ] make all sizes more human-friendly 

[packr]: https://github.com/gobuffalo/packr
