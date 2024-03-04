```go
import (
	"github.com/stanistan/veun"
	"github.com/stanistan/veun/el"

	"github.com/stanistan/veun-http-demo/internal/docs"
)
```

We want to make a component that lists out all of the docs
that we've written, by their filenames so that we can attach
it to the server.

### Getting the tree

This is in the [`internal/docs`](/docs/internal/docs/tree.md) package.

## Making the view

Now that we have the files always available, we can make an index page that includes
the directory.

We pass in the current url so we can set an active node on the tree.

```go
func View(current string) veun.AsView {
    return el.Div{
        el.Class("doc-tree"),
        treeView(docs.Tree(), current),
    }
}
```

And our tree view function. This is recursive and walks the
entire tree to build out the nav.

```go
func treeView(n docs.Node, current string) el.Div {
	childPages := n.SortedKeys()
	name, href := n.LinkInfo()
	attrs := el.Attrs{}

	var elName el.Param
	if len(childPages) == 0 || name == "/" {
		elName = el.A{
			el.Attr{"href", href},
			el.Text(name),
		}
	} else {
		elName = el.Text(name)
	}

	if current == href {
		elName = el.Text(name + " ↞")
		attrs["class"] += " current"
	}

	var childContent el.Param = el.Fragment{}
	if len(childPages) > 0 {
		var children []el.Param
		for _, name := range childPages {
			children = append(children, el.Li{
				treeView(n.Children[name], current),
			})
		}
		childContent = el.Ul(children)
	}

	return el.Div{
		el.Div{attrs, elName},
		childContent,
	}
}
```
