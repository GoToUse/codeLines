# 用法

```bash
go run staticCodeLine.go -p [root path] -s [suffix name] -e [exclude dirs]
```

如果 `-e` 有多个参数，多次输入 `-e [suffix name]`。 

```bash
╰─± go run statisticCodeLine.go -p /Users/dapeng/Desktop/code/react/practice/newhand -s .ts -e node_modules -e public

# Output
root Path: /Users/dapeng/Desktop/code/react/practice/newhand
Suffix name: .ts
Exclude Dirs: [node_modules public]
type |file-name                                |status         |line-nums
-------------------------------------------------------------------------
file |./src/react-app-env.d.ts                |complete       |line = 1
file |./src/reportWebVitals.ts                |complete       |line = 15
file |./src/setupTests.ts                     |complete       |line = 5
-------------------------------------------------------------------------
Total files:  3
Total lines:  21
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

type |file-name                               |status         |line-nums
------------------------------------------------------------------------
file |/decorators.py                          |complete       |line = 111
file |/sites.py                               |complete       |line = 606
file |/tests.py                               |complete       |line = 234
file |/options.py                             |complete       |line = 2506
file |/__init__.py                            |complete       |line = 50
file |/models.py                              |complete       |line = 190
file |/utils.py                               |complete       |line = 578
file |/actions.py                             |complete       |line = 96
file |/views/__init__.py                      |complete       |line = 0
file |/templatetags/admin_urls.py             |complete       |line = 66
file |/helpers.py                             |complete       |line = 560
file |/templatetags/admin_list.py             |complete       |line = 533
file |/exceptions.py                          |complete       |line = 13
file |/templatetags/log.py                    |complete       |line = 69
file |/views/decorators.py                    |complete       |line = 19
file |/templatetags/__init__.py               |complete       |line = 0
file |/filters.py                             |complete       |line = 550
file |/checks.py                              |complete       |line = 1350
file |/apps.py                                |complete       |line = 27
file |/templatetags/base.py                   |complete       |line = 45
file |/templatetags/admin_modify.py           |complete       |line = 150
file |/views/autocomplete.py                  |complete       |line = 122
file |/widgets.py                             |complete       |line = 587
file |/forms.py                               |complete       |line = 31
file |/views/main.py                          |complete       |line = 585
------------------------------------------------------------------------
Total files:  25
Total lines:  9078
```
