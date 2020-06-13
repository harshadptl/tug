package tug

import "sync"

var once sync.Once
var ttug *Tug

func Init(n, h, p string) {
	once.Do(func() {
		ttug = NewTug(n, h, p)
	})
}

func Pause(vars ...interface{}) {
	ttug.Pause(vars)
}
