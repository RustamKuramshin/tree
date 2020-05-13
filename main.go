package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func getChildItem(out *bytes.Buffer, dir string, prefix string, printFiles bool) error {

	// Чтение каталога
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	// Нужно ли выводить файлы. Фильтрация среза.
	if !printFiles {
		for idx := len(files) - 1; idx >= 0; idx-- {
			if !files[idx].IsDir() {
				files = append(files[:idx], files[idx+1:]...)
			}
		}
	}

	// Обход каталога
	for i := 0; i < len(files); i++ {

		isDir := files[i].IsDir()

		// Вывод в терминал
		var size string
		if !isDir {
			if files[i].Size() > 0 {
				size = " (" + strconv.FormatInt(files[i].Size(), 10) + "b)"
			} else {
				size = " (empty)"
			}
		}
		if i < len(files)-1 {
			fmt.Fprintf(out, "%v├───%v%v\n", prefix, files[i].Name(), size)
		} else {
			fmt.Fprintf(out, "%v└───%v%v\n", prefix, files[i].Name(), size)
		}

		// Уходим в рекурсию, если встречаем каталог
		if isDir {
			nextDir := dir + string(os.PathSeparator) + files[i].Name()

			// Формирование префикса из символов графики "│" и "\t"
			var nextPrefix string
			if i < len(files)-1 {
				nextPrefix = prefix + "│" + "\t"
			} else {
				nextPrefix = prefix + "\t"
			}

			err := getChildItem(out, nextDir, nextPrefix, printFiles)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func dirTree(out *bytes.Buffer, path string, printFiles bool) error {
	err := getChildItem(out, path, "", printFiles)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	out := new(bytes.Buffer)
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}

	fmt.Print(out.String())
}
