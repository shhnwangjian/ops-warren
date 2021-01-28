package playbook

import (
	"gopkg.in/yaml.v3"
)

func Run(t string) (error, []string) {
	var (
		res []string
	)
	m := make([]map[string]interface{}, 0)
	err := yaml.Unmarshal([]byte(t), &m)
	if err != nil {
		return err, res
	}
	for _, line := range m {
		for k, _ := range line {
			b, err := yaml.Marshal(line)
			if err != nil {
				res = append(res, err.Error())
				continue
			}
			switch k {
			case "file":
				res = append(res, fileBook(string(b))...)
			case "shell":

			case "copy":

			case "template":

			}
		}
	}
	return nil, res
}

func fileBook(s string) []string {
	f := FileInfo{}
	f.DoOne(s)
	return f.Msg
}
