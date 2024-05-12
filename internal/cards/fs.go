package cards

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"path"

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
