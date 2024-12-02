package constants

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const ECHO = "echo"
const EXIT = "exit"
const TYPE = "type"

// const CAT = "cat"

const TOKEN_DELIMITER = " "
const LINE_DELIMITER = '\n'

const BUILTIN = "a shell builtin"

var mapCommand2Path = map[string]string{
	ECHO: BUILTIN,
	EXIT: BUILTIN,
	TYPE: BUILTIN,
	// CAT:  BUILTIN,
}

func readEntry(entry fs.DirEntry, prefix string, level int) {
	curEntry := filepath.Join(prefix, entry.Name())
	// fmt.Println(strings.Repeat("| ", level), entry.Name())
	if entry.IsDir() {
		childEntries, err := os.ReadDir(curEntry)
		if err != nil {
			fmt.Println(strings.Repeat("| ", level+1), entry.Name(), err)
		}
		for _, entry := range childEntries {
			readEntry(entry, curEntry, level+1)
		}
	}
	if mapCommand2Path[entry.Name()] != BUILTIN {
		// fmt.Println(entry.Name())
		mapCommand2Path[strings.TrimSpace(entry.Name())] = curEntry
	}
}

func GetCommand(command string) (string, bool) {
	if path, isOk := mapCommand2Path[command]; isOk {
		return path, isOk
	}
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ":")
	for _, path := range paths {
		fileEntries, _ := os.ReadDir(path)
		for _, entry := range fileEntries {
			readEntry(entry, path, 0)
		}
	}
	path, isOk := mapCommand2Path[command]
	return path, isOk
}
