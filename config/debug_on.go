//go:build !nodebug

package config

import "fmt"

func Info(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}
