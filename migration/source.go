package main

import "os"

type host interface {
	base() string
	up(file *os.File) error
	down(name string) ([]byte, error)
}
