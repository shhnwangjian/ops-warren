package main

import (
	"fmt"

	"github.com/shhnwangjian/ops-warren/playbook"
)

var t string = `
- name: Give insecure permissions to an existing file
  file:
    path: /Users/wangjian/go/src/github.com/shhnwangjian/ops-warren/test.txt
    owner: root
    group: root
    mode: u=r,g-rw,o-x

- name: test2
  file:
    path: /Users/wangjian/go/src/github.com/shhnwangjian/ops-warren/test2.txt
    src: /Users/wangjian/go/src/github.com/shhnwangjian/ops-warren/test.txt
    state: hard

- name: test3
  file:
    path: /Users/wangjian/go/src/github.com/shhnwangjian/ops-warren/test.txt
    state: touch

- name: test4
  file:
    path: /Users/wangjian/go/src/github.com/shhnwangjian/ops-warren/test.txt
    state: absent
`

func main() {
	PlaybookFileTest()
}

func PlaybookFileTest() {
	f := playbook.FileInfo{}
	err := f.Do(t)
	fmt.Println(err)
	fmt.Println(f.Msg)
}
