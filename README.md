# Log Rotate

Go语言编写的类似于`logrotate`功能的日志切割工具。

使用示例：

```go
package main

import (
	"log"
	"time"

	"github.com/scue/go-logrotate"
)

func main() {
	writer := logrotate.New("/tmp/x.log", "0 * * * * *", 3) // 最多保留3个日志文件
	log.SetOutput(writer)
	go writer.CronTask() // 后台定时任务

	for {
		log.Println(time.Now())
		time.Sleep(time.Second)
	}
}
```

目录结构：

```txt
/tmp
├── x.log
├── x.log.2018-06-25_201000
├── x.log.2018-06-25_201100
└── x.log.2018-06-25_201200
```

参考链接：https://stackoverflow.com/a/28797984
