package utilities

import (
	"bufio"
	"log"
	"os"
	"path"
	"path/filepath"
)

func ReadLines(filename string) ([]string, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	ex = filepath.ToSlash(ex)
	file, err := os.Open(path.Join(path.Dir(ex) + "/input/" + filename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func WriteLines(filename string, line string) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	ex = filepath.ToSlash(ex)
	f, err := os.OpenFile(path.Join(path.Dir(ex)+"/input/"+filename), os.O_RDWR|os.O_APPEND, 0660)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	_, err2 := f.WriteString(line + "\n")
	if err2 != nil {
		log.Fatal(err2)
	}
	return nil

}
