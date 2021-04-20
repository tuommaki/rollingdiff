package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tuommaki/filediff/rollingdiff"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: " + os.Args[0] + " <oldfile> <newfile>")
		os.Exit(1)
	}

	oldData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	newData, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		panic(err)
	}

	oldChunks := rollingdiff.Signatures(oldData)
	newChunks := rollingdiff.Signatures(newData)

	changes := rollingdiff.Delta(oldChunks, newChunks)

	fmt.Printf("len(changes): %d\n", len(changes))
	for i, c := range changes {
		fmt.Printf("%d: %#v\n", i, c)
	}
}

// ByteCountBinary can be used to print "human readable" sizes.
func ByteCountBinary(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
