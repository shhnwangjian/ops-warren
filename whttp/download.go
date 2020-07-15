package whttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Download 下载
func Download(url, writeFile string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http.Get Error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(fmt.Sprint("download failed, response code:", resp.StatusCode))
	}
	f, err := os.Create(writeFile)
	if err != nil {
		return fmt.Errorf("os.Create Error, %s", err.Error())
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("io.Copy Error, %s", err.Error())
	}
	return nil
}

// DownloadTimeout 下载，超时控制
func DownloadTimeout(url, writeFile string, timeout int) error {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest Error: %s", err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do Error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(fmt.Sprint("download failed, response code:", resp.StatusCode))
	}
	f, err := os.Create(writeFile)
	if err != nil {
		return fmt.Errorf("os.Create Error, %s", err.Error())
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("io.Copy Error, %s", err.Error())
	}
	return nil
}

// DownloadContextTimeout 下载，超时控制
func DownloadContextTimeout(url, writeFile string, timeout int) error {
	d := time.Now().Add(time.Duration(timeout) * time.Second) // deadline max
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest Error: %s", err.Error())
	}
	req = req.WithContext(ctx)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do Error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(fmt.Sprint("download failed, response code:", resp.StatusCode))
	}
	f, err := os.Create(writeFile)
	if err != nil {
		return fmt.Errorf("os.Create Error, %s", err.Error())
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("io.Copy Error, %s", err.Error())
	}
	return nil
}
