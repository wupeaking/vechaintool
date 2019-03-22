package view

import (
	"fmt"
	"github.com/andlabs/ui"
	"reflect"
	"sync"
)

var refresh = struct {
	chs       map[string]chan string
	refreshCb map[string]func()
	selects   []reflect.SelectCase
	sync.Mutex
}{
	refreshCb: make(map[string]func()),
	selects:   make([]reflect.SelectCase, 0),
	chs:       make(map[string]chan string),
}

func RegistRefreshPage(page string, cb func()) error {
	refresh.Lock()
	defer refresh.Unlock()
	if _, ok := refresh.chs[page]; ok {
		return fmt.Errorf("can't regist same page refresh: %s", page)
	}
	ch := make(chan string, 1)
	refresh.selects = append(refresh.selects, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ch),
	})
	refresh.chs[page] = ch
	refresh.refreshCb[page] = cb
	return nil
}

func StartRefresh() {
	for {
		_, recv, _ := reflect.Select(refresh.selects)
		page := recv.String()
		tkFunc := refresh.refreshCb[page]
		ui.QueueMain(tkFunc)
	}
}

func RefreshfPage(page string) {
	ch, ok := refresh.chs[page]
	if !ok {
		return
	}
	go func() {
		ch <- page
	}()
}
