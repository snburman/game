package api

type Response struct {
	Success bool
	Status  int
	Headers map[string][]string
	Body    []byte
	Error   error
}

func (r *Response) GetHeader(key string) (string, bool) {
	if r.Headers == nil {
		return "", false
	}
	if v, ok := r.Headers[key]; ok {
		return v[0], true
	}
	return "", false
}
