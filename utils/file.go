package utils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gososy/sorpc/log"
	"github.com/pkg/errors"
)

func FileRead(filename string) (content string, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Errorf("read file %s err %s", filename, err)
		return "", err
	}
	return string(buf), err
}

func FileRead2LineSet(filename string) (m map[string]bool, err error) {
	content, err := FileRead(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(content, "\n")
	m = make(map[string]bool)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			m[line] = true
		}
	}
	return m, nil
}

func FileWrite(filename string, content string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = out.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}

func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func IsPermission(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsPermission(err) {
			return false
		}
	}
	return true
}

func WriteAll(f *os.File, buf []byte) error {
	if len(buf) > 0 {
		n, err := f.Write(buf)
		if err != nil {
			return err
		}
		if n == len(buf) {
			return nil
		}
		for n < len(buf) {
			x, err := f.Write(buf[n:])
			if err != nil {
				return err
			}
			n += x
		}
	}
	return nil
}

func GetFileMd5(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	m := md5.New()
	_, err = io.Copy(m, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", m.Sum(nil)), nil
}

func FileReadLineByLine(fp *os.File, logic func(line string) bool) error {
	r := bufio.NewReader(fp)
	for {
		line, err := r.ReadBytes('\n')
		if len(line) > 0 {
			// 行不完整
			if line[len(line)-1] != '\n' {
				break
			}
			x := strings.TrimSpace(string(line))
			if x != "" {
				if !logic(x) {
					break
				}
			}
		}
		if err != nil {
			if err != io.EOF {
				log.Errorf("read fail %s", err)
				return err
			} else {
				break
			}
		}
	}
	return nil
}

func ProcessFileLineByLine(filePath string, routineCnt uint32, logic func(line string) error) error {
	if routineCnt == 0 {
		err := errors.New("invalid routine count 0")
		log.Error(err)
		return err
	}
	ch := make(chan string, 1000)
	doneMap := make(map[string]bool)
	var doneMapMu sync.RWMutex
	var err error
	doneFilePath := filePath + ".done"
	if FileExist(doneFilePath) {
		doneMap, err = FileRead2LineSet(doneFilePath)
		if err != nil {
			return err
		}
	}
	doneFp, err := os.OpenFile(doneFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Error("open done file err:", err)
		return err
	}
	defer doneFp.Close()
	fp, err := os.Open(filePath)
	if err != nil {
		log.Errorf("open %s err %s", filePath, err)
		return err
	}
	defer fp.Close()
	producerFinished := false
	var wg sync.WaitGroup
	for i := uint32(0); i < routineCnt; i++ {
		wg.Add(1)
		go func() {
			ticker := time.NewTicker(time.Second)
			for !producerFinished || len(ch) > 0 {
				select {
				case line := <-ch:
					doneMapMu.RLock()
					existed := doneMap[line]
					doneMapMu.RUnlock()
					if existed {
						break
					}
					err := logic(line)
					if err == nil {
						doneMapMu.Lock()
						doneMap[line] = true
						doneMapMu.Unlock()
						_, _ = doneFp.WriteString(fmt.Sprintf("%s\n", line))
					} else {
						log.Errorf("process line %s err %s", line, err)
					}
				case <-ticker.C:
				}
			}
			wg.Done()
		}()
	}
	err = FileReadLineByLine(fp, func(line string) bool {
		ch <- line
		return true
	})
	if err != nil {
		log.Errorf("read file err %s", err)
		return err
	}
	producerFinished = true
	wg.Wait()
	log.Infof("finished")
	return nil
}

func GetFileLastModTime(filePath string) (time.Time, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		return time.Unix(0, 0), err
	}
	return file.ModTime(), nil
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
