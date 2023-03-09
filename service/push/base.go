package push

import "github.com/ThreeCatsLoveFish/medalhelper/util"

var pushMap map[string]Push

// InitPush bind endpoints with config file
func InitPush() {
	for _, endpoint := range util.GlobalConfig.Endpoints {
		SetEndpoint(endpoint)
	}
}

func SetEndpoint(endpoint util.Endpoint) {
	switch endpoint.Type {
	case PushDeerName:
		addPush(endpoint.Name, PushDeerPush{endpoint})
	case PushPlusName:
		addPush(endpoint.Name, PushPlusPush{endpoint})
	case TelegramName:
		addPush(endpoint.Name, TelegramPush{endpoint})
	case BarkName:
		addPush(endpoint.Name, BarkPush{endpoint})
	}
}

// Data represents data needed for push
type Data struct {
	Title   string
	Content string
}

// Push contain all info needed for push action
type Push interface {
	Submit(data Data) error
}

func addPush(name string, push Push) {
	if pushMap == nil {
		pushMap = make(map[string]Push)
	}
	pushMap[name] = push
}

func NewPush(name string) Push {
	if push, ok := pushMap[name]; ok {
		return push
	}
	panic("push not found")
}
