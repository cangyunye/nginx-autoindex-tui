# nginx-autoindex-tui

一个基于 Go + Bubble Tea + Lip Gloss 的终端文件浏览与下载工具，用于访问启用了 nginx autoindex 的静态文件服务器。

## 术语

**根路径（Root URL）：**
程序启动时传入的第一个位置参数，即用户最初传入的那个 nginx autoindex URL 所指向的目录。所有后续的目录浏览操作都不能脱离这个范围，即使链接中包含 `../` 也不会真正脱离根路径。
_Avoid：_ 起始 URL、初始路径、根目录

**历史栈（History Stack）：**
记录用户在浏览过程中进入过的目录层级的线性结构，每次进入一个新的子目录就向栈顶追加，按返回上级键（索引输入为空时的 `Backspace` / `Esc` / `b`）时弹出栈顶并回退到上一层。
_Avoid：_ 浏览历史、访问栈、导航栈

**索引导航（Index Navigation）：**
通过数字键输入最多三位的索引编号来快速定位到表格中对应行的交互方式，输入过程中光标实时移动到对应位置；上下方向键与数字输入保持同步——当前索引值会随着光标移动而增减；`Esc` 清空索引输入并保持在当前位置；`Backspace` 在索引输入非空时删除最后一位数字。
_Avoid：_ 数字选择、行号跳转、光标定位

**标记模式（Mark Mode）：**
通过 `v` 键进入的一种交互状态，用于在文件列表中标记多个文件以便后续批量下载。标记模式下按 `Enter` / `Space` 切换当前行的标记状态；按 `x` 开始下载所有已标记的文件并退出标记模式；按 `v` / `Esc` / `q` 退出标记模式但保留已标记；按 `c` 清空所有标记。
_Avoid：_ 选择模式、多选模式、批量模式

**下载队列（Download Queue）：**
在标记模式下批量标记后，按 `x` 开始下载的文件集合。队列中的文件会被并发执行 wget 命令，同时下载数量受并发数限制，超出部分排队等待。
_Avoid：_ 任务队列、下载列表、批量下载

**并发数（Concurrency Limit）：**
同一时间允许同时运行的 wget 下载进程的最大数量。默认值为 3，可通过命令行 flag 调整。
_Avoid：_ 并行数、线程数、同时下载数

**消息框（Status Bar）：**
TUI 底部常驻的一行区域，用于显示错误信息、下载进度提示、覆盖确认提示、当前索引输入状态等交互消息。
_Avoid：_ 状态栏、底部提示、消息区域

**覆盖确认（Overwrite Confirmation）：**
当目标文件已存在于下载目录时，在开始下载前提示用户确认的交互流程。批量下载时会先汇总统计已存在的文件数量，然后一次性给出覆盖 / 跳过 / 逐个确认 的选择。
_Avoid：_ 文件存在确认、冲突处理

**强制覆盖（Force Overwrite）：**
一种全局设置，开启后所有下载直接覆盖已存在的同名文件，不再提示。可通过命令行 flag `--force` 设置初始值，也可在运行时通过按键切换。
_Avoid：_ 覆盖模式、强制模式、覆盖开关

**HTML 格式（HTML Format）：**
nginx autoindex 以 HTML 输出的目录列表，格式为 `<h1>Index of ...</h1>` + `<pre>` 标签内的多行 `<a href="...">...</a>` 链接。
_Avoid：_ 网页格式、浏览器格式

**JSON 格式（JSON Format）：**
nginx autoindex 以 `application/json` 输出的目录列表，格式为 JSON 数组，每个元素包含 `name` / `type` / `mtime` / `size` 字段。
_Avoid：_ API 格式、数据接口格式

**目录（Directory / Dir）：**
在 nginx autoindex 输出中，`href` 以 `/` 结尾的条目，或 JSON 中 `type` 为 `directory` 的条目。进入目录后会发起新的 HTTP 请求并刷新列表。
_Avoid：_ 文件夹、路径

**文件（File）：**
在 nginx autoindex 输出中，`href` 不以 `/` 结尾的条目，或 JSON 中 `type` 为 `file` 的条目。选择文件后会发起下载。
_Avoid：_ 文件条目、文件项

**URL 编码（URL Encoding）：**
访问路径时保持 URL 的百分号编码（如 `%E6%96%87%E4%BB%B6.txt`），用于构造 HTTP 请求；显示给用户的名称则解码为原始字符（如 `文件.txt`）。
_Avoid：_ 百分号编码、URL 转义

**wget 命令（wget Command）：**
程序不直接实现文件写入，而是通过 `os/exec` 调用系统中的 `wget` 命令执行下载。每次下载会在消息框或下载历史中显示实际执行的命令字符串。
_Avoid：_ 下载命令、外部下载

**主题（Theme）：**
TUI 的颜色方案，支持 `dark`（深色）、`light`（浅色）、`mono`（黑白无彩色）三种预设。默认值为 `dark`。
_Avoid：_ 配色、颜色方案、样式

**边框样式（Border Style）：**
表格边框的渲染风格，支持 `normal`（标准 Unicode 边框）、`rounded`（圆角边框）、`ascii`（纯 ASCII 边框）。默认值为 `normal`。
_Avoid：_ 边框类型、边框风格

**截断显示（Column Truncation）：**
当文件名过长超出表格列宽时，超出部分以省略号省略的显示模式。默认启用，按 `Tab` 键可切换到完全显示模式（自动换行）。
_Avoid：_ 省略、截断

**完全显示（Full Display）：**
按 `Tab` 键从截断显示切换后的模式，文件名在列宽范围内自动换行显示全部内容。再次按 `Tab` 恢复截断。
_Avoid：_ 展开显示、完整显示、换行模式

**目录缓存（Directory Cache）：**
对最近浏览过的目录的解析结果进行内存缓存，当用户在相同目录间频繁往返时避免重复发起 HTTP 请求。采用 LRU 策略，最多缓存 5 个目录。
_Avoid：_ 页面缓存、结果缓存、浏览缓存

**下载历史面板（Download History Panel）：**
按 `l` 键打开的一个独立视图，显示本次程序运行期间所有下载任务的状态与结果。按 `q` 或 `Esc` 关闭面板返回列表视图。
_Avoid：_ 下载日志、下载记录、历史视图

**Insecure 模式（Insecure Mode）：**
通过 `--insecure` / `-k` flag 启用的模式，跳过 HTTPS 的 SSL 证书验证，用于测试自签名证书的服务器。默认关闭。
_Avoid：_ 忽略证书、跳过验证

**User-Agent（User Agent）：**
HTTP 请求中用于标识客户端的字符串，默认使用 Chrome 浏览器的 User-Agent，可通过 `--user-agent` flag 自定义。
_Avoid：_ UA、客户端标识

**超时（Timeout）：**
HTTP 请求的最大等待时间，固定为 30 秒。
_Avoid：_ 响应超时、请求超时

**下载目录（Output Directory）：**
下载文件保存的本地目录，默认为当前工作目录（`pwd`），可通过 `--output-dir` flag 指定。程序启动时会自动创建指定目录。
_Avoid：_ 保存目录、目标目录、输出路径
