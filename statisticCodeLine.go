package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

var (
	space           = 40
	fileTotal       int
	lineSum         int
	mutex           = new(sync.Mutex)
	rootPath        string
	suffixName      string
	storeFormatOut  []string
	relativePathArr []string
	excludeDirs     excludeDirArray
)

type excludeDirArray []string

func (e *excludeDirArray) String() string {
	return fmt.Sprint(*e)
}

func (e *excludeDirArray) Set(value string) error {
	*e = append(*e, value)
	return nil
}

func main() {
	flag.StringVar(&rootPath, "p", ".", "Root path.")
	flag.StringVar(&suffixName, "s", `.go`, "Suffix name of file, startswith `.`; Such as `.go`")
	flag.Var(&excludeDirs, "e", "Exclude directories.")
	flag.Parse()

	absPath, _ := filepath.Abs(rootPath)
	fmt.Println("root Path:", absPath)
	fmt.Println("Suffix name:", suffixName)
	fmt.Println("Exclude Dirs:", excludeDirs)
	title := fmt.Sprintf("%s |%s"+"|%-15s"+"|%s\n", "type", "file-name", "status", "line-nums")
	storeFormatOut = append(storeFormatOut, title)

	done := make(chan bool)
	go codeLineSum(rootPath, done)
	<-done

	rSpace := getMaximumLinesNumber(relativePathArr)
	if rSpace < space {
		rSpace = space
	}
	length, newOut := formatOutput(storeFormatOut, rSpace, suffixName)
	for _, v := range newOut {
		fmt.Println(v)
	}
	fmt.Println(strings.Repeat("-", length))
	fmt.Println("Total files: ", fileTotal)
	fmt.Println("Total lines: ", lineSum)
}

func codeLineSum(root string, done chan bool) {
	var goes int
	goDone := make(chan bool)

	isDstDir := checkDir(root)

	defer func() {
		if pan := recover(); pan != nil {
			fmt.Printf("root: %s, panic:%#v\n", root, pan)
		}

		for i := 0; i < goes; i++ {
			<-goDone
		}

		done <- true
	}()

	if !isDstDir {
		return
	}

	rootFileInfo, err := os.Lstat(root)
	checkErr(err)

	rootFile, err := os.Open(root)
	checkErr(err)
	defer func(rootFile *os.File) {
		err := rootFile.Close()
		if err != nil {
			fmt.Printf("filename: %s, panic:%#v\n", root, err)
		}
	}(rootFile)

	if rootFileInfo.IsDir() {
		fis, err := rootFile.ReadDir(0)
		checkErr(err)

		for _, fi := range fis {
			if strings.HasPrefix(fi.Name(), ".") {
				continue
			}
			goes++

			if fi.IsDir() {
				go codeLineSum(path.Join(root, fi.Name()), goDone)
			} else {
				go readFile(path.Join(root, fi.Name()), goDone)
			}
		}
	} else {
		goes = 1
		go readFile(root, goDone)
	}
}

func readFile(fileName string, done chan bool) {
	var line int

	isDstFile := strings.HasSuffix(fileName, suffixName)
	defer func() {
		if pan := recover(); pan != nil {
			fmt.Printf("filename: %s, panic:%#v\n", fileName, pan)
		}

		if isDstFile {
			addLineNum(line)
			relativePath := strings.Split(fileName, rootPath)[1]
			relativePathArr = append(relativePathArr, relativePath)
			rline := fmt.Sprintf("file |.%s"+"|%-15s"+"|line = %d\n", relativePath, "complete", line)
			storeFormatOut = append(storeFormatOut, rline)
		}

		done <- true
	}()

	if !isDstFile {
		return
	}

	fileTotal++
	file, err := os.Open(fileName)
	checkErr(err)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("filename: %s, panic:%#v\n", fileName, err)
		}
	}(file)

	reader := bufio.NewReader(file)
	for {
		_, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		if !isPrefix {
			line++
		}
	}
}

func checkDir(dirPath string) bool {
	for _, dir := range excludeDirs {
		if path.Join(rootPath, dir) == path.Join(dirPath) {
			return false
		}
	}
	return true
}

func addLineNum(num int) {
	mutex.Lock()
	defer mutex.Unlock()

	lineSum += num
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func getMaximumLinesNumber(lineArr []string) int {
	var nums []int
	for i := 0; i < len(lineArr); i++ {
		nums = append(nums, len(lineArr[i]))
	}
	maxNum := getMaxNumber(nums)

	return maxNum
}

func getMaxNumber(arr []int) int {
	maxVal := arr[0]

	for i := 1; i < len(arr); i++ {
		if maxVal < arr[i] {
			maxVal = arr[i]
		}
	}

	return maxVal
}

func formatTitle(space int) string {
	title := fmt.Sprintf("%s |%-"+strconv.Itoa(space)+"s"+"|%-15s"+"|%s\n", "type", "file-name", "status", "line-nums")
	return title
}

func formatLine(space int, relativePath string, line int) string {
	rLine := fmt.Sprintf("file |%-"+strconv.Itoa(space)+"s"+"|%-15s"+"|line = %d\n", relativePath, "complete", line)
	return rLine
}

func formatOutput(storeOutStr []string, space int, suffixName string) (int, []string) {
	newStoreOutStr := storeOutStr
	newStoreOutStr[0] = strings.TrimRight(formatTitle(space), "\n")
	maxLength := len(newStoreOutStr[0])
	newStoreOutStr = append(newStoreOutStr[:1], append([]string{strings.Repeat("-", maxLength)}, newStoreOutStr[1:]...)...)
	restStoreOutStr := newStoreOutStr[2:]

	for i := 0; i < len(restStoreOutStr); i++ {
		re, _ := regexp.Compile("./(.*?)" + suffixName)
		filePath := re.FindString(restStoreOutStr[i])
		lineNum, _ := strconv.Atoi(strings.TrimSpace(strings.Split(restStoreOutStr[i], "line =")[1]))
		restStoreOutStr[i] = strings.TrimRight(formatLine(space, filePath, lineNum), "\n")
	}

	return maxLength, newStoreOutStr
}
