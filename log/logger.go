package log

import (
	"fmt"
	"gollect/utils"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	name     string
	receiver chan string
	wg       *sync.WaitGroup
}

func NewLogger(name string, receiver chan string, wg *sync.WaitGroup) *Logger {
	return &Logger{name: name, receiver: receiver, wg: wg}
}

func (l *Logger) Add(delta int) {
	l.wg.Add(delta)
}

func (l *Logger) done() {
	l.wg.Done()
}

func (l *Logger) Wait() {
	l.wg.Wait()
}

func (l *Logger) Close() {
	close(l.Receiver())
}

func (l *Logger) Name() string {
	return l.name
}

func (l *Logger) SetName(name string) {
	l.name = name
}

func (l *Logger) Receiver() chan string {
	return l.receiver
}

func (l *Logger) SetReceiver(receiver chan string) {
	l.receiver = receiver
}

func (l *Logger) Send(msg string) {
	l.receiver <- msg
}

func (l *Logger) GetLogPath() string {
	userDir := utils.GetEnvHOME()
	if strings.EqualFold(userDir, "") {
		panic(fmt.Errorf("The environment variable (HOME) is not defined.\n"))
	}
	return utils.MakePath(userDir, ".gollect", "log")
}

func (l *Logger) OpenLogFile() (*os.File, error) {
	logPath := l.GetLogPath()
	if err := utils.Mkdir(logPath); err != nil {
		panic(err)
	}

	return utils.OpenFile(utils.MakePath(logPath, utils.AppendPrefixTimeStampShort(".log")))
}

func (l *Logger) Do() {
	defer func() {
		switch recover() {
		default:
			l.done()
		}
	}()

	log, err := l.OpenLogFile()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := log.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Fprintln(log, time.Now().Format("2006/01/02 15:04:05.06")+" "+strings.Join(os.Args, " "))

	for {
		select {
		case text, ok := <-l.Receiver():
			if ok {
				fmt.Fprintln(log, text)
			} else {
				return
			}
		}
	}
}
