package dom

import (
	"github.com/owais/rendr/components"
	"honnef.co/go/js/dom"
)

var children map[string]dom.Node = map[string]dom.Node{}

func Render(selector string, r components.Renderer) {
	root := dom.GetWindow().Document().QuerySelector(selector)
	node := r.Render().Html()
	/*
		if old, ok := children[selector]; ok {
			root.ReplaceChild(node, old)
		} else {
			root.AppendChild(node)
		}
	*/
	root.SetInnerHTML("")
	root.AppendChild(node)
	children[selector] = node
}
