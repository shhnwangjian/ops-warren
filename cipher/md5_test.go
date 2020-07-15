package cipher

import (
	"fmt"
	"testing"
)

func TestFileMD5(t *testing.T) {
	fmt.Println(FileMD5(`C:\wangjian\go\project\conf\test.txt`))
}
