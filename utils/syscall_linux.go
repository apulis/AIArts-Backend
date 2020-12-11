// +build linux

package utils

func Umask(mask int) (oldmask int)  {
	return syscall.Umask(mask)
}
