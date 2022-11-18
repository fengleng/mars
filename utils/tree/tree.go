package tree

import (
	"errors"
	"fmt"
	"github.com/gososy/sorpc/log"
)

type NodeData interface {
	GetId() string
	GetParentId() string
	String() string
}
type Node struct {
	Data    NodeData
	Parent  *Node
	Child   *Node
	Sibling *Node
	// cache
	id       string
	parentId string
}

func (p *Node) GetId() string {
	return p.id
}
func (p *Node) GetParentId() string {
	return p.parentId
}
func NewNode(data NodeData) *Node {
	n := &Node{
		Data:     data,
		id:       data.GetId(),
		parentId: data.GetParentId(),
	}
	if n.id == "" {
		log.Fatal("invalid id")
	}
	return n
}
func (p *Node) ResetData(data NodeData) {
	p.Data = data
	p.id = data.GetId()
	if p.id == "" {
		log.Fatal("invalid id")
	}
	p.parentId = data.GetParentId()
}
func (p *Node) Dump(level int) {
	for i := 0; i < level; i++ {
		fmt.Printf("  ")
	}
	fmt.Printf("%s\n", p.Data.String())
	c := p.Child
	for c != nil {
		c.Dump(level + 1)
		c = c.Sibling
	}
}
func (p *Node) VisitChild(v func(node *Node), recursive bool) {
	c := p.Child
	for c != nil {
		v(c)
		c = c.Sibling
	}
	if recursive {
		c = p.Child
		for c != nil {
			c.VisitChild(v, true)
			c = c.Sibling
		}
	}
}
func (p *Node) VisitParent(v func(node *Node), recursive bool) {
	c := p.Parent
	if c == nil {
		return
	}
	v(c)
	if recursive {
		for {
			c = c.Parent
			if c == nil {
				break
			}
			v(c)
		}
	}
}
func (p *Node) VisitAll(v func(node *Node)) {
	v(p)
	p.VisitChild(v, true)
}
func (p *Node) IsSub(sub *Node) bool {
	if sub == p {
		return false
	}
	i := sub.Parent
	for i != nil {
		if i == p {
			return true
		}
		i = i.Parent
	}
	return false
}
func (p *Node) Attach(sub *Node) {
	if p.Child == nil {
		p.Child = sub
		sub.Sibling = nil
	} else {
		sub.Sibling = p.Child
		p.Child = sub
	}
	sub.Parent = p
}
func (p *Node) Detach(sub *Node) {
	if p.Child == sub {
		p.Child = nil
	} else {
		c := p.Child
		var pre *Node
		for c != nil {
			if c == sub {
				if pre != nil {
					pre.Sibling = sub.Sibling
					sub.Sibling = nil
					sub.Parent = nil
				} else {
					p.Child = nil
				}
				break
			}
			pre = c
			c = c.Sibling
		}
	}
}

type Tree struct {
	rootMap        map[string]*Node
	hash           map[string]*Node
	allowMultiRoot bool
}

func NewTree() *Tree {
	return &Tree{
		rootMap: map[string]*Node{},
		hash:    make(map[string]*Node),
	}
}
func BuildTree(nodes []NodeData, allowMultiRoot bool) (*Tree, error) {
	t := &Tree{
		hash: make(map[string]*Node),
	}
	err := t.Build(nodes, allowMultiRoot)
	if err != nil {
		return nil, err
	}
	return t, nil
}
func (p *Tree) IsEmpty() bool {
	return len(p.rootMap) == 0
}
func (p *Tree) GetHash() map[string]*Node {
	return p.hash
}
func (p *Tree) Build(nodes []NodeData, allowMultiRoot bool) error {
	// to hash table
	var m = make(map[string]*Node)
	var rootMap = map[string]*Node{}
	for _, v := range nodes {
		node := NewNode(v)
		m[node.id] = node
		if node.parentId == "" {
			if len(rootMap) > 0 && !allowMultiRoot {
				return ErrDupRoot
			}
			rootMap[node.id] = node
		}
	}
	if len(rootMap) == 0 && len(nodes) > 0 {
		return errors.New("missed root")
	}
	// to tree
	for _, v := range m {
		pid := v.parentId
		if pid != "" {
			parentNode := m[pid]
			node := m[v.id]
			if parentNode == nil {
				return fmt.Errorf("parent not found %s %s", pid, v.id)
			}
			// 看看树有没有环绕
			tmp := parentNode
			for tmp != nil {
				if tmp.id == node.id {
					return fmt.Errorf("tree node dead loop %s", node.id)
				}
				tmp = tmp.Parent
			}
			parentNode.Attach(node)
		}
	}
	p.rootMap = rootMap
	p.hash = m
	p.allowMultiRoot = allowMultiRoot
	return nil
}

var ErrInvalidParentId = errors.New("ErrInvalidParentId")
var ErrParentNotFound = errors.New("ErrParentNotFound")
var ErrDepartmentNotFound = errors.New("ErrDepartmentNotFound")
var ErrDepartmentNotLeaf = errors.New("ErrDepartmentNotLeaf")
var ErrBeforeCheckFail = errors.New("ErrBeforeCheckFail")
var ErrCanNotMove = errors.New("ErrCanNotMove")
var ErrDupRoot = errors.New("dup root")

func (p *Tree) Dump() {
	if len(p.rootMap) == 0 {
		fmt.Printf("<nil tree>\n")
	} else {
		for _, v := range p.rootMap {
			v.Dump(0)
		}
	}
}
func (p *Tree) Add(d NodeData, beforeAdd func(parentNode *Node) bool) error {
	var parentNode *Node
	if d.GetParentId() == "" {
		if len(p.rootMap) > 0 && !p.allowMultiRoot {
			return ErrInvalidParentId
		}
	} else {
		parentNode = p.hash[d.GetParentId()]
		if parentNode == nil {
			return ErrParentNotFound
		}
	}
	if beforeAdd != nil && !beforeAdd(parentNode) {
		return ErrBeforeCheckFail
	}
	node := NewNode(d)
	p.hash[node.id] = node
	if parentNode != nil {
		parentNode.Attach(node)
	}
	if d.GetParentId() == "" {
		p.rootMap[node.id] = node
	}
	return nil
}
func (p *Tree) Del(id string, beforeDel func(node *Node) bool) error {
	node := p.hash[id]
	if node == nil {
		return ErrDepartmentNotFound
	}
	// 非页子不能删除
	if node.Child != nil {
		return ErrDepartmentNotLeaf
	}
	if beforeDel != nil && !beforeDel(node) {
		return ErrBeforeCheckFail
	}
	parent := node.Parent
	if parent != nil {
		parent.Detach(node)
	}
	delete(p.hash, node.id)
	// 如果是顶级结点
	delete(p.rootMap, node.id)
	return nil
}
func (p *Tree) Move(id string, newParentId string, beforeMove func() bool) error {
	node := p.hash[id]
	if node == nil {
		return ErrDepartmentNotFound
	}
	oldParent := node.Parent
	if oldParent == nil {
		return ErrCanNotMove
	}
	newParent := p.hash[newParentId]
	if newParent == nil {
		return ErrParentNotFound
	}
	if node.IsSub(newParent) {
		// 我不能移去我的子部门里边
		return ErrCanNotMove
	}
	if beforeMove != nil && !beforeMove() {
		return ErrBeforeCheckFail
	}
	oldParent.Detach(node)
	node.parentId = newParentId
	newParent.Attach(node)
	return nil
}
func (p *Tree) Get(id string) *Node {
	return p.hash[id]
}
func (p *Tree) GetRoot() *Node {
	if len(p.rootMap) > 0 {
		for _, v := range p.rootMap {
			return v
		}
	}
	return nil
}
func (p *Tree) GetRootMap() map[string]*Node {
	return p.rootMap
}
func (p *Tree) VisitAll(v func(node *Node)) {
	for _, n := range p.rootMap {
		n.VisitAll(v)
	}
}
