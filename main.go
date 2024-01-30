package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	ct "github.com/daviddengcn/go-colortext"
	"github.com/spf13/cobra"
)

// vriables for flags.
var (
	STARTDIRFLAG, DELIM string
	FLAGALL             bool
	FILES, DIRS         int
)

func checkOS() {
	if runtime.GOOS == "windows" {
		DELIM = `\`
	} else {
		DELIM = `/`
	}
}

// split full path, call building branch and
// return buit branch and name of current directory
func buildTree(path string) (*bytes.Buffer, string) {
	s := strings.Split(path, DELIM)
	b := buildBranch(s)
	return b, filepath.Base(path)
}

// exchange full path for spaces and separator
func buildBranch(s []string) *bytes.Buffer {
	b := new(bytes.Buffer)
	for i := 0; i < len(s)-1; i++ {
		b.WriteString("|  ")
	}
	return b
}

// the function shows new branch
func outputBranch(separator *bytes.Buffer, name string, color ct.Color, brightness bool) {
	separator.WriteString("|___" + name)
	ct.Foreground(color, brightness)
	fmt.Println(separator.String())
	ct.ResetColor()
}

func tree(path string) {
	fs, _ := os.ReadDir(path)
	// get spaces instead of full path
	separator, name := buildTree(path)
	if name == "." {
	} else {
		outputBranch(separator, name, ct.Blue, true)
		// increment amount of directory
		DIRS++
	}
	for _, v := range fs {
		// File may be hidden. The checking finds out if the file
		// is hidden and hide it from output if flag --all is upset
		if strings.HasPrefix(v.Name(), ".") && !FLAGALL {
			continue
		} else if v.IsDir() {
			// if current element is a directory, should
			// call this function recursively with path to the element.
			tree(filepath.Join(path, v.Name()))
		} else {
			// if current element is a regular file, should
			// get spaces with separator instead of full path
			// and single name of directory
			separator, name := buildTree(filepath.Join(path, v.Name()))
			outputBranch(separator, name, ct.White, false)
			// increment amount of flags
			FILES++
		}
	}
}

// preparing and start tree
func init() {
	checkOS()
	root := &cobra.Command{
		Use:  "mytree",
		Long: "Defaultly the program make tree for home directory",
		Run: func(cmd *cobra.Command, args []string) {
			// switch to the start directory otherwise,
			// output will be beside the middle of terminal
			err := os.Chdir(STARTDIRFLAG)
			if err != nil {
				log.Fatal(err)
			}
			// show start directiory on the top of tree
			ct.Foreground(ct.Blue, true)
			fmt.Println(filepath.Base(STARTDIRFLAG))
			ct.ResetColor()
			tree(".")
		},
	}
	// get default directory to show its tree if flag is upset
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	// get flags or set default if it's upset
	root.Flags().StringVarP(&STARTDIRFLAG, "path", "p", home, "set path to build tree")
	root.Flags().BoolVarP(&FLAGALL, "all", "a", false, "use for see hidden dirs and files")
	root.Execute()
}

func main() {
	// final message
	fmt.Println(DIRS, "directories,", FILES, "files")
}
