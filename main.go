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
	TARGET, DELIM                 string
	ALL, PERMS, DIRONLY, FULLPATH bool
	FILES, DIRS                   int
)

func checkOS() {
	if runtime.GOOS == "windows" {
		DELIM = `\`
	} else {
		DELIM = `/`
	}
}

// exchange full path for spaces and separator
func buildBranch(path string, v os.DirEntry) *bytes.Buffer {
	buffer := new(bytes.Buffer)
	piecesOfPath := strings.Split(path, DELIM)
	for i := 0; i < len(piecesOfPath)-1; i++ {
		buffer.WriteString("|  ")
	}
	return buffer
}

// the function shows new branch
func outputBranch(separator *bytes.Buffer, name string, v os.DirEntry) {
	separator.WriteString(name)
	var color ct.Color
	bright := false
	if v == nil {
		color = ct.Blue
		bright = true
	} else {
		color = ct.None
	}
	ct.Foreground(color, bright)
	fmt.Println(separator.String())
	ct.ResetColor()
}

// Performance of fullpath flag
func showAbsPath(path string) (string, string) {
	var nfiles string
	ndirs := filepath.Base(path)
	if FULLPATH {
		nf, err := filepath.Abs(path)
		if err != nil {
			log.Fatal(err)
		}
		nd, err := filepath.Abs(path)
		if err != nil {
			log.Fatal(err)
		}
		ndirs = nd
		nfiles = nf
	}
	return nfiles, ndirs
}

func tree(path string) {
	nfiles, ndirs := showAbsPath(path)
	fs, _ := os.ReadDir(path)
	var v os.DirEntry
	// get spaces instead of full path
	separator := buildBranch(path, v)
	if filepath.Base(path) == "." {
	} else {
		outputBranch(separator, "|__ "+showPerms(path)+ndirs, v)
		// increment amount of directory
		DIRS++
	}
	for _, v = range fs {
		// File may be hidden. The checking finds out if the file
		// is hidden and hide it from output if flag --all is upset
		if strings.HasPrefix(v.Name(), ".") && !ALL {
			continue
		} else if v.IsDir() {
			// if current element is a directory, should
			// call this function recursively with path to the element.
			tree(filepath.Join(path, v.Name()))
		} else {
			if DIRONLY {
				continue
			}
			// if current element is a regular file, should
			// get spaces with separator instead of full path
			// and single name of directory
			separator = buildBranch(filepath.Join(path, v.Name()), v)
			defer outputBranch(separator, "|__ "+showPerms(v)+filepath.Join(nfiles, v.Name()), v)
			// increment amount of flags
			FILES++
		}
	}
}

func showPerms(v interface{}) (mode string) {
	if PERMS {
		switch T := v.(type) {
		case os.DirEntry:
			info, err := T.Info()
			if err != nil {
				mode = ""
				return
			}
			mode = "[" + info.Mode().String() + "] "
			return
		case string:
			file, err := os.OpenFile(T, os.O_RDONLY, 7777)
			defer file.Close()
			if err != nil {
				mode = ""
				return
			}
			info, err := file.Stat()
			if err != nil {
				mode = ""
				return
			}
			mode = "[" + info.Mode().String() + "] "
			return
		}
	}
	return ""
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
			err := os.Chdir(TARGET)
			if err != nil {
				log.Fatal(err)
			}
			// show start directiory on the top of tree
			ct.Foreground(ct.Blue, true)
			fmt.Println(filepath.Base(TARGET))
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
	root.Flags().StringVarP(&TARGET, "target", "t", home, "set target path to build tree")
	root.Flags().BoolVarP(&ALL, "all", "a", false, "use for see hidden dirs and files")
	root.Flags().BoolVarP(&PERMS, "permissions", "p", false, "use for include file mode in the output")
	root.Flags().BoolVarP(&DIRONLY, "dirs", "d", false, "use for output directories only")
	root.Flags().BoolVarP(&FULLPATH, "full-path", "f", false, "set for output full path")

	root.Execute()
}

func main() {
	// final message
	fmt.Println(DIRS, "directories,", FILES, "files")
}
