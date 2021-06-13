package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type file struct {
	OriginalName string
	UpdatedName  string
}

type dir struct {
	OriginalName string
	UpdatedName  string
	Files        []file
	SubDirs      []dir
}

type commandInfo struct {
	Recurse          bool
	RemoveWhitespace bool
	DryRun           bool
	Patterns         []string
}

const (
	currentDirectory     = "."
	recurseFlag          = "-r"
	removeWhiteSpaceFlag = "-w"
	dryRunFlag           = "-n"

	removePatternFlag = "-p"
)

func main() {
	info := commandInfo{}
	args := os.Args[1:]
	if contains(args, recurseFlag) {
		info.Recurse = true
	}

	if contains(args, removeWhiteSpaceFlag) {
		info.RemoveWhitespace = true
	}

	if contains(args, removePatternFlag) {
		info.Patterns = append(info.Patterns, getPatterns(args)...)
	}

	if contains(args, dryRunFlag) {
		info.DryRun = true
	}

	dirs, err := getDirectories(currentDirectory, info.Recurse)
	if err != nil {
		log.Panic(err)
	}

	for i, d := range dirs {
		updatedDir, err := getFiles(currentDirectory, d, info.Recurse)
		if err != nil {
			log.Panic(err)
		}

		dirs[i] = updatedDir
	}

	if len(info.Patterns) > 0 {
		for i, d := range dirs {
			dirs[i] = removePatternsFromDirectory(d, info.Patterns)
		}
	}

	if info.RemoveWhitespace {
		for i, d := range dirs {
			dirs[i] = removeWhitespaceFromDirectory(d)
		}
	}

	for _, d := range dirs {
		printDir(currentDirectory, currentDirectory, d)
	}

	if info.DryRun {
		return
	}

	for _, d := range dirs {
		err := renameDirectory(currentDirectory, d)
		if err != nil {
			log.Panic(err)
		}
	}
}

func printDir(origParent, updParent string, d dir) {
	orig := fmt.Sprintf("%s/%s", origParent, d.OriginalName)
	upd := fmt.Sprintf("%s/%s", updParent, d.UpdatedName)

	log.Printf("%s -> %s\n", orig, upd)
	for _, f := range d.Files {
		origFN := fmt.Sprintf("%s/%s", orig, f.OriginalName)
		updFN := fmt.Sprintf("%s/%s", upd, f.UpdatedName)

		log.Printf("\t%s -> %s\n", origFN, updFN)
	}

	for _, sd := range d.SubDirs {
		printDir(orig, upd, sd)
	}
}

func removePatternsFromDirectory(d dir, patterns []string) dir {
	d.UpdatedName = removePatterns(d.UpdatedName, patterns)
	for i, f := range d.Files {
		d.Files[i].UpdatedName = removePatterns(f.UpdatedName, patterns)
	}
	for i, sd := range d.SubDirs {
		d.SubDirs[i] = removePatternsFromDirectory(sd, patterns)
	}
	return d
}

func removeWhitespaceFromDirectory(d dir) dir {
	d.UpdatedName = removeWhitespace(d.UpdatedName)
	for i, f := range d.Files {
		d.Files[i].UpdatedName = removeWhitespace(f.UpdatedName)
	}
	for i, sd := range d.SubDirs {
		d.SubDirs[i] = removeWhitespaceFromDirectory(sd)
	}
	return d
}

func getPatterns(args []string) []string {
	var patterns []string

	for i := 0; i < len(args); i++ {
		if args[i] != removePatternFlag {
			continue
		}

		if i == len(args)-1 {
			return patterns
		}

		patterns = append(patterns, args[i+1])
		i++
	}

	return patterns
}

func removePatterns(original string, patterns []string) string {
	updated := original
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		updated = re.ReplaceAllString(updated, "")
	}

	return updated
}

func removeWhitespace(original string) string {
	return strings.ReplaceAll(original, " ", "")
}

func getDirectories(root string, recurse bool) ([]dir, error) {
	var dirs []dir
	files, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, dir{
				OriginalName: f.Name(),
				UpdatedName:  f.Name(),
			})
		}
	}

	if recurse {
		for i, d := range dirs {
			newRoot := fmt.Sprintf("%s/%s", root, d.OriginalName)
			subDirs, err := getDirectories(newRoot, true)
			if err != nil {
				return nil, err
			}

			dirs[i].SubDirs = subDirs
		}
	}

	return dirs, nil
}

func getFiles(parent string, d dir, recurse bool) (dir, error) {
	path := fmt.Sprintf("%s/%s", parent, d.OriginalName)
	log.Println(path)
	directoryFiles, err := os.ReadDir(path)
	if err != nil {
		return d, err
	}

	for _, f := range directoryFiles {
		if f.IsDir() {
			continue
		}

		d.Files = append(d.Files, file{OriginalName: f.Name(), UpdatedName: f.Name()})
	}

	if !recurse {
		return d, nil
	}

	for i, sd := range d.SubDirs {
		usd, err := getFiles(path, sd, true)
		if err != nil {
			return d, err
		}

		d.SubDirs[i] = usd
	}

	return d, nil
}

func renameDirectory(parent string, d dir) error {
	orig := fmt.Sprintf("%s/%s", parent, d.OriginalName)
	upd := fmt.Sprintf("%s/%s", parent, d.UpdatedName)

	err := os.Rename(orig, upd)
	if err != nil {
		return err
	}

	for _, f := range d.Files {
		origFN := fmt.Sprintf("%s/%s", parent, f.OriginalName)
		updFN := fmt.Sprintf("%s/%s", parent, f.UpdatedName)

		err := os.Rename(origFN, updFN)
		if err != nil {
			return err
		}
	}

	for _, sd := range d.SubDirs {
		err := renameDirectory(upd, sd)
		if err != nil {
			return err
		}
	}

	return nil
}

func contains(list []string, value string) bool {
	for _, el := range list {
		if el == value {
			return true
		}
	}

	return false
}
