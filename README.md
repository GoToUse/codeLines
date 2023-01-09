# 用法

```bash
go run main.go -p [root path] -s [suffix name] -e [exclude dirs]
```

如果 `-e` 有多个参数，多次输入 `-e [suffix name]`。 

```bash
╰─± go run main.go -h                                              
Usage of /var/folders/mq/kc_gdgd54rb08z8h4bbly5p40000gn/T/go-build3741289270/b001/exe/main:
  -e value
        Exclude directories.
  -p string
        Root path. (default ".")
  -s .
        Suffix name of file, starts with .; Such as `.go` (default ".go")
```

**Examples:**

```bash
╰─± go run main.go -p ../word-cloud/frontend -s .js -e node_modules
root Path: /Users/dapeng/Desktop/code/Git/word-cloud/frontend
Suffix name: .js
Exclude Dirs: [node_modules]

type |file-name                                |status         |line[blank]    |line[code]
------------------------------------------------------------------------------------------
file |src/App.js                               |complete       |blank = 2      |line = 23
file |build/static/js/runtime-main.ff5b8ea4.js |complete       |blank = 0      |line = 2
file |src/components/wordCloud/index.js        |complete       |blank = 9      |line = 111
file |build/static/js/main.95d45d10.chunk.js   |complete       |blank = 0      |line = 2
file |src/index.js                             |complete       |blank = 1      |line = 19
file |src/assets/errorImg.js                   |complete       |blank = 0      |line = 2
file |build/static/js/2.772ed22b.chunk.js      |complete       |blank = 0      |line = 3
------------------------------------------------------------------------------------------
Summary: total files: 7 blanks: 12 codes: 162
```
