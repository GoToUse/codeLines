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
