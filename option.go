package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// パスを相対パスとして綺麗な形に
func cleanerPath(pwd, directory string) (string, error) {
	if filepath.IsAbs(directory) {
		dirname, err := filepath.Rel(pwd, directory)
		return dirname, err
	} else {
		dirname, err := filepath.Rel(pwd, filepath.Join(pwd, directory))
		return dirname, err
	}
}

// オプション解析
func OptParse() ([]string, []string, bool) {
	// テストする言語を指定する
	// ここは flag を使わずに自前処理でも良さそう
	var flagLanguages string
	flag.StringVar(
		&flagLanguages,
		"lang", "", "languages to test (ex. --lang=ruby,python3,java)")

	// async オプション（実験的）
	var flagAsync bool
	flag.BoolVar(
		&flagAsync,
		"async", false, "execute all tests asynchronously (experimantal)")

	flag.Parse()
	args := flag.Args()

	directories := []string{}
	languages := []string{}

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
			dirname, err := cleanerPath(pwd, directory)
			if err != nil {
				panic(err)
			}
			directories = append(directories, dirname)
		}

	}

	if flagLanguages != "" {
		languages = strings.Split(flagLanguages, ",")
	}

	return directories, languages, flagAsync
}
