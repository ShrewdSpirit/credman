package methods

type MethodResult map[string]interface{}

type MethodInterface interface {
	Do() (MethodResult, error)
}
