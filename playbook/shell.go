package playbook

type ShellInfo struct {
	Msg []string
}

type ShellBook struct {
	Name         string    `yaml:"name"`
	Shell        string    `yaml:"shell"`
	Args         ShellArgs `yaml:"args"`
	IgnoreErrors bool      `yaml:"ignore_errors"`
}

type ShellArgs struct {
	Chdir      string `yaml:"chdir"`
	Creates    string `yaml:"creates"` // A filename, when it already exists, this step will not be run.
	Removes    string `yaml:"removes"` // A filename, when it does not exist, this step will not be run.
	Executable string `yaml:"executable"`
	Warn       string `yaml:"warn"`
	Cmd        string `yaml:"cmd"`
}
