package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// .../path/<sourceId>/users.json -> []Config

// using this to import configs from an existing path
// walk this path of directories (source id) and extract
// login []Config from a 'users.json' file at each location
func ImportFrom(path string) ([]Config, error) {
	var cfgs []Config
	fmt.Println("walking:", path)
	walk := func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		id := filepath.Base(filepath.Dir(filename))
		base := filepath.Base(filename)
		if !strings.HasSuffix(base, ".json") {
			return nil
		}
		var tmp []Config
		// eat all these errors since we're searching for matching json
		file, _ := os.Open(filename)
		blob, _ := ioutil.ReadAll(file)
		_ = file.Close()
		_ = json.Unmarshal(blob, &tmp)
		for _, c := range tmp {
			if c.LoginURL == "" || c.Username == "" {
				continue
			}
			c.ID = fmt.Sprintf("%s:%s", id, c.Username)
			fmt.Println("importing...", c.ID)
			cfgs = append(cfgs, c)
		}
		return nil
	}
	err := filepath.Walk(path, walk)
	fmt.Printf("loaded %d configs\n", len(cfgs))
	return cfgs, err
}
