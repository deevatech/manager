package types

type RunParams struct {
	Language string `json:"lang" binding:"required"`
	Source   string `json:"source" binding:"required"`
	Spec     string `json:"spec" binding:"required"`
}

type RunResults struct {
	Output interface{} `json:"output"`
}

type JsonResult map[string]interface{}

type TestSubmitParams struct {
	Code string `json:"code" binding:"required"`
}
