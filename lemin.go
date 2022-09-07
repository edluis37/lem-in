package leminmod

import (
	"bufio"
	"fmt"
	"io"
	"leminmod/anthive"
	"os"
	"strings"
)

// RunProgramWithFile - path is filepath,
// writes result to output. Close program if has error.
func RunProgramWithFile(path string, showContent bool) {
	err := WriteResultByFilePath(os.Stdout, path, showContent)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err.Error())
		os.Exit(1)
	}
}

// WriteResultByFilePath - path is filepath.
func WriteResultByFilePath(w io.Writer, path string, writeContent bool) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("WriteResultByFilePath: %w", err)
	} else if fInfo, _ := file.Stat(); fInfo.IsDir() {
		err = fmt.Errorf("%v is directory", fInfo.Name())
		return fmt.Errorf("WriteResultByFilePath: %w", err)
	}

	scanner := bufio.NewScanner(file)
	result, err := GetResult(scanner)
	if err != nil {
		return fmt.Errorf("WriteResultByFilePath: %w", err)
	}
	if writeContent {
		file.Seek(0, io.SeekStart)
		_, err = io.Copy(w, file)
		if err != nil {
			return fmt.Errorf("WriteResultByFilePath: %w", err)
		}
		fmt.Fprint(w, "\n\n# result\n")
	}
	result.WriteResult(w)

	return nil
}

// WriteResultByContent - using for Web,
//inputs writer for write result, writes nothing if returns error
func WriteResultByContent(w io.Writer, content string, writeContent bool) error {
	scanner := bufio.NewScanner(strings.NewReader(content))
	result, err := GetResult(scanner)
	if err != nil {
		return fmt.Errorf("WriteResultByContent: %w", err)
	}

	if writeContent {
		fmt.Fprintf(w, "%v\n\n", content)
	}
	result.WriteResult(w)

	return nil
}

func errInvalidDataFormat(err error) error {
	return fmt.Errorf("invalid data format, %s", err)
}

// GetResult - returns result,
//nil if shortest disjoint paths was found
func GetResult(scanner *bufio.Scanner) (*anthive.Result, error) {
	terrain := anthive.Createanthive()
	var err error
	for scanner.Scan() {
		err = terrain.ReadDataFromLine(scanner.Text())
		if err != nil {
			return nil, errInvalidDataFormat(err)
		}
	}
	err = terrain.ValidateByFieldInfo()
	if err != nil {
		return nil, errInvalidDataFormat(err)
	}
	err = terrain.Match()
	if err != nil {
		return nil, errPaths(err)
	}
	return terrain.Result, nil
}

func errPaths(err error) error {
	return fmt.Errorf("path error, %s", err)
}
