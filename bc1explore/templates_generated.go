// Code generated by "esc "; DO NOT EDIT.

package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// _escFS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func _escFS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// _escDir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func _escDir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// _escFSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func _escFSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		_ = f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// _escFSMustByte is the same as _escFSByte, but panics if name is not present.
func _escFSMustByte(useLocal bool, name string) []byte {
	b, err := _escFSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// _escFSString is the string version of _escFSByte.
func _escFSString(useLocal bool, name string) (string, error) {
	b, err := _escFSByte(useLocal, name)
	return string(b), err
}

// _escFSMustString is the string version of _escFSMustByte.
func _escFSMustString(useLocal bool, name string) string {
	return string(_escFSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/templates/block.html": {
		local:   "templates/block.html",
		size:    2005,
		modtime: 1536411333,
		compressed: `
H4sIAAAAAAAC/6xVQY/bLBC9768YWTl+XvSt9lBVXkvb9tDTqlK36hmHiaElkAJ2vEX+7xXGTlgn6VpK
Lw7MPN6bGSaD9ww3QiFka60cKpf1/U3BRAtrSa19yIzeZ+UNAEDB78oPUq9/wid0VEhbEH5X3hSEiba8
dCaxrrXMtyx/N7oGt6OVxAkQN8M3r7RhaJCNW65bNOPaOiN2yKCq8z0XDhO+yFlp9jK3mdeGaOTlZxQ1
dwVx/Jyfld7fDinfRmDfF8SxGTeZk19Wo5Yv0aKWX6n0LLa4QCnArlUyVFm6dkIru0Sx+6gbdW0hv4rf
f01PGxj1AhKyJ/KYXSn5/c1OOYpG7D+R/WKwFbqxUAXmy/IFBW5w85CFSlOL34zsezIcIofiB7LYXFl5
xlgQWl4V7BN27q1AvRebqU4BH6VPoAN8SU5HjiSnozHkdELuPSq28GLIbJwUZJhC43iLs2/JpOP3s78K
v4/etCDPnU1qUTRyYpTCurw2utklnN6vAluLzx28f4Dbx3GTUHi/qmLlBsShignAUFUjrMR/sHKdYBF3
JpohIilOI8qFwy2wfCOxgx+NdWLzko+vSV6h2yMqkDyYGCqLLKaMv+AY/iDd93E/Xk92poFCuZO2mJLr
e+I64v1IEzphXMaePtzSKzK7o+rw9mDn8m3jkIXDlDFYCfg/nA+oeU9Ikd7CEG3SII2cbhalxdSVKlaU
1QjDN99To4Sqs/F9ZdRRUNoBbamQod3SMFLB6fGNP5PnTwAAAP//ik4TbtUHAAA=
`,
	},

	"/templates/overview.html": {
		local:   "templates/overview.html",
		size:    4826,
		modtime: 1536501577,
		compressed: `
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
`,
	},

	"/templates/tx.html": {
		local:   "templates/tx.html",
		size:    3518,
		modtime: 1536411333,
		compressed: `
H4sIAAAAAAAC/6RWT4/bthO976cg+MsPaA+ynEWQw0YWsJsgqAN0t8i6Tm8FZY4tpjSpkiP/qaDvXpCS
bMkreWP3YphD8j1y5s2jioLDUiggdKEVgkJaljcRFxuykMzaCTV6S+MbQgiJ0tt4ZpiybIFCK/IJkAlp
ozC9jW+ikItNPLQTWSKhiVcD/xsk2nAwwOthqjdg6v8WjciAk2QVbFOBUGNVeInm+/bYHAdVII1nf0w/
RSGmpzM8LorRbDea8rKMQuQtmLCN0wv6C7PpOVA3fwXsR62WwqyZy6sdwteGjB6kXvw16iwndEyvoJwq
4sH62TqBohDLmrosOzN+OSOpgeWEFsXogVn43ciyDBO3OnQhf+IqLyHuwmPyaXwyHYXslBikhT5OmzHV
6ClhfAXE/waZEWtm9jT+FdaZ1jIK3coXqIqfgF6cvmfxD5wTgpu/oirz13DnVwJ/A7FK8RxyteIK6JlG
JslUZTkOahd1rgQSx+OXTxWhD7OP9Gq6pxwv4XvK8WrCzwA/QPMZ4L8wGIY/yPKVIRBqGYYPr3GFLZ+M
Qm+rrzl1K7rQsu256bu4qXH6rhXPZbNBCovByug8a+0rijdJ5QrkbkKODtFaYJhagb/eXCh72plSvMQP
BMKa8GApYUe+5xbFch/UD1iQAG4BFJGpC3FQFjjtc7TfDGyecuzzF5+i06jftULykwR12D2659yAtWB/
JuMeqBOLbJLRMkPByzJU7v9cu9M4XxSKw66Hw1G8sMm2VT6G933+Vimg91YdL0XYYbDO0WXsKLzmGHMm
85bKe611wLAd9eVpqMPuxj2n7zdyKdri666Iwlw27XCAO6/5g9FcKHrl5f44JHSd4wVKb0vZqxD+Jm+U
hyfua2wD9U1pfMVb6fK8z2CopK5VnjNQeMlDvGVGCbWisXU7L9PKEKbNFwuwlsa5OofaL/5e7VdFfh93
O7vd0aQsD93Y6cKaKArT9wO4/X3VySfx9zjU7thxr3baQENfKP/+t4CsMbhtHoSspfokR9SK4D6DCa0G
9FAoVCRBdVAV4QxZgHq1kuD7SrLMAq0d4H9N4M/v1oEYLY+QLy7KjGAB7DKmOPAJXTLpkHzUmb7R0h4p
KsRuWp5TvSVf77+RL1arVjoqwiYl2WBCjscXvJ/JWUVDcPSKyNWdGWBEsTVMqGHb6r4W9+7CW8ExvXs7
Hv//Q+o/v+7ejmH9oba+0Zfnp0f/xNcoh/M1hf03AAD//ziZca++DQAA
`,
	},

	"/": {
		isDir: true,
		local: "",
	},

	"/templates": {
		isDir: true,
		local: "templates",
	},
}
