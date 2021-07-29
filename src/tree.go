package engine

import "strings"

// 前缀树路由实现，用以支持动态路由
// 通配符 `*`：例如 /rootCA/* 可以匹配 /rootCA/show、/rootCA/verify/log
// 参数匹配 `:`：例如 /rootCA/:/log 可以匹配 /rootCA/a/log、/rootCA/b/log

// 树节点
type node struct {
	pattern  string  //存储动态路由
	part     string  // 存储动态路由中 `/` 分割的一段
	children []*node // 子节点组
	isWild   bool    // 精确匹配
}

// 查找匹配的第一个子节点
func (n *node) findChild(part string) *node {
	// 遍历子节点
	for _, child := range n.children {
		// 找到子节点，或者子节点为模糊匹配，认为匹配到
		if child.part == part || child.isWild {
			return n
		}
	}
	return nil
}

func (n *node) findChildren(part string) []*node {
	children := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			children = append(children, child)
		}
	}
	return children
}

// 路由注册
// parts: 一条路由将被拆分成字符切片{"rootCA","verify","log"}
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.findChild(part)
	if child == nil {
		/*
		 * 若子节点不包含 part，则新建一个节点存放 part。
		 * 若 part 中包含 `*` 和 `:`，则认为是模糊匹配
		 */
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		//新建节点添加在子节点组中
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		/*
		 * len(parts) 代表动态路由的层次，
		 * 如果查到了最后一层或者当前路由部分包含通配符，
		 * pattern还为空，表明路由表中没有这条路由
		 */
		if n.pattern == "" {
			return nil
		}
		return n
	}
	/*
	 * 若还没有查到最后一层
	 */
	part := parts[height]
	//存在模糊匹配，因此先把模糊匹配和精确匹配的节点都收集起来
	children := n.findChildren(part)
	for _, child := range children {
		//遍历，查找下一个 part 是否相同
		//若也相同，则返回该节点
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
