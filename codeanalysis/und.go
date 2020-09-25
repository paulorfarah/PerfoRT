package codeanalysis

import (
	"fmt"
	"os/exec"
)

func Understand(file string){
	c := exec.Command("cmd", "/D", "und", file)
	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}
