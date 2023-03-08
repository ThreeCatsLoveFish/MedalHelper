package push

import (
	"github.com/ThreeCatsLoveFish/medalhelper/manager"
	"github.com/ThreeCatsLoveFish/medalhelper/util"
	"net/url"
)

const BarkName = "bark"

type BarkPush struct {
	util.Endpoint
}

func (push BarkPush) Submit(pd Data) error {
	data := url.Values{
		"device_key": []string{push.Token},
		"title":      []string{pd.Title},
		"body":       []string{pd.Content},
	}

	_, err := manager.Post(push.URL, data)
	return err
}
