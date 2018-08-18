bc1toolkit
==========

A set of diverse, Bitcoin-related tools. All minimal, focused and kept in UNIX style. 

## Legend

Click on name to get more detailed description.

| name      | short desc                |
|----------:|:--------------------------|
| [bc1isup] | Check status of BTC nodes |

[bc1isup]: https://github.com/meeDamian/bc1toolkit/tree/master/bc1isup

## Installation

1. Go to [releases] page
2. Download bundle for your OS
3. Extract it to your favourite destination
4. Make sure your favourite destination is added to `$PATH`

[releases]: https://github.com/meeDamian/bc1toolkit/releases

## From sources

**Note:** You need to have a _recent-ish_ version of [Go] installed

```bash
git clone git@github.com:meeDamian/bc1toolkit.git
cd bc1toolkit
make install
```

[Go]: https://golang.org/

### Other useful tools

The table below contains external high-quality tools.  

| name     | short desc                                                            | 
|---------:|:----------------------------------------------------------------------|
| [btcdeb] | A set of tools used to debug or construct scripts for use in Bitcoin. |
| [xpub-converter] | Coverts between any extended key formats |

[btcdeb]: https://github.com/kallewoof/btcdeb
[xpub-converter]: https://jlopp.github.io/xpub-converter/
