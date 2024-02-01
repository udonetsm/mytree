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

// exchange full path for spaces and separator
func buildBranch(path string) (*bytes.Buffer, string) {
	buffer := new(bytes.Buffer)
	piecesOfPath := strings.Split(path, DELIM)
	for i := 0; i < len(piecesOfPath)-1; i++ {
		buffer.WriteString("|  ")
	}
	return buffer, path
}

// the function shows new branch
func outputBranch(separator *bytes.Buffer, name string, isDir bool) {
	separator.WriteString("|___" + name)
	var color ct.Color
	if isDir {
		color = ct.Blue
	} else {
		color = ct.None
	}
	ct.Foreground(color, isDir)
	fmt.Println(separator.String())
	ct.ResetColor()
}

func tree(path string) {
	fs, _ := os.ReadDir(path)
	// get spaces instead of full path
	separator, name := buildBranch(path)
	if name == "." {
	} else {
		outputBranch(separator, name, true)
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
			separator, name := buildBranch(filepath.Join(path, v.Name()))
			defer outputBranch(separator, name, false)
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
