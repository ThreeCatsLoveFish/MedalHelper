package push

import (
	"net/url"

	"github.com/ThreeCatsLoveFish/medalhelper/manager"
	"github.com/ThreeCatsLoveFish/medalhelper/util"
)

const PushPlusName string = "push_plus"

// WeChat PushPlus push
type PushPlusPush struct {
	util.Endpoint
}

// Submit data to endpoint and finish one task
func (push PushPlusPush) Submit(pd Data) error {
	// Prepare content and header
	data := url.Values{
		"token":    []string{push.Token},
		"title":    []string{pd.Title},
		"content":  []string{pd.Content},
		"template": []string{"markdown"},
	}
	// Submit info
	_, err := manager.Get(push.URL, data)
	return err
}
