package logrotate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron"
)

type RotateWriter struct {
	lock     sync.Mutex
	filename string
	fp       *os.File
	cron     string
	count    int
}

func New(filename, cron string, count int) *RotateWriter {
	writer := &RotateWriter{
		filename: filename,
		cron:     cron,
		count:    count,
	}
	err := writer.rotate()
	if err != nil {
		log.Println(`RotateWriter rotate error:`, err)
		return nil
	}
	return writer
}

func (writer *RotateWriter) Write(output []byte) (int, error) {
	writer.lock.Lock()
	defer writer.lock.Unlock()
	return writer.fp.Write(output)
}

func (writer *RotateWriter) rotate() (err error) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	// Close
	if writer.fp != nil {
		err = writer.fp.Close()
		if err != nil {
			return
		}
	}

	// Rename if exists
	_, err = os.Stat(writer.filename)
	if err == nil {
		name := writer.filename + "." + time.Now().Format("2006-01-02_150405")
		err = os.Rename(writer.filename, name)
		if err != nil {
			return
		}
	}

	// Create a file
	writer.fp, err = os.Create(writer.filename)
	return
}

// 仅保留count个数的日志文件，其余删除之
func cleanOlderFiles(filename string, count int) {
	info, e := os.Stat(filename)
	if e != nil || info.IsDir() {
		return
	}

	dir := filepath.Dir(filename)
	base := filepath.Base(filename)

	// read directory
	infos, e := ioutil.ReadDir(dir)
	if e != nil {
		return
	}

	// find history logs
	prefix := fmt.Sprintf(`%s.`, base)
	historyLogs := make([]string, 0)
	for _, file := range infos {
		if strings.HasPrefix(file.Name(), prefix) {
			historyLogs = append(historyLogs, fmt.Sprintf(`%s/%s`, dir, file.Name()))
		}
	}

	if len(historyLogs) < count {
		return
	}

	// remove history logs
	sort.Strings(historyLogs)
	for i := count; i < len(historyLogs); i++ {
		os.Remove(historyLogs[i])
	}
}

// 定时任务
func (writer *RotateWriter) CronTask() {
	c := cron.New()
	c.AddFunc(writer.cron, func() {
		log.Println(`CronTask start ...`)
		writer.rotate()
		cleanOlderFiles(writer.filename, writer.count)
	})
	c.Start()
	defer c.Stop()
	select {}
}
