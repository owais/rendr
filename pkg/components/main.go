package components

import (
	"fmt"

	"honnef.co/go/js/dom"
)

type Handler func(dom.Event)
type Events map[string]Handler

type Renderer interface {
	Key() string
	Value() string

	Append(...Renderer) Renderer

	On(string, Handler) Renderer

	Render() Renderer

	Html() dom.Node
	Text() string
}

type Component struct {
	Events   Events
	Children []Renderer
	Body     string
	Tag      string
}

func (c Component) Key() string {
	return c.Tag
}

func (c Component) Value() string {
	return c.Body
}

func (c Component) Append(renderers ...Renderer) Renderer {
	for _, r := range renderers {
		c.Children = append(c.Children, r)
	}
	return c
}

func (c Component) Render() Renderer {
	return c
}

func (c Component) Text() string {
	tag := c.Tag
	if tag == "" {
		tag = "div"
	}

	attrs := ""
	style := ""
	body := ""
	if c.Body != "" {
		body = c.Body
	} else {
		for _, child := range c.Children {
			if child == nil {
				continue
			}
			switch child.(type) {
			case Component:
				body += child.Render().Text()
			case AttrProp:
				attrs += child.Key() + "=\"" + child.Value() + "\""
			case StyleProp:
				style += child.Text()
			}
		}
	}
	return fmt.Sprintf("<%s %s style=\"%s\">%s</%s>", tag, attrs, style, body, tag)
}

func (c Component) Html() dom.Node {
	document := dom.GetWindow().Document()

	tag := c.Tag
	if tag == "" {
		tag = "div"
	}

	node := document.CreateElement(tag)

	if c.Body != "" {
		node.SetTextContent(c.Body)
	} else {
		styles := ""
		classes := ""
		for _, child := range c.Children {
			if child == nil {
				continue
			}
			switch child.(type) {
			case ClassProp:
				classes += " " + child.Value()
			case AttrProp:
				node.SetAttribute(child.Key(), child.Value())
			case StyleProp:
				styles += child.Text()
			case Component:
				node.AppendChild(child.Render().Html())
			}

		}
		if styles != "" {
			node.SetAttribute("style", styles)
		}

		if classes != "" {
			node.SetAttribute("class", classes)
		}
	}

	for event, handler := range c.Events {
		node.AddEventListener(event, false, handler)
	}
	return node
}

func (c Component) On(event string, handler Handler) Renderer {
	c.Events[event] = handler
	return c
}

func Any(tag string, children ...Renderer) Renderer {
	//return Component{Tag: tag, Attrs: attrs, Events: Events{}, Children: children}
	return Component{
		Tag:    tag,
		Events: Events{},
	}.Append(children...)
}

func H1(children ...Renderer) Renderer {
	return Any("h1", children...)
}

func H2(children ...Renderer) Renderer {
	return Any("h2", children...)
}

func H3(children ...Renderer) Renderer {
	return Any("h3", children...)
}

func H4(children ...Renderer) Renderer {
	return Any("h4", children...)
}

func H6(children ...Renderer) Renderer {
	return Any("h6", children...)
}

func Div(children ...Renderer) Renderer {
	return Any("div", children...)
}

func Span(children ...Renderer) Renderer {
	return Any("span", children...)
}

func Section(children ...Renderer) Renderer {
	return Any("section", children...)
}

func Ul(children ...Renderer) Renderer {
	return Any("ul", children...)
}

func Li(children ...Renderer) Renderer {
	return Any("li", children...)
}

func Button(children ...Renderer) Renderer {
	return Any("button", children...)
}

func Img(children ...Renderer) Renderer {
	return Any("img", children...)
}

func Strong(children ...Renderer) Renderer {
	return Any("strong", children...)
}

func Small(children ...Renderer) Renderer {
	return Any("small", children...)
}

func Text(body string) Renderer {
	return Component{Body: body}
}

// attributes

type AttrProp struct {
	Component
}

type ClassProp struct {
	Component
}

type StyleProp struct {
	Component
}

func (s StyleProp) Text() string {
	return s.Key() + ":" + s.Value() + ";"
}

func Attr(key string, value string) Renderer {
	return AttrProp{Component{Tag: key, Body: value}}
}

func Class(class string) Renderer {
	return ClassProp{Component{Body: class}}
}

func Style(key string, value string) Renderer {
	return StyleProp{Component{Tag: key, Body: value}}
}

func Styles(values map[string]string) []Renderer {
	styles := []Renderer{}
	for key, value := range values {
		styles = append(styles, Style(key, value))
	}
	return styles
}
