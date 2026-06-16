# 一次完整的请求与响应样例

针对 URL `http://example.com/files/` 发起 GET：

## Request

```
GET /files/ HTTP/1.1
Host: example.com
User-Agent: Go-http-client/1.1
Accept-Encoding: gzip
```

## Response Body（HTML 格式）

```html
<html>
<head><title>Index of /files/</title></head>
<body bgcolor="white">
<h1>Index of /files/</h1><hr><pre>
<a href="../">../</a>                                             01-Jan-2024 10:00                   -
<a href="docs/">docs/</a>                                          15-Mar-2024 12:30                   -
<a href="file1.txt">file1.txt</a>                                  01-Jan-2024 10:00                1.2K
<a href="image.png">image.png</a>                                  02-Jan-2024 11:30                  45K
</pre><hr>
</body>
</html>
```

## Response Body（JSON 格式，nginx autoindex_format json;）

```json
[
  {
    "name": "../",
    "type": "directory",
    "mtime": "2024-01-01T10:00:00+00:00",
    "size": ""
  },
  {
    "name": "docs/",
    "type": "directory",
    "mtime": "2024-03-15T12:30:00+00:00",
    "size": ""
  },
  {
    "name": "file1.txt",
    "type": "file",
    "mtime": "2024-01-01T10:00:00+00:00",
    "size": "1229"
  },
  {
    "name": "image.png",
    "type": "file",
    "mtime": "2024-01-02T11:30:00+00:00",
    "size": "46080"
  }
]
```

## 解析后得到的 Page（HTML 格式）

```
Title: "Index of /files/"

Entries:
  1. Href="../"         Name="../"         DateTime="01-Jan-2024 10:00"  Size="-"    IsDir=true
  2. Href="docs/"       Name="docs/"       DateTime="15-Mar-2024 12:30"  Size="-"    IsDir=true
  3. Href="file1.txt"   Name="file1.txt"   DateTime="01-Jan-2024 10:00"  Size="1.2K" IsDir=false
  4. Href="image.png"   Name="image.png"   DateTime="02-Jan-2024 11:30"  Size="45K"  IsDir=false
```

## Terminal UI View（概念样式）

```
Index of /files/
URL: http://example.com/files/

┌─────┬──────┬───────────┬───────────────────┬───────┐
│ Idx │ Type │ Name      │ Date/Time         │ Size  │
├─────┼──────┼───────────┼───────────────────┼───────┤
│ 1   │ DIR  │ ../       │ 01-Jan-2024 10:00 │ -     │
│ 2   │ DIR  │ docs/     │ 15-Mar-2024 12:30 │ -     │
│ 3   │ FILE │ file1.txt │ 01-Jan-2024 10:00 │ 1.2K  │
│ 4   │ FILE │ image.png │ 02-Jan-2024 11:30 │ 45K   │
└─────┴──────┴───────────┴───────────────────┴───────┘

↑/↓ or j/k 移动 · Enter 进入目录/下载文件 · R 刷新 · q/ctrl+c 退出
```
