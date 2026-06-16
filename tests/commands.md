# 命令行与 API 用例

## Case 1: URL 作为参数启动

**Command:**
```
nginx-autoindex-tui http://example.com/files/
```

**Expected behavior:**
- 发送 GET `http://example.com/files/`，解析页面；
- 在表格中展示 `Idx | Type | Name | Date/Time | Size`；
- 移动光标选中 `docs/` 条目并回车 → 发起 GET `http://example.com/files/docs/` 并刷新表格；
- 移动光标选中 `file1.txt` 并回车 → 在表格下方打印：

```
$ wget -N -q "http://example.com/files/file1.txt"
```

然后执行下载，下载完成后继续保留在终端 UI 中。

## Case 2: 交互模式输入 URL

**Command:**
```
nginx-autoindex-tui
```

**Prompt/Input:**
```
Enter nginx autoindex URL: http://example.com/files/
```

**Expected behavior:** 与 Case 1 一致。

## Case 3: 错误 URL

**Command:**
```
nginx-autoindex-tui http://nonexistent.invalid/
```

**Expected output in TUI:**
```
ERROR: Get "http://nonexistent.invalid/": dial tcp: lookup nonexistent.invalid: no such host
按 R 重试，ctrl+c 退出
```

## Case 4: 返回非 2xx

**Command:**
```
nginx-autoindex-tui http://example.com/404-path-xyz
```

**Expected output in TUI:**
```
ERROR: status 404 Not Found from http://example.com/404-path-xyz
```

## Case 5: 父目录链接

**Initial URL:** `http://example.com/files/docs/`
**在表格中选中 `../` 并回车 →**
**Expected:** 发起请求 `http://example.com/files/` 并重新渲染。
