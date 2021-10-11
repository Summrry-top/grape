package grape

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(len(a), len(b))
	for i < max && a[i] == b[i] {
		i++
	}
	return i
}

type nodeType uint8

const (
	static nodeType = iota // default
	root
)

type node struct {
	path     string
	fullPath string
	indices  string
	nType    nodeType
	priority uint32
	children []*node
	handlers HandlersChain
}

// 增加给定子节点的优先级，并在必要时重新排序
// prio优先级
// Increments priority of the given child and reorders if necessary
func (n *node) incrementChildPrio(pos int) int {
	cs := n.children
	cs[pos].priority++
	prio := cs[pos].priority

	// Adjust position (move to front)
	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		// Swap node positions
		//交换节点的位置
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	// Build new index char string
	if newPos != pos {
		n.indices = n.indices[:newPos] + // Unchanged prefix, might be empty
			n.indices[pos:pos+1] + // The index char we move
			n.indices[newPos:pos] + n.indices[pos+1:] // Rest without char at 'pos'
	}

	return newPos
}

// addRoute adds a node with the given handle to the path.
// Not concurrency-safe!
func (n *node) addRoute(path string, handlers HandlersChain) {
	fullPath := path
	n.priority++

	// Empty tree
	if n.path == "" && n.indices == "" {
		n.insertChild(path, fullPath, handlers)
		n.nType = root
		return
	}

walk:
	for {
		// Find the longest common prefix.
		// This also implies that the common prefix contains no ':' or '*'
		// since the existing key can't contain those chars.
		i := longestCommonPrefix(path, n.path)

		// Split edge
		if i < len(n.path) {
			child := node{
				path:     n.path[i:],
				nType:    static,
				indices:  n.indices,
				children: n.children,
				handlers: n.handlers,
				priority: n.priority - 1,
			}

			n.children = []*node{&child}
			// []byte for proper unicode char conversion, see #65
			n.indices = string([]byte{n.path[i]})
			n.path = path[:i]
			n.handlers = nil
		}

		// Make new node a child of this node
		if i < len(path) {
			path = path[i:]

			idxc := path[0]

			// Check if a child with the next path byte exists
			for i, c := range []byte(n.indices) {
				if c == idxc {
					i = n.incrementChildPrio(i)
					n = n.children[i]
					continue walk
				}
			}

			// []byte for proper unicode char conversion, see #65
			n.indices += string([]byte{idxc})
			child := &node{}
			n.children = append(n.children, child)
			n.incrementChildPrio(len(n.indices) - 1)
			n = child
			n.insertChild(path, fullPath, handlers)
			return
		}
		n.handlers = handlers
		return
	}
}

func (n *node) insertChild(path, fullPath string, handlers HandlersChain) {
	n.path = path
	n.fullPath = fullPath
	n.handlers = handlers
}

// Returns the handle registered with the given path (key). The values of
// wildcards are saved to a map.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (n *node) getValue(path string) (handlers HandlersChain, tsr bool) {
walk: // Outer loop for walking the tree
	for {
		prefix := n.path
		if len(path) > len(prefix) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]

				// If this node does not have a wildcard (param or catchAll)
				// child, we can just look up the next child node and continue
				// to walk down the tree
				idxc := path[0]
				for i, c := range []byte(n.indices) {
					if c == idxc {
						n = n.children[i]
						continue walk
					}
				}

				// Nothing found.
				// We can recommend to redirect to the same URL without a
				// trailing slash if a leaf exists for that path.
				tsr = path == "/" && n.handlers != nil
				return

			}
		} else if path == prefix {
			// We should have reached the node containing the handle.
			// Check if this node has a handle registered.
			if handlers = n.handlers; handlers != nil {
				return
			}

			// No handle found. Check if a handle for this path + a
			// trailing slash exists for trailing slash recommendation
			for i, c := range []byte(n.indices) {
				if c == '/' {
					n = n.children[i]
					tsr = len(n.path) == 1 && n.handlers != nil
					return
				}
			}
			return
		}

		// Nothing found. We can recommend to redirect to the same URL with an
		// extra trailing slash if a leaf exists for that path
		tsr = (path == "/") ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && n.handlers != nil)
		return
	}
}
