package world

import "syscall/js"

type JSObjects struct {
	Context js.Value
	Doc     js.Value
	Canvas  js.Value
}
