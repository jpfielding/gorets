package explorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

//JSONExist ...
func JSONExist(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

// JSONStore raw file storage
func JSONStore(filename string, data interface{}) error {
	dir := path.Dir(filename)
	// TODO dont repeat this for every write
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(filename + ".tmp")
	if err != nil {
		return err
	}
	defer f.Close()
	// TODO verify key order on output raw json  (alphabetical by design????)
	row, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// formatted
	var out bytes.Buffer
	json.Indent(&out, row, "", "\t")
	_, err = out.WriteTo(f)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	err = os.Rename(f.Name(), filename)
	if err != nil {
		return err
	}
	fmt.Println("wrote:", filename)
	return nil
}

// JSONLoad raw file load
func JSONLoad(filename string, data interface{}) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}
	blob, err := ioutil.ReadAll(file)
	err = json.Unmarshal(blob, data)
	if err != nil {
		return err
	}
	return nil
}
