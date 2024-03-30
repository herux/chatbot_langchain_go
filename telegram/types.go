package telebot

type LoginpinToken struct {
	Token string
}

type DataPinResp struct {
	Data LoginpinToken `json:data`
}

type CommonResp struct {
	Status  string
	Message string
}

type LoginpinResp struct {
	CommonResp
	DataPinResp
}

type ContentChatResp struct {
	Content string
}

type DataChatResp struct {
	Data ContentChatResp
}

type ChatResp struct {
	CommonResp
	Data DataChatResp
}

type ChatTextResp struct {
	CommonResp
	Data string
}
