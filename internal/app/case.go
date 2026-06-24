package app

type TestCase struct {
	Name    string  `json:"name"`
	Request Request `json:"request"`
	Extract Extract `json:"extract"`
	Asserts Assert  `json:"asserts"`
}

type TestCaseResult struct {
	Name          string         `json:"name"`
	Success       bool           `json:"success"`
	Error         error          `json:"error"`
	Response      *Response      `json:"response"`
	ExtractResult *ExtractResult `json:"extract-result"`
	AssertResult  *AssertResult  `json:"assert-result"`
}

func (t *TestCase) Do(ctx *TestContext, baseURL string) *TestCaseResult {
	var result TestCaseResult
	result.Name = t.Name
	httpResp, respTime, err := t.Request.Send(ctx, baseURL)

	if err != nil {
		result.Error = err
		result.Success = false
		return &result
	}

	resp, err := NewResponse(httpResp, respTime)

	if err != nil {
		result.Error = err
		result.Success = false
		return &result
	}

	extractResult := t.Extract.Do(resp)

	if !extractResult.Success {
		result.Error = extractResult.Error
		result.Success = false
		return &result
	}

	ctx.SetManyVars(extractResult.Variables)

	assertResult := t.Asserts.Check(resp)

	if !assertResult.Success {
		result.Success = false
	}

	result.Response = resp
	result.ExtractResult = extractResult
	result.AssertResult = assertResult

	return &result
}
