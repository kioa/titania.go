package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// オプション解析
func OptParse() ([]string, []string) {
	// テストする言語を指定する
	// ここは flag を使わずに自前処理でも良さそう
	var flagLanguages string
	flag.StringVar(
		&flagLanguages,
		"lang", "", "languages to test (ex. --lang=ruby,python3,java)")

	flag.Parse()
	args := flag.Args()

	var directories []string

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(args) == 0 {
		entries, err := ioutil.ReadDir(pwd)
		if err != nil {
			panic(err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				directories = append(directories, entry.Name())
			}
		}

	} else {
		for _, directory := range args {
			// パスを相対パスとして綺麗な形に
			dirname, err := filepath.Rel(pwd, filepath.Join(pwd, directory))
			if err != nil {
				panic(err)
			}
			directories = append(directories, dirname)
		}
	}

	var languages []string

	if flagLanguages != "" {
		languages = strings.Split(flagLanguages, ",")
	}

	return directories, languages
}