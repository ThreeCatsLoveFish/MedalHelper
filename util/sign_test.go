package util

import (
	"testing"
)

func TestClientSign(t *testing.T) {
	dataStr := `{"platform":"android","uuid":"cfdecf5e-4a85-4054-a6ca-f2eecffb186c","buvid":"LNEW8JPBAAVLHDVMCYGC4TGXPU3IUFSSDR56Z","seq_id":"1","room_id":"33989","parent_id":"6","area_id":"283","timestamp":"1653671556","secret_key":"axoaadsffcazxksectbbb","watch_time":"60","up_id":"63231","up_level":"40","jump_from":"30000","gu_id":"QMNPXZD0YMJDVKGZHAFOS4WT9GCAESIX16U8PNV7QLI","play_type":"0","play_url":"","s_time":"0","data_behavior_id":"","data_source_id":"","up_session":"l:one:live:record:33989:1653582728","visit_id":"QN32KHZYM8REJ0XFYSM1V7JK9VDGBXL5","watch_status":"%7B%22pk_id%22%3A0%2C%22screen_status%22%3A1%7D","click_id":"c69a5c90-39b2-4e5a-9538-124aff562f6d","session_id":"","player_type":"0","client_ts":"1653671616"}`
	if ClientSign(dataStr) != "3d1997245ebf46203ea5912ff69938160320f3ca5e528490922b9ce1af6c387932dba8ba7212173dc4fdc582beb2fcd570139e7e885b76033056747bce9c3e9e" {
		t.Fatal("Wrong client sign")
	}
}
