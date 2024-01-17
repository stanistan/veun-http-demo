# Docs Tree

I want to have a tree representation (in memory)
of our documentation, to be able to show it like a tree
in the html.

```go
import (
    "io/fs"
    "path/filepath"
    "sort"
    "strings"
    "sync"

    "github.com/stanistan/veun-http-demo/docs"
)
```

## Node

That seems mostly fine to represent what we need.

```go
type Node struct {
	Name     string          `json:"name"`
	Href     string          `json:"href"`
	Children map[string]Node `json:"children,omitempty"`
}
```

### Tree Construction

```go
func (n *Node) insert(path string) {
	n.insertPath(strings.Split(path, string(filepath.Separator)), 0)
}

func (n *Node) insertPath(pieces []string, i int) {
    if len(pieces[i:]) == 0 {
        return
    }

    name := pieces[i]

    if n.Children == nil {
        n.Children = map[string]Node{}
    }

    node, exists := n.Children[name]
	if !exists {
		node = Node{
			Name: name,
			Href: "/" + strings.Join(pieces[:i], "/"),
		}
	}

    node.insertPath(pieces, i + 1)
    n.Children[name] = node
}
```

### Lookup, links, etc

Go maps don't have consistent ordering so we have to sort our own
keys.

```go
func (n *Node) SortedKeys() []string {
	keys := make([]string, len(n.Children))
	i := 0
	for k := range n.Children {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	return keys
}

func (n *Node) LinkInfo() (string, string) {
    name := strings.TrimSuffix(n.Name, ".go.md")
    href := filepath.Join(n.Href, name)
    if len(n.Children) > 0 {
        return name + "/", href + "/"
    } else {
        return name + ".md", href + ".md"
    }
}
```

And our `Tree` constructor is memoized to only execute one time
for the duration server runtime.

```go
var Tree = sync.OnceValue(func() Node {
    root := Node{Name: "", Href: ""}

    for _, filename := range DocFilenames() {
        root.insert(filepath.Join("docs", filename))
    }

    return root
})
```

## Parsing the docs

We can extract our entire doc tree as strings, and have
them be memoized.

```go
var DocFilenames = sync.OnceValue(func() []string {
	var filenames []string
	if err := fs.WalkDir(docs.Docs, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !entry.IsDir() && strings.HasSuffix(path, ".go.md") {
			filenames = append(filenames, path)
		}

		return nil
	}); err != nil {
		panic(err)
	}

	return filenames
})
```
