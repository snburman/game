package api

type Response struct {
	Success bool
	Status  int
	Headers map[string][]string
	Body    []byte
	Error   error
}

func (r *Response) GetHeader(key string) string {
	if r.Headers == nil {
		return ""
	}
	if v, ok := r.Headers[key]; ok {
		return v[0]
	}
	return ""
}
