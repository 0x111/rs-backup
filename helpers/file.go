package helpers

import (
	"fmt"
	"io/ioutil"
)

func ReadFileContent(fileName string) []byte {
	b, err := ioutil.ReadFile(fileName) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	return b
}
