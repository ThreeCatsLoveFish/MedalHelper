package push

import (
	"MedalHelper/manager"
	"MedalHelper/util"
	"net/url"
)

const PushDeerName = "push_deer"

// PushDeer push
type PushDeerPush struct {
	util.Endpoint
}

// Submit data to endpoint and finish one task
func (push PushDeerPush) Submit(pd Data) error {
	// Prepare content and header
	data := url.Values{
		"pushkey": []string{push.Token},
		"text":    []string{pd.Title},
		"desp":    []string{pd.Content},
	}
	// Submit info
	_, err := manager.Post(push.URL, data)
	return err
}

func (push PushDeerPush) Info() util.Endpoint {
	return push.Endpoint
}