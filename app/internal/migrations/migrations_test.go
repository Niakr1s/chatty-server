package migrations

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_checkFormat(t *testing.T) {
	tests := []struct {
		name                  string
		fileNameWithoutSuffix string
		errWant               bool
	}{
		{"valid", "0000", false},
		{"valid", "0001", false},
		{"valid", "0020", false},
		{"valid", "1000", false},
		{"valid", "1111", false},
		{"valid", "9999", false},
		{"valid", "0999", false},
		{"empty", "", true},
		{"less than 4 digits", "000", true},
		{"more than 4 digits", "00000", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.errWant, checkFormat(tt.fileNameWithoutSuffix) != nil)
		})
	}
}

func Test_checkGaps(t *testing.T) {
	tests := []struct {
		name       string
		migrations Migrations
		wantErr    bool
	}{
		{"valid", Migrations{&Migration{Num: 0}, &Migration{Num: 1}, &Migration{Num: 2}}, false},
		{"starting from 1", Migrations{&Migration{Num: 1}, &Migration{Num: 2}, &Migration{Num: 3}}, true},
		{"gap between 0 and 2", Migrations{&Migration{Num: 0}, &Migration{Num: 2}, &Migration{Num: 3}}, true},
		{"gap between 1 and 3", Migrations{&Migration{Num: 0}, &Migration{Num: 1}, &Migration{Num: 3}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkGaps(tt.migrations)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGetMigrations(t *testing.T) {
	const contents = "hello, world"

	checkAlways := func(migrations Migrations) {}
	check3Migr := func(migrations Migrations) {
		assert.Len(t, migrations, 3)
		for i, m := range migrations {
			assert.Equal(t, i, m.Num)
		}
	}

	tests := []struct {
		name         string
		prepareFiles func(dir string)
		check        func(migrations Migrations)
		wantErr      bool
	}{
		{
			"valid",
			func(dir string) {
				createFiles(t, dir, []string{"0000.sql", "0001.sql", "0002.sql"})
			},
			check3Migr,
			false,
		},
		{
			"valid with some contents",
			func(dir string) {
				createFiles(t, dir, []string{"0000.sql", "0001.sql", "0002.sql"})
				if err := ioutil.WriteFile(path.Join(dir, "0000.sql"), []byte(contents), os.ModePerm); err != nil {
					t.Fatal(err)
				}
			},
			func(migrations Migrations) {
				check3Migr(migrations)
				assert.Equal(t, contents, migrations[0].Contents)
				assert.Empty(t, migrations[1].Contents)
				assert.Empty(t, migrations[2].Contents)
			},
			false,
		},
		{
			"valid with other files",
			func(dir string) {
				createFiles(t, dir, []string{"0000.sql", "0001.sql", "0002.sql", "0002.xml", "README.md"})
			},
			check3Migr,
			false,
		},
		{
			"with gap",
			func(dir string) {
				createFiles(t, dir, []string{"0000.sql", "0002.sql", "0003.sql"})
			},
			checkAlways,
			true,
		},
		{
			"with first not 0000",
			func(dir string) {
				createFiles(t, dir, []string{"0001.sql", "0002.sql", "0003.sql"})
			},
			checkAlways,
			true,
		},
		{
			"with wrong formatted",
			func(dir string) {
				createFiles(t, dir, []string{"0000.sql", "001.sql", "0002.sql"})
			},
			checkAlways,
			true,
		},
	}
	for i, tt := range tests {
		i := i
		t.Run(tt.name, func(t *testing.T) {
			dir := strconv.Itoa(i)
			if err := os.Mkdir(dir, os.ModePerm); err != nil {
				t.Fatalf("cannot create dir %s", dir)
			}
			defer os.RemoveAll(dir)
			tt.prepareFiles(dir)

			migrations, err := GetMigrations(dir)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantErr {
				return
			}
			tt.check(migrations)
		})
	}
}

func createFiles(t *testing.T, dir string, paths []string) {
	t.Helper()
	for _, p := range paths {
		if _, err := os.Create(path.Join(dir, p)); err != nil {
			t.Fatal(err)
		}
	}
}
