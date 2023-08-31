package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
)

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func main() {
	var startPath string
	fmt.Print("Enter the path to start scanning from (leave blank for current directory): ")
	fmt.Scanln(&startPath)

	if startPath == "" {
		startPath = "."
	}

	fmt.Println("Scanning for node_modules directories...")

	var nodeModulesDirs []string
	var nodeModulesSizes []int64
	err := filepath.WalkDir(startPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == "node_modules" {
			nodeModulesDirs = append(nodeModulesDirs, path)
			size, err := dirSize(path)
			if err != nil {
				return err
			}
			size = size / 1048576
			path = strings.TrimPrefix(path, startPath)
			nodeModulesSizes = append(nodeModulesSizes, size)
			fmt.Printf("%s (%d MB)\n", path, size)
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	if len(nodeModulesDirs) == 0 {
		color.Red.Println("No node_modules directories found.")
		return
	}

	prompt := promptui.Select{
		Label: "Do you want to delete all found node_modules directories?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		panic(err)
	}

	if result == "No" {
		color.Yellow.Println("No directories were deleted.")
		return
	}

	for _, dir := range nodeModulesDirs {
		fmt.Printf("Deleting %s...\n", dir)
		err := os.RemoveAll(dir)
		if err != nil {
			panic(err)
		}
	}

	color.Green.Println("All node_modules directories were deleted.")
	color.Cyan.Println("Made with love by cajereta using Go!")
}
