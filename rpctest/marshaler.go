package rpctest

//请求数据结构
type Req struct {
	Id string
}

func (this *Req) Marshal() ([]byte, error) {
	return []byte(this.Id), nil
}
func (this *Req) Unmarshal(data []byte) error {
	this.Id = string(data)
	return nil
}
func (this *Req) String() string {
	return "req"
}

//响应数据结构
type Rsp struct {
	Msg string
}

func (this *Rsp) Marshal() ([]byte, error) {
	return []byte(this.Msg), nil
}
func (this *Rsp) Unmarshal(data []byte) error {
	this.Msg = string(data)
	return nil
}
func (this *Rsp) String() string {
	return "rsp"
}
