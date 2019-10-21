package framework

import (
	"fmt"
	"io/ioutil"
	"strings"

	color "github.com/multiverse-os/cli/text/ansi/color"
)

func InitializeFile(path, content string) error {
	fmt.Println("  " + color.Green("CREATE") + " " + path)
	contentBytes := []byte(content)
	return ioutil.WriteFile(path, contentBytes, 0644)
}

func InitializeREADME(name string) string {
	return `## ` + name + `
This README.md is an example to provide structure for your future starship 
project. This is where you will describe to others what your project is about
and how they can use it. The starship project is designed to be very familiar 
Rails users.`
}

func InitializeModel(name string, attributes map[string]string) string {
	return `import ` + strings.ToLower(name) +
		`
import (

)

type ` + name + ` struct {
	loop through attributes
}
`

}
