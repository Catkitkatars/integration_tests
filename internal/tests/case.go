package tests

type TestCase struct {
	Name    string  `json:"name"`
	Request Request `json:"request"`
	Extract Extract `json:"extract"`
	Asserts Assert  `json:"asserts"`
}

type TestCaseResult struct {
	Name          string         `json:"name"`
	Error         error          `json:"error"`
	Response      *Response      `json:"response"`
	ExtractResult *ExtractResult `json:"extract-result"`
	AssertResult  *AssertResult  `json:"assert-result"`
}

func (t *TestCase) Do(ctx *TestContext, baseURL string) *TestCaseResult {
	var result TestCaseResult
	result.Name = t.Name
	url := baseURL + t.Request.URL
	httpResp, respTime, err := t.Request.Send(ctx, url)

	if err != nil {
		result.Error = err
		return &result
	}

	resp, err := NewResponse(httpResp, respTime)

	if err != nil {
		result.Error = err
		return &result
	}

	extractResult := t.Extract.Do(resp, t.Extract)
	ctx.SetMany(extractResult.Variables)

	result.Response = resp
	result.ExtractResult = extractResult
	result.AssertResult = t.Asserts.Check(resp)

	return &result
}
