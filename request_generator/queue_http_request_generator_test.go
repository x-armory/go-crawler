package request_generator

import (
	"net/http"
	"testing"
)

func TestQueueHttpRequestGenerator_GenRequest(t *testing.T) {
	generator := QueueHttpRequestGenerator{
		Requests: make(chan *Request, 1),
	}

	go func() {
		defer generator.Finish()
		for i := 1; i <= 10; i++ {
			generator.AddRequest(&Request{
				Method: "GET",
			})
		}
	}()

	for true {
		request := generator.GenRequest()
		if request == nil {
			break
		}
		println(request.(*http.Request).Method)
	}
	println("no more")
}
