package utils

import "fmt"

func Padding0(otp string, expectedLength int) string {
	fmtTemplate := fmt.Sprintf("%%0%ds", expectedLength)
	return fmt.Sprintf(fmtTemplate, otp)
}
