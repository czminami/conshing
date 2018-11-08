package conshing

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_Add_Remove(t *testing.T) {
	consh := NewConsh()

	consh.NewNode("node1")
	consh.NewNode("node2")
	consh.NewNode("node3")
	t.Log(consh.Empty(), consh.List())

	consh.RemoveNode("node1")
	t.Log(consh.Empty(), consh.List())

	consh.RemoveNode("node2")
	t.Log(consh.Empty(), consh.List())

	consh.RemoveNode("node3")
	t.Log(consh.Empty(), consh.List())
}

func Test_LookUp(t *testing.T) {
	consh := NewConsh()

	consh.NewNode("node1")
	consh.NewNode("node2")
	consh.NewNode("node3")
	consh.NewNode("node4")

	for k := 0; k < 10; k++ {
		t.Log("$", k, consh.LookUpNode([]byte(strconv.Itoa(k))))
	}

	for k := 0; k < 10; k++ {
		t.Log("#", k, consh.LookUpNode([]byte(strconv.Itoa(k)), "node1", "node2")) // excepts node1 and node2
	}
}

func Test_DistributionBalance(t *testing.T) {
	consh := NewConsh()

	consh.NewNode("node1")
	consh.NewNode("node2")
	consh.NewNode("node3")
	consh.NewNode("node4")

	total := 1000000
	dist := make(map[string]int, 4)

	for k := 0; k < total; k++ {
		node := consh.LookUpNode([]byte(strconv.Itoa(k)))
		if size, ok := dist[node]; ok {
			dist[node] = size + 1
		} else {
			dist[node] = 1
		}
	}

	for name, size := range dist {
		t.Log(fmt.Sprintf("%s %.6f", name, float32(size)/float32(total)))
	}
}
