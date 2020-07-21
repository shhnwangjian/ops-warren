package playbook

var (
	PlayInfo = make(map[string]PlayBook)
)

type PlayBook interface {
	Do(string) error
	Res(int, ResPlayBook)
}

type ResPlayBook struct {
	Name   string
	Model  string
	Msg    string
	Status int8
}

func Register(name string, collect PlayBook) {
	if collect == nil {
		panic("config: Register PlayBook is nil")
	}
	if _, ok := PlayInfo[name]; ok {
		panic("config: Register PlayBook twice for adapter " + name)
	}
	PlayInfo[name] = collect
}

func getStatus(i int8) string {
	switch i {
	case 0:
		return "success"
	case -1:
		return "fail"
	default:
		return "Unknown"
	}
}
