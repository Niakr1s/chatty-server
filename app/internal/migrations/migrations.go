package migrations

import (
	"errors"
	"io/ioutil"
	"sort"
	"strings"
)

const sqlSuffix = ".sql"

// errors
var (
	ErrGapsFound           = errors.New("migrations: gaps are found")
	ErrFoundWrongFormatted = errors.New("migrations: found sql filenames with wrong formatting")
)

// GetMigrations gets CONTENTS of files from folder, sorted in descending order.
// Migration - is a file with ".sql" extension.
// Subfolders will be ignored. Files with other extensions will be ignored.
// Filenames should be formatted: "0000.sql", "0001.sql", "0002.sql"... .
// Files with other extensions will be ignored.
// Gaps or ".sql" files, formatted differently will be errored.
func GetMigrations(folder string) (Migrations, error) {
	allItemsInDir, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	migrations := Migrations{}
	for _, item := range allItemsInDir {
		if item.IsDir() {
			continue
		}
		fileName := item.Name()
		if !strings.HasSuffix(fileName, ".sql") {
			continue
		}
		// we don't want formats, other than "0000.sql" in our folder because it can be error-prone
		if err := checkFormat(removeSuffix(fileName, sqlSuffix)); err != nil {
			return nil, err
		}
		migr, err := NewMigration(folder, fileName)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migr)
	}
	sort.Sort(migrations)
	if err := checkGaps(migrations); err != nil {
		return nil, err
	}
	if err := migrations.fillContentsAll(); err != nil {
		return nil, err
	}
	return migrations, nil
}

// checkGaps checks gaps. Input migrations should be sorted.
func checkGaps(migrations Migrations) error {
	counter := 0
	for _, m := range migrations {
		if m.Num != counter {
			return ErrGapsFound
		}
		counter++
	}
	return nil
}

func checkFormat(fileNameWithoutSuffix string) error {
	if len(fileNameWithoutSuffix) != 4 {
		return ErrFoundWrongFormatted
	}
	for _, r := range fileNameWithoutSuffix {
		if r < '0' || r > '9' {
			return ErrFoundWrongFormatted
		}
	}
	return nil
}
