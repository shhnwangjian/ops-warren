package playbook

type PlayBook interface {
	HandlerAll() error
	HandlerOne(fb FileBook)
}

var PlayInfo = make(map[string]PlayBook)

func Register(name string, collect PlayBook) {
	if collect == nil {
		panic("config: Register PlayBook is nil")
	}
	if _, ok := PlayInfo[name]; ok {
		panic("config: Register PlayBook twice for adapter " + name)
	}
	PlayInfo[name] = collect
}