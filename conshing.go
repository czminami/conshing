package conshing

import (
	"fmt"
	"sort"
	"strconv"
	"sync"

	"github.com/spaolacci/murmur3"
)

const (
	default_virtual_size = 200
)

type Consh struct {
	hashings    Hashings
	lock        sync.RWMutex
	nodes       int
	node        map[string]int
	virtualNode map[uint32]string
}
type Hashings []uint32

func (this Hashings) Len() int {
	return len(this)
}
func (this Hashings) Less(i, j int) bool {
	return this[i] < this[j]
}

func (this Hashings) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func hash(raw []byte) uint32 {
	return murmur3.Sum32(raw)
}

func NewConsh() *Consh {
	return &Consh{
		node:        make(map[string]int, 20),
		virtualNode: make(map[uint32]string, 20),
	}
}

func (this *Consh) NewNode(name string, virtualSize ...int) {
	this.lock.Lock()
	defer this.lock.Unlock()

	_, ok := this.node[name]
	if !ok {
		virtual_size := default_virtual_size
		if len(virtualSize) > 0 && virtualSize[0] > 0 {
			virtual_size = virtualSize[0]
		}

		this.node[name] = virtual_size

		for index := 0; index < virtual_size; index++ {
			digest := hash([]byte(fmt.Sprint(name, "-", index)))

			this.virtualNode[digest] = name
			this.hashings = append(this.hashings, digest)
		}

		this.nodes += 1
		sort.Sort(this.hashings)
	}
}

func (this *Consh) RemoveNode(name string) {
	if name == "" {
		return
	}

	this.lock.Lock()
	defer this.lock.Unlock()

	virtual_size, ok := this.node[name]
	if ok {
		delete(this.node, name)
		this.nodes -= 1

		for index := 0; index < virtual_size; index++ {
			digest := hash([]byte(fmt.Sprint(name, "-", index)))
			delete(this.virtualNode, digest)

			for i, v := range this.hashings {
				if v == digest {
					this.hashings = append(this.hashings[:i], this.hashings[i+1:]...)
					break
				}
			}
		}
	}
}

func (this *Consh) Empty() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.nodes == 0
}

func (this *Consh) LookUpNode(obj []byte, excepts ...string) string {
	if len(obj) == 0 || this.nodes == 0 || this.nodes == len(excepts) {
		return ""
	}

	var except map[string]byte
	if len(excepts) > 0 {
		except = make(map[string]byte, len(excepts))
		for _, v := range excepts {
			except[v] = 1
		}
	}

	this.lock.RLock()
	defer this.lock.RUnlock()

	var target string
	var index, loop int
	var excepted bool

	for {
		index = sort.Search(len(this.hashings), func(i int) bool {
			return this.hashings[i] >= hash(append(obj, []byte(strconv.Itoa(loop))...))
		})

		loop += 1
		if index == len(this.hashings) {
			index = 0
		}

		target = this.virtualNode[this.hashings[index]]

		if except != nil {
			if _, excepted = except[target]; !excepted {
				break
			}

		} else {
			break
		}
	}

	return target
}

func (this *Consh) List() []string {
	this.lock.RLock()
	defer this.lock.RUnlock()

	nodes := make([]string, this.nodes)

	index := 0
	for node, _ := range this.node {
		nodes[index] = node
		index += 1
	}

	return nodes
}
