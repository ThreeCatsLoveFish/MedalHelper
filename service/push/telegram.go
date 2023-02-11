package push

import (
	"fmt"
	"net/url"

	"github.com/ThreeCatsLoveFish/medalhelper/manager"
	"github.com/ThreeCatsLoveFish/medalhelper/util"
)

const TelegramName string = "telegram"

// Telegram BOT push
type TelegramPush struct {
	util.Endpoint
}

// Submit data to endpoint and finish one task
func (push TelegramPush) Submit(pd Data) error {
	// Prepare content and header
	data := url.Values{
		"chat_id": []string{push.Token},
		"text":    []string{fmt.Sprintf("%s, %s", pd.Title, pd.Content)},
	}
	// Submit info
	_, err := manager.Get(push.URL, data)
	return err
}
