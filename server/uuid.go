package server

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	ioutil "io/ioutil"
	"os"

	"github.com/lithammer/shortuuid"
	homedir "github.com/mitchellh/go-homedir"
)

func GenerateID() string {
	hasher := md5.New()
	hasher.Write([]byte(shortuuid.New()))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateDeviceID() string {
	deviceID := GenerateID()
	home, _ := homedir.Dir()
	deviceFile := fmt.Sprintf("%s/.device_id", home)
	if _, err := os.Stat(deviceFile); os.IsNotExist(err) {
		ioutil.WriteFile(deviceFile, []byte(deviceID), 0644)
	}
	ID, _ := ioutil.ReadFile(deviceFile)
	return string(ID)
}
