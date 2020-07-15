package playbook

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type FileInfo struct {
}

// https://docs.ansible.com/ansible/latest/modules/file_module.html#examples

type FileBook struct {
	Name string `yaml:"name"`
	Op   FileOp `yaml:"file"`
}

type FileOp struct {
	Src   string `yaml:"src"`
	Dest  string `yaml:"dest"`
	Owner string `yaml:"owner"`
	Group string `yaml:"group"`
	State string `yaml:"state"`
	Mode  string `yaml:"mode"`
	Path  string `yaml:"path"`
}

var t = `
- name: Give insecure permissions to an existing file
  file:
    path: /work
    owner: root
    group: root
    mode: '1777'
`

func ReadYamlConfig(s string) ([]FileBook, error) {
	conf := make([]FileBook, 0)
	err := yaml.Unmarshal([]byte(s), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func Parse() ([]FileBook, error) {
	return ReadYamlConfig(t)
}

func (f *FileInfo) HandlerAll() error {
	l, err := Parse()
	if err != err {
		return err
	}
	if len(l) == 0 {
		return errors.New("file book no data")
	}
	for _, line := range l {
		f.HandlerOne(line)
	}
	return nil
}

func (f *FileInfo) HandlerOne(fb FileBook) {
}

func init() {
	Register("file", &FileInfo{})
}
