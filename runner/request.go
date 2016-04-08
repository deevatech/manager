package runner

import (
	"encoding/json"
	"fmt"
	. "github.com/deevatech/manager/types"
	"github.com/parnurzeal/gorequest"
)

func (ctx Context) Request(p RunParams) (*RunResults, error) {
	req := gorequest.New()
	_, body, _ := req.Post(fmt.Sprintf("http://%s/run", ctx.ContainerHostPort)).Send(p).End()

	var output JsonResult
	if err := json.Unmarshal([]byte(body), &output); err != nil {
		return nil, err
	}

	return &RunResults{Output: output}, nil
}
