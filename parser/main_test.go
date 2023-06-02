package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	//"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	resultLines := []string{"http://api.tech.com/item/228", "http://api.tech.com/item/976", "http://api.tech.com/item/112", "http://api.tech.com/item/521", "http://api.tech.com/item/696", "http://api.tech.com/item/914", "http://api.tech.com/item/418", "http://api.tech.com/item/940", "http://api.tech.com/item/859", "http://api.tech.com/item/360", ""}

	var err error
	cmd := exec.Command("sh", "-c", "echo '../data/test_sample_1000.txt' | ./parser")
	out, err := cmd.CombinedOutput()
	sout := string(out) // because out is []byte

	if strings.Join(resultLines, "\n") != sout {
		fmt.Println("Expected:\n", strings.Join(resultLines, "\n"))
		fmt.Println("Got:\n", sout) // so we can see the full output
		t.Errorf("%v", err)
	}
}
