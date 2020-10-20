package playbook

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type TemplateInfo struct {
	Msg []string
}

type TemplateBook struct {
	Name string      `yaml:"name"`
	Op   *TemplateOp `yaml:"template"`
}

type TemplateOp struct {
	Src      string `yaml:"src"`      // 文件源
	Dest     string `yaml:"dest"`     // 目标
	Owner    string `yaml:"owner"`    // 文件属主
	Group    string `yaml:"group"`    // 文件属主组
	Mode     string `yaml:"mode"`     // 文件权限
	Validate string `yaml:"validate"` // 验证命令
	Backup   bool   `yaml:"backup"`   // 是否备份
	Force    bool   `yaml:"force"`    // 是否强制替换
}

func (f *TemplateInfo) readYamlConfigList(s string) ([]*TemplateBook, error) {
	conf := make([]*TemplateBook, 0)
	err := yaml.Unmarshal([]byte(s), &conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (f *TemplateInfo) parseList(s string) ([]*TemplateBook, error) {
	return f.readYamlConfigList(s)
}

func (f *TemplateInfo) readYamlConfig(s string) (*TemplateBook, error) {
	conf := &TemplateBook{}
	err := yaml.Unmarshal([]byte(s), conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func (f *TemplateInfo) parse(s string) (*TemplateBook, error) {
	return f.readYamlConfig(s)
}

func (f *TemplateInfo) Res(r ResPlayBook) {
	f.Msg = append(f.Msg, fmt.Sprintf("执行步骤名称:%s,执行模块:%s,执行结果:%s,%s\n", r.Name, r.Model, getStatus(r.Status), r.Msg))
}
