# transfer

主要用于将命名捕获分组的重写替换规则转换为普通数字分组重写替换规则

### 用法

```bash
$ ./transfer input_json_path output_json_path
```

### 输入 json 文件格式

```json
{
	"hosts": [
		"itisatest.qiniudn.com"
	],
	"routers": [
		{
			"comment": "x_0 / convert_format",
			"pattern": "^(?P<key_prefix>(?:/[^/]+)*/[^.]+)[.]thumb[.](?P<xxx>[1-9]\\d*)_0_c[.](?P<from_format>[^_]+)_(?P<to_format>jpeg|webp)$",
			"repl": "${key_prefix}.${from_format}?imageMogr2/format/${to_format}/quality/90/thumbnail/${xxx}x/ignore-error/1"
		}
	]
	"version": "71848b9ef9074fbf9c5cfec206f8e27b"
}
```

### 输出 json 文件格式

```json
{
	"hosts": [
		"itisatest.qiniudn.com"
	],
	"routers": [
		{
			"Pattern": "^((?:/[^/]+)*/[^.]+)[.]thumb[.]([1-9]\\d*)_0_c[.]([^_]+)_(jpeg|webp)$",
			"Repl": "${1}.${3}?imageMogr2/format/${4}/quality/90/thumbnail/${2}x/ignore-error/1",
			"Comment": "x_0 / convert_format"
		}
	]
	"version": "71848b9ef9074fbf9c5cfec206f8e27b"
}
```