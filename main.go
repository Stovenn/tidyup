package main

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("You should provide a path")
		return
	}

	workDir, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	filesByExtMap := make(map[string][]string)
	err = filepath.Walk(workDir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && (strings.HasPrefix(info.Name(), ".") || path != workDir) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return err
		}
		path = strings.TrimPrefix(path, workDir)
		if strings.HasPrefix(path, ".") {
			return err
		}
		parts := strings.SplitN(info.Name(), ".", 2)
		var ext string
		if len(parts) != 2 {
			ext = "other"
		} else {
			ext = parts[1]
		}
		filesByExtMap[ext] = append(filesByExtMap[ext], path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(filesByExtMap)

	for ext, paths := range filesByExtMap {
		extensionFolder := filepath.Join(workDir, ext)
		_, err := os.Stat(extensionFolder)
		if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(extensionFolder, fs.FileMode(0755))
		}
		if err == nil {
			err = moveFiles(paths, workDir, extensionFolder)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func moveFiles(paths []string, srcDir string, destDir string) error {
	for _, path := range paths {
		f, err := os.OpenFile(filepath.Join(srcDir, path), os.O_RDONLY, 0655)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, path)
		stat, err := os.Stat(destPath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		// the file exists check its size to know if it's safe to copy
		if stat != nil && stat.Size() > 0 {
			fmt.Printf("[WARNING]: '%s' already exists do you want to overwrite it ? [y/n]\n", destPath)
			sc := bufio.NewScanner(os.Stdin)
			_ = sc.Scan()
			if !strings.EqualFold(sc.Text(), "y") {
				fmt.Printf("skipping '%s'...\n", destPath)
				continue
			}
		}
		dst, err := os.OpenFile(destPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0655)
		if err != nil {
			return err
		}
		_, err = io.Copy(dst, f)
		if err != nil {
			return err
		}
	}
	return nil
}
