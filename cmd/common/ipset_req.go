package common

type IPSetReq struct {
	Token   string   `json:"token"`
	SetName string   `json:"set_name"`
	IPList  []string `json:"ip_list"`
	Timeout uint32   `json:"timeout"`
}
