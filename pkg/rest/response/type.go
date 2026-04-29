package response

type response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var Default = &response{}

func (r *response) SetCode(code int) {
	r.Code = code
}

func (r *response) SetMsg(msg string) {
	r.Msg = msg
}

func (r *response) SetData(data interface{}) {
	r.Data = data
}

func (r *response) Clone() response {
	return *r
}

type pageData struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
}
