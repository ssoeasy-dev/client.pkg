package client

type Response[Data any, Meta any] struct {
	data      Data
	meta      Meta
	errors    []string
	path      string
	timestamp string
	code      int
}
