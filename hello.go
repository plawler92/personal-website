package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// err := recreate_deploy("deploy2")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	if err := recreate_deploy("deploy2"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// for i, file := range matches {
	// 	s := strconv.Itoa(i) + ": " + file
	// 	fmt.Println(s)
	// }

	// fmt.Println(matches)

	// var files []string
	// walk_func := func(path string, info os.FileInfo, err error) error {
	// 	// if !info.IsDir() {
	// 	// 	files = append(files, path)
	// 	// }
	// 	fmt.Println(path)
	// 	return nil
	// }

	if err := generate_deploy("site", "deploy2"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("At End")
}

func generate_deploy(site_path string, deploy_path string) error {
	walk_func := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if err := os.MkdirAll(deploy_path+"/"+path, os.ModePerm); err != nil {
				return err
			}

		} else {
			if err := copy_file(path, deploy_path+"/"+path); err != nil {
				return err
			}
		}

		return nil
	}

	if err := filepath.Walk(site_path, walk_func); err != nil {
		return err
	}

	return nil
}

func copy_file(oldpath string, newpath string) error {
	src, err := os.Open(oldpath)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(newpath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func recreate_deploy(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = os.Mkdir(path, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
