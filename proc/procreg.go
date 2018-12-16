package proc

import (
	"fmt"
	"github.com/0990/gorpc"
	"sort"
)

type ProcessorBinder func(bundle ProcessorBundle, userCallback gorpc.EventCallback)

var (
	procByName = map[string]ProcessorBinder{}
)

func RegisterProcessor(procName string, f ProcessorBinder) {
	procByName[procName] = f
}

func ProcessorList() (ret []string) {
	for name := range procByName {
		ret = append(ret, name)
	}

	sort.Strings(ret)
	return
}

func BindProcessorHandler(peer gorpc.Peer, procName string, userCallback gorpc.EventCallback) {
	if proc, ok := procByName[procName]; ok {
		bundle := peer.(ProcessorBundle)
		proc(bundle, userCallback)
	} else {
		panic(fmt.Sprintf("processor not found,name:`%s`", procName))
	}
}
