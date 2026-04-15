package rum

import "context"

// IRegister is a callable agent function
type IRegister[in, out any] struct {
	Fn func(ctx context.Context, req in) (out, error)
}

// IDispatchResult holds the result of a completed dispatch call
type IDispatchResult struct {
	IsReady bool
	Metric  ProfileMetric
	Result  []byte
}

func NewDispatchResult() *IDispatchResult {
	return &IDispatchResult{
		IsReady: false,
		Metric:  NewProfileMetric(),
	}
}
