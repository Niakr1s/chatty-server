package migrations

import (
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

// Migration ...
type Migration struct {
	// in format of "0000.sql"
	FullFileName string

	// Num will hold int, extracted from FileName
	// 0 for "0000.sql", 1 for "0001.sql" etc.
	Num int

	// full contents of file
	Contents string
}

// NewMigration will constructs Migration and generates num.
// Doesn't read contents from file
func NewMigration(folder, filename string) (*Migration, error) {
	num, err := fileNameToNum(filename)
	if err != nil {
		return nil, err
	}
	fullFileName := path.Join(folder, filename)
	return &Migration{FullFileName: fullFileName, Num: num}, nil
}

// FillContents reads contents from file
func (m *Migration) FillContents() error {
	contentsBytes, err := ioutil.ReadFile(m.FullFileName)
	if err != nil {
		return err
	}
	m.Contents = string(contentsBytes)
	return nil
}

func fileNameToNum(filename string) (int, error) {
	filename = removeSuffix(filename, sqlSuffix)
	num, err := strconv.ParseInt(filename, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(num), nil
}

func removeSuffix(filename string, suffix string) string {
	return strings.TrimRight(filename, suffix)
}

// Migrations is slice of Migration.
// For sort to work properly, all migrations should be properly formatted (use formatIsOk function before).
type Migrations []*Migration

// Len is sort.Interface impl.
func (m Migrations) Len() int {
	return len(m)
}

// Less is sort.Interface impl.
func (m Migrations) Less(i, j int) bool {
	return m[i].Num < m[j].Num
}

// Swap is sort.Interface impl.
func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m Migrations) fillContentsAll() error {
	for _, migr := range m {
		if err := migr.FillContents(); err != nil {
			return err
		}
	}
	return nil
}
