package cards

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/yyewolf/rwbyadv3"
)

var CardFS = rwbyadv3.GetCardFS()

func GetEmbeddableImage(cardType string, imageType string, format string) (io.Reader, error) {
	p := path.Join("cards/img/", cardType, imageType+"."+format)

	var err error
	var data []byte

	for len(data) == 0 {
		data, err = CardFS.ReadFile(p)
		if err == nil {
			break
		}
		fmt.Println(p)
		p = path.Join(path.Dir(path.Dir(p)), imageType+"."+format)
		if p == "cards/img" {
			return nil, err
		}
	}

	return bytes.NewBuffer(data), nil
}

func GetImageURI(cardType string, imageType string, format string) (string, error) {
	p := path.Join("cards/img/", cardType, imageType+"."+format)

	for {
		_, err := CardFS.Open(p)
		if err == nil {
			break
		}

		p = path.Join(path.Dir(path.Dir(p)), imageType+"."+format)
		if p == "cards/img" {
			return "", err
		}
	}

	// remove cards/img prefix
	p, _ = strings.CutPrefix(p, "cards/img/")
	// URI is /cdn/cards/p
	uri := fmt.Sprintf("/cdn/cards/%s", p)

	return uri, nil
}

func MustGetImageURI(cardType string, imageType string, format string) string {
	uri, err := GetImageURI(cardType, imageType, format)
	if err != nil {
		logrus.Fatalf("Could not find image (%s, %s, %s)", cardType, imageType, format)
	}

	return uri
}
