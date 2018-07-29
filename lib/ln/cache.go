package ln

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path"

	"github.com/meeDamian/bc1toolkit/lib/common"
	"github.com/pkg/errors"
)

const lnCacheFileName = "ln-pubkeys"

var cacheFilePath = path.Join(common.GetCacheDir(), lnCacheFileName)

func initFile() error {
	cacheDir := common.GetCacheDir()
	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "unable to create cache directory %s", cacheDir)
	}

	f, err := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		if os.IsExist(err) {
			return nil // all good
		}

		return errors.Wrapf(err, "unable to open file %s", cacheFilePath)
	}

	defer f.Close()

	bytes := make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		return errors.Wrap(err, "unable to get secure randomness")
	}

	salt := base64.StdEncoding.EncodeToString(bytes)
	_, err = f.WriteString(fmt.Sprintf("%s\n", salt))
	if err != nil {
		return errors.Wrapf(err, "unable to init file %s with salt %s", cacheFilePath, salt)
	}

	return nil
}

func GetPubKeyFor(host, port string) (pubKey string, err error) {
	err = initFile()
	if err != nil {
		return
	}

	//btc.DoubleSha256()

	return "", nil
}

func SavePubkey(host, port, pubKey string) (err error) {
	err = initFile()
	if err != nil {
		return
	}

	return nil
}
