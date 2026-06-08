package ssoeasy

import "fmt"

func NewError(format string, args... any) error {
	err := fmt.Errorf(format, args...)
	return fmt.Errorf("ssoeasy: %w", err)
}
