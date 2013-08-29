package execution

import (
	"fmt"
	"reflect"
)

func (parent *scope) visit() {
	parent.action()

	if len(parent.children) == 0 {
		parent.cleanup()
	} else {
		parent.visitChild()
	}
}
func (parent *scope) allChildrenVisited() bool {
	return parent.child >= len(parent.birthOrder)
}
func (parent *scope) visitChild() {
	child := parent.children[parent.birthOrder[parent.child]]
	child.visit()
	if child.visited() {
		parent.cleanup()
	}
}
func (parent *scope) cleanup() {
	for _, reset := range parent.resets {
		reset()
	}
	parent.child++
}

func (parent *scope) adopt(child *scope) {
	if parent.hasChild(child) {
		return
	}
	name := functionName(child.action)
	parent.birthOrder = append(parent.birthOrder, name)
	parent.children[name] = child
}
func (parent *scope) hasChild(child *scope) bool {
	for _, name := range parent.birthOrder {
		if name == functionName(child.action) {
			return true
		}
	}
	return false
}

func (self *scope) visited() bool {
	return self.child >= len(self.birthOrder)
}

func (self *scope) registerReset(action func()) {
	self.resets[functionId(action)] = action
}

func newScope(name string, action func()) *scope {
	fmt.Sprintf("")

	self := scope{name: name, action: action}
	self.children = make(map[string]*scope)
	self.birthOrder = []string{}
	self.resets = make(map[uintptr]func())
	return &self
}

type scope struct {
	name       string
	action     func()
	children   map[string]*scope
	birthOrder []string
	child      int
	resets     map[uintptr]func()
}

func functionId(action func()) uintptr {
	return reflect.ValueOf(action).Pointer()
}
