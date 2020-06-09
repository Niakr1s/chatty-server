package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/html/charset"
)

// errors
var (
	ErrNoContent = errors.New("no content field in json")
)

// AnecdotCommand for string "/help"
func AnecdotCommand() CommandFunc {
	return func() (string, error) {
		c := http.Client{}
		w, err := c.Get(`http://rzhunemogu.ru/RandJSON.aspx?CType=1`)
		if err != nil {
			return "", err
		}
		defer w.Body.Close()

		// getting utf8 reader
		utf8, err := charset.NewReader(w.Body, w.Header.Get("Content-Type"))
		if err != nil {
			return "", err
		}
		b, err := ioutil.ReadAll(utf8)
		if err != nil {
			return "", err
		}

		// replacing bad characters
		b = bytes.ReplaceAll(b, []byte("\r"), []byte(""))
		b = bytes.ReplaceAll(b, []byte("\n"), []byte(" "))

		res := make(map[string]string)

		// getting content
		if err := json.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
			return "", err
		}
		content, ok := res["content"]
		if !ok {
			return "", ErrNoContent
		}
		return content, nil
	}
}
