package core

type Response interface {
	Error(code ResponseError)
	Success(model interface{})
}

type CardHandler interface {
	Get(id string, resp Response)
	Search(criteria Criteria, resp Response)
	Create(req *Request, resp Response)
	Revoke(req *Request, resp Response)
}
