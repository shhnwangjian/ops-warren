package main

import "os"

func main() {
	filenames := []string{"/tmp/demo1.txt", "/tmp/demo2.txt"}
	_ = runForDemo1(filenames)
	_ = runForDemo2(filenames)
}

/*
这段代码其实很危险，很可能会用尽所有文件描述符。
因为 defer 语句不到函数的最后一刻是不会执行的，也就是说文件始终得不到关闭。
所以切记，一定不要在 for 循环中使用 defer 语句。
 */
func runForDemo1(filenames []string) error{
	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()
	}
	return nil
}

/*
优化
将循环体单独写一个函数，这样每次循环的时候都会调用关闭函数
 */
func runForDemo2(filenames []string) error {
	for _, filename := range filenames {
		if err := doFile(filename); err != nil {
			return err
		}
	}
	return nil
}

func doFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
