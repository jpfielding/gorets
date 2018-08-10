package explorer

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// GZIP the output

//JSONExist ...
func JSONExist(filename string, ifNewerThan time.Duration) bool {
	stat, err := os.Stat(filename + ".gz")
	if os.IsNotExist(err) {
		return false
	}
	return time.Since(stat.ModTime()) <= ifNewerThan
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
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// formatted
	var out bytes.Buffer
	json.Indent(&out, raw, "", "\t")
	z := gzip.NewWriter(f)
	_, err = out.WriteTo(z)
	if err != nil {
		return err
	}
	err = z.Close()
	if err != nil {
		return err
	}
	f.Close()
	err = os.Rename(f.Name(), filename+".gz")
	if err != nil {
		return err
	}
	fmt.Println("wrote:", filename)
	return nil
}

// JSONLoad raw file load
func JSONLoad(filename string, data interface{}) error {
	file, err := os.Open(filename + ".gz")
	defer file.Close()
	if err != nil {
		return err
	}
	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	blob, err := ioutil.ReadAll(gz)
	if err != nil {
		return err
	}
	err = json.Unmarshal(blob, data)
	if err != nil {
		return err
	}
	return nil
}
