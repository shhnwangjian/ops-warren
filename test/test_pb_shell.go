package main

import (
	"fmt"

	"github.com/shhnwangjian/ops-warren/playbook"
)

var s = `
- name: Execute the command in remote shell; stdout goes to the specified file on the remote.
  shell: ls -lrt
  args:
    group: staff

- name: Change the working directory to somedir/ before executing the command.
  shell: passwd
  args:
    chdir: somedir/
    group: staff
`

func main() {
	f := playbook.ShellInfo{}
	f.DoAll(s)
	fmt.Println(f.Msg)
}
