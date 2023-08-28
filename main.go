package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
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

// go variables
var (
	fileTotal       int
	lineSum         int
	blankSum        int
	mutex           = new(sync.Mutex)
	storeFormatOut  []string
	relativePathArr []string
	bsPool          = sync.Pool{New: func() interface{} { return make([]byte, 0, 128*1024) }}
)

// go const variables
var (
	space       = 40 // width of `file-name` column
	statusWidth = 15 // width of `status` column
	numberWidth = 15 // width of `line or blank number`
)

// go flag variables
var (
	rootPath    string          // root path of directory
	excludeDirs excludeDirArray // Multiple -e flags or string seperated by commas or spaces are supported
	suffixName  string
)

type excludeDirArray []string

func (e *excludeDirArray) String() string {
	return fmt.Sprint(*e)
}

func (e *excludeDirArray) Set(value string) error {
	commaRegex := regexp.MustCompile(`,`)
	commaMatched := commaRegex.MatchString(value)

	spaceRegex := regexp.MustCompile(`\s`)
	spaceMatched := spaceRegex.MatchString(value)

	switch {
	case commaMatched && spaceMatched:
		return errors.New("spaces and commas cannot be included at the same time, there is only one type: `spaces` or `commas`")
	case commaMatched && !spaceMatched:
		commas := strings.Split(value, ",")
		*e = commas
		return nil
	case spaceMatched && !commaMatched:
		spaces := strings.Split(value, " ")
		*e = spaces
		return nil
	default:
		*e = append(*e, value)
		return nil
	}
}

func init() {
	flag.StringVar(&rootPath, "p", ".", "Root path.")
	flag.StringVar(&suffixName, "s", `.go`, "Suffix name of file, starts with `.`; Such as `.go`")
	flag.Var(&excludeDirs, "e", "Exclude directories, Multiple -e flags or string seperated by commas or spaces are supported.")
	flag.Parse()
}

func main() {
	rootPath, err := convertToAbsPath(rootPath)
	checkErr(err)

	fmt.Println("root Path:", rootPath)
	fmt.Println("Suffix name:", suffixName)
	fmt.Printf("Exclude Dirs: %s\n\n", excludeDirs.String())
	title := formatTitle(0)
	storeFormatOut = append(storeFormatOut, title)

	done := make(chan bool)
	go codeLineSum(rootPath, done)
	<-done

	rSpace := getMaximumLinesNumber(relativePathArr)
	if rSpace < space {
		rSpace = space
	}
	length, newOut := formatOutput(storeFormatOut, rSpace, suffixName)

	if len(newOut) == 2 {
		/* Only this two lines
		type |file-name                               |status         |line[blank]    |line[code]
		-----------------------------------------------------------------------------------------
		Do nothing.
		*/
	} else {
		for _, v := range newOut {
			fmt.Println(v)
		}

		fmt.Println(strings.Repeat("-", length))
	}
	// https://stackoverflow.com/questions/4842424/list-of-ansi-color-escape-sequences/33206814#33206814
	fmt.Printf("\033[1mSummary:\033[0m total files: \033[31m%d\033[0m blanks: \033[32m%d\033[0m codes: \033[33m%d\033[0m\n",
		fileTotal, blankSum, lineSum)
}

func convertToAbsPath(root string) (path string, err error) {
	path, err = filepath.Abs(root)
	return path, err
}

func codeLineSum(root string, done chan bool) {
	var goes int
	goDone := make(chan bool)

	isDstDir := excludeDir(root)

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
	var line, blank int
	rootPath, err := convertToAbsPath(rootPath)
	checkErr(err)

	isDstFile := strings.HasSuffix(fileName, suffixName)
	defer func() {
		if pan := recover(); pan != nil {
			fmt.Printf("filename: %s, panic:%#v\n", fileName, pan)
		}

		if isDstFile {
			addNumWithLock(line, &lineSum)
			addNumWithLock(blank, &blankSum)

			relativePath := strings.Split(fileName, rootPath)[1]
			relativePathArr = append(relativePathArr, relativePath)
			rLine := formatLine(0, relativePath, line, blank)
			storeFormatOut = append(storeFormatOut, rLine)
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

	buf := getByteSlice()
	defer putByteSlice(buf)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 1024*1024)

	for scanner.Scan() {
		lineOrg := scanner.Text()

		lineOrgLen := len(strings.TrimSpace(lineOrg))

		if lineOrgLen == 0 {
			// blank lines
			blank++
		} else {
			// codes and comments
			line++
		}
	}
}

// excludeDir return true if the dirPath is contained in `excludeDirs`, otherwise false.
func excludeDir(dirPath string) bool {
	for _, dir := range excludeDirs {
		if strings.Contains(dirPath, dir) {
			return false
		}
	}
	return true
}

// addNumWithLock Adds `num` to `*sumVar` with Lock
func addNumWithLock(num int, sumVar *int) {
	mutex.Lock()
	defer mutex.Unlock()

	*sumVar += num
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

	if len(nums) == 0 {
		return 0
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
	title := fmt.Sprintf(
		"%s |%-"+strconv.Itoa(space)+"s"+"|%-"+strconv.Itoa(statusWidth)+"s"+"|%-"+strconv.Itoa(numberWidth)+"s"+"|%s\n",
		"type", "file-name", "status", "line[blank]", "line[code]")
	return title
}

func formatLine(space int, relativePath string, line, blank int) string {
	rLine := fmt.Sprintf(
		"file |%-"+strconv.Itoa(space)+"s"+"|%-"+strconv.Itoa(statusWidth)+"s"+"|%-"+strconv.Itoa(numberWidth)+"s"+"|line = %d\n",
		relativePath, "complete", fmt.Sprintf("blank = %d", blank), line)
	return rLine
}

// formatOutput format the output contents.
func formatOutput(storeOutStr []string, space int, suffixName string) (int, []string) {
	var newStoreOutStr *[]string
	err := deepCopy(&newStoreOutStr, storeOutStr)
	checkErr(err)

	(*newStoreOutStr)[0] = strings.TrimRight(formatTitle(space), "\n")
	maxLength := len((*newStoreOutStr)[0])
	*newStoreOutStr = append((*newStoreOutStr)[:1], append([]string{strings.Repeat("-", maxLength)}, (*newStoreOutStr)[1:]...)...)
	restStoreOutStr := (*newStoreOutStr)[2:]

	for i := 0; i < len(restStoreOutStr); i++ {
		var err error
		re, err := regexp.Compile("/(.*?)\\" + suffixName)
		checkErr(err)

		content := restStoreOutStr[i]
		filePath := re.FindString(content)
		blankCompile := regexp.MustCompile(`blank = \\d+`)
		lineCompile := regexp.MustCompile(`line = \\d+`)

		lineNum, err := strconv.Atoi(trimStringSpace(strings.Split(lineCompile.FindString(content), "=")[1], false))
		checkErr(err)

		blankNum, err := strconv.Atoi(trimStringSpace(strings.Split(blankCompile.FindString(content), "=")[1], false))
		checkErr(err)

		filePath = strings.TrimPrefix(filePath, "/")
		restStoreOutStr[i] = strings.TrimRight(formatLine(space, filePath, lineNum, blankNum), "\n")
	}

	return maxLength, *newStoreOutStr
}

func getByteSlice() []byte {
	return bsPool.Get().([]byte)
}

func putByteSlice(bs []byte) {
	bsPool.Put(&bs)
}

// deepCopy
// dst can't be a non-pointer type.
func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// trimStringSpace trim all space if all is true, trim leading and trailing space otherwise.
func trimStringSpace(s string, all bool) string {
	if !all {
		return strings.TrimSpace(s)
	} else {
		return strings.ReplaceAll(s, " ", "")
	}
}
