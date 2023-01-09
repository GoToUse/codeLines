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

## 下载为可执行文件

```bash
go install github.com/GoToUse/codeLines@latest
```

然后在本地shell，可以直接使用 `codeLines -h` 查看使用命令，以及代码量检索功能。

```bash
╰─± codeLines -p /Users/dapeng/.virtualenvs/py310/lib/python3.10/site-packages/django/contrib/admin/ -s .py -e "migrations"
# Output
root Path: /Users/dapeng/.virtualenvs/py310/lib/python3.10/site-packages/django/contrib/admin
Suffix name: .py
Exclude Dirs: [migrations]

type |file-name                               |status         |line[blank]    |line[code]
-----------------------------------------------------------------------------------------
file |models.py                               |complete       |blank = 22     |line = 168
file |__init__.py                             |complete       |blank = 3      |line = 47
file |widgets.py                              |complete       |blank = 92     |line = 495
file |apps.py                                 |complete       |blank = 8      |line = 19
file |checks.py                               |complete       |blank = 115    |line = 1235
file |forms.py                                |complete       |blank = 6      |line = 25
file |templatetags/admin_urls.py              |complete       |blank = 17     |line = 49
file |tests.py                                |complete       |blank = 38     |line = 196
file |templatetags/__init__.py                |complete       |blank = 0      |line = 0
file |actions.py                              |complete       |blank = 12     |line = 84
file |utils.py                                |complete       |blank = 77     |line = 501
file |sites.py                                |complete       |blank = 98     |line = 508
file |views/__init__.py                       |complete       |blank = 0      |line = 0
file |options.py                              |complete       |blank = 274    |line = 2232
file |templatetags/admin_modify.py            |complete       |blank = 17     |line = 133
file |views/decorators.py                     |complete       |blank = 2      |line = 17
file |templatetags/base.py                    |complete       |blank = 5      |line = 40
file |templatetags/log.py                     |complete       |blank = 12     |line = 57
file |templatetags/admin_list.py              |complete       |blank = 60     |line = 473
file |helpers.py                              |complete       |blank = 68     |line = 492
file |views/autocomplete.py                   |complete       |blank = 18     |line = 104
file |exceptions.py                           |complete       |blank = 6      |line = 7
file |filters.py                              |complete       |blank = 70     |line = 480
file |decorators.py                           |complete       |blank = 23     |line = 88
file |views/main.py                           |complete       |blank = 42     |line = 543
-----------------------------------------------------------------------------------------
Summary: total files: 25 blanks: 1085 codes: 7993

```
