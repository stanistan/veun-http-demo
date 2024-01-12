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

## The node

That seems mostly fine to represent what we need.

```go
type Node struct {
	Name     string          `json:"name"`
	Href     string          `json:"href"`
	Children map[string]Node `json:"children,omitempty"`
}

func (n *Node) insert(path string) {
	n.insertPath(strings.Split(path, string(filepath.Separator)), 0)
}

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
    href := filepath.Join("/docs", n.Href, name)
    return name, href
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
			Href: strings.Join(pieces[:i], "/"),
		}
	}

    node.insertPath(pieces, i + 1)
    n.Children[name] = node
}
```

And our `Tree` constructor is memoized to only execute one time
for the duration of the server.


```go
var Tree = sync.OnceValue(func() Node {
	root := Node{Name: "docs"}

    for _, filename := range DocFilenames() {
        root.insert(filename)
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


