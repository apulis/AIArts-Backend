// +build linux

package utils

import (
	"syscall"
)

func Umask(mask int) (oldmask int)  {
	return syscall.Umask(mask)
}
