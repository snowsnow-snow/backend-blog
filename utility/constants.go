package utility

import (
	"fmt"
)

var (
	Point     = "."
	Delimiter = "-"
)

// RunWithRecover 是一个通用方法，用于在协程中执行任务并处理 panic
func RunWithRecover(task func() error) error {
	errChan := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errChan <- fmt.Errorf("recovered from panic: %v", r)
			}
			close(errChan)
		}()
		// 执行传入的任务
		errChan <- task()
	}()
	// 等待协程结束并获取错误信息
	if err := <-errChan; err != nil {
		return err
	}
	return nil
}
