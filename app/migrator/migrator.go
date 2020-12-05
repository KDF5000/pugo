package migrator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/inconshreveable/log15"
	"github.com/kdf5000/pugo/app/helper"
	"github.com/unknwon/com"
	"gopkg.in/yaml.v2"
)

var (
	hexoBlockSeparator = []byte("---")
)

// PostHeader header of pugo post
type PostHeader struct {
	Title      string   `toml:"title" ini:"title" yaml:"title"`
	Date       string   `toml:"date" ini:"date" yaml:"date"`
	Update     string   `toml:"update_date" ini:"update_date" yaml:"update_date"`
	AuthorName string   `toml:"author" ini:"author" yaml:"author"`
	Thumb      string   `toml:"thumb" ini:"thumb" yaml:"thumb"`
	Tags       []string `toml:"tags" ini:"-" yaml:"tags"`
	Category   []string `toml:"-" ini:"-" yaml:"categories"`
	Draft      bool     `toml:"draft" ini:"draft" yaml:"draft"`
}

// Migrator migrate Hexo source into Pugo
type Migrator struct {
	srcDir  string
	destDir string
}

// NewMigrator create a migrator object to migrate
// Hexo format source file to Pugo
func NewMigrator(src, dest string) *Migrator {
	return &Migrator{
		srcDir:  src,
		destDir: dest,
	}
}

func getFirstBreakByte(data []byte) int {
	for i, v := range data {
		if v == 10 {
			return i
		}
	}
	return 0
}

// Migrate start migration
func (m *Migrator) Migrate() error {
	if !com.IsDir(m.srcDir) || !com.IsDir(m.destDir) {
		return fmt.Errorf("src or dest dir is not a directory")
	}

	return filepath.Walk(m.srcDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		data, err := m.convertFormat(p)
		if err != nil {
			// log15.Warn("Migrate|Convert|%s|%s", p, err)
			fmt.Printf("%s\n", p)
			return nil
		}

		relFile, _ := filepath.Rel(m.srcDir, p)
		dstFile := filepath.Join(m.destDir, relFile)
		if com.IsExist(dstFile) {
			hash1, _ := helper.Md5File(dstFile)
			hash2 := helper.Md5(string(data))
			if hash1 == hash2 {
				log15.Debug("Migrate|Keep|%s", dstFile)
				return nil
			}
		}

		err = com.WriteFile(dstFile, data)
		if err != nil {
			return err
		}

		log15.Debug("Migrate|Write|%s", dstFile)
		return nil
	})
}

func (m *Migrator) convertFormat(file string) ([]byte, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(fileBytes) < 3 {
		return nil, fmt.Errorf("post content is too less")
	}

	dataSlice := bytes.SplitN(fileBytes, hexoBlockSeparator, 3)
	if len(dataSlice) != 3 {
		return nil, fmt.Errorf("post need front-matter block and markdown block")
	}

	var header PostHeader
	if err := yaml.Unmarshal(bytes.Trim(dataSlice[1], "\n"), &header); err != nil {
		return nil, err
	}
	header.AuthorName = "KDF5000"
	header.Update = header.Date
	header.Draft = false

	var tomlData bytes.Buffer
	tomlData.Write([]byte("```toml\n"))
	if err := toml.NewEncoder(&tomlData).Encode(&header); err != nil {
		return nil, err
	}
	tomlData.Write([]byte("```"))
	tomlData.Write(dataSlice[2])
	return tomlData.Bytes(), nil
}
