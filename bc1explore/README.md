bc1explore
=======

A minimal & self contained Bitcoin Block explorer w/o dependencies. Inspired by [Dumb Block Explorer] by [@jonasschnelli].    

[Dumb Block Explorer]: https://github.com/jonasschnelli/dumb-block-explorer
[@jonasschnelli]: https://github.com/jonasschnelli

### Usage:

```bash
go get -u github.com/meeDamian/bc1toolkit/bc1explore
cd $GOPATH/src/github.com/meeDamian/bc1toolkit/bc1explore
bc1explore
```

Open `http://127.0.0.1:8080` in your local web browser.

> TODO: Add `bc1explore` to releases.

### Demo

Soon™.

In the meantime you can browse [this explorer] by [@jonasschnelli]. 

[this explorer]: https://bitcointools.jonasschnelli.ch

### Build

1. Install Go
2. Run this:

```bash
go get -u github.com/meedamian/bc1toolkit/bc1explore
cd $GOPATH/src/github.com/meedamian/bc1toolkit/bc1explore
go build main.go
```

[@jonasschnelli]: https://github.com/jonasschnelli
