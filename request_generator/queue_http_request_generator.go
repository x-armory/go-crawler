package request_generator

type Request struct {
	Method  string
	Url     string
	Headers map[string][]string
	Values  map[string][]string
}

type QueueHttpRequestGenerator struct {
	Requests chan *Request
}

func (g *QueueHttpRequestGenerator) AddRequest(request *Request) {
	g.Requests <- request
}

func (g *QueueHttpRequestGenerator) Finish() {
	close(g.Requests)
}

func (g *QueueHttpRequestGenerator) GenRequest() interface{} {
	req, ok := <-g.Requests
	if !ok {
		return nil
	}
	return ParseHttpRequest(req)
}
