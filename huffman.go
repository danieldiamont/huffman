package huffman

import (
	"container/heap"
	"fmt"
	"strconv"
)

type Huffman struct {
    Data *byte
    Sum int
    Left *Huffman
    Right *Huffman
}

type EncodedData struct {
    Codes map[byte]uint32
    Data []uint8
    Padding uint8
}

// An Item is something we manage in a priority queue.
type Item struct {
	value    *Huffman // The value of the item; arbitrary.
	priority int    // The priority of the item in the queue.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    a := pq[i]
    b := pq[j]

    if a.priority == b.priority {
        if a.value.Data != nil && b.value.Data != nil {
            return *a.value.Data < *b.value.Data
        }

        return a.value.Data != nil && b.value.Data == nil
    }

    return a.priority < b.priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value *Huffman, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}


func Print(hf *Huffman, name string) {
    fmt.Println("Tree name: ", name)
    
    var helper func(hf *Huffman, tabs string) 

    helper = func(hf *Huffman, tabs string) {
        if hf == nil {
            return
        }
        ch := byte('*')
        if hf.Data != nil {
            ch = *hf.Data
        }
        fmt.Printf("%s%+v\tcharacter = %c\n", tabs, hf, ch)
        tabs += "\t"
        helper(hf.Left, tabs)
        helper(hf.Right, tabs)
    }

    helper(hf, "")
}

func CompareHuffman(a, b *Huffman) bool {
    if (a == nil && b == nil) {
        return true
    }
    if (a != nil && b == nil || a == nil && b != nil) {
        return false
    }

    if (a.Data == nil || b.Data == nil) {
        if (a.Data == nil && b.Data == nil) {
            return a.Sum == b.Sum &&
                CompareHuffman(a.Left, b.Left) && 
                CompareHuffman(a.Right, b.Right)
        }
        return false
    }

    return *a.Data == *b.Data && 
        a.Sum == b.Sum && 
        CompareHuffman(a.Left, b.Left) && 
        CompareHuffman(a.Right, b.Right)
}


func BuildHuffman(f map[byte]int) (*Huffman, error) {
    if f == nil || len(f) < 2 {
        e := fmt.Errorf("Expected f to not be nil or len(f) > 2")
        return nil, e
    }

    pq := PriorityQueue{}

    for ch, priority := range f {
        node := &Huffman{
            &ch,
            priority,
            nil,
            nil,
        }
        pq = append(pq, &Item{
            value: node,
            priority: priority,
        })
    }

    heap.Init(&pq)
    for pq.Len() >= 2 {

        first := heap.Pop(&pq).(*Item)
        second := heap.Pop(&pq).(*Item)
        total := first.value.Sum + second.value.Sum

        node := &Huffman{
            nil,
            total,
            first.value,
            second.value,
        }
        item := &Item{
            value: node,
            priority: total,
        }

        heap.Push(&pq, item)
    }

    return heap.Pop(&pq).(*Item).value, nil
}

func GetTable(hf *Huffman) (map[byte]uint32, error) {
    f := make(map[byte]uint32)

    var table func(hf * Huffman, code string) error
    table = func(hf * Huffman, code string) error {
        if hf == nil {
            return nil
        }

        if hf.Data != nil {
            ch := *hf.Data
            val, err := strconv.ParseUint(code, 2, 32)
            if err != nil {
                return fmt.Errorf("%v\n", err)
            }
            f[ch] = uint32(val)
            return nil
        } else {
            err := table(hf.Left, code + "0")
            if err != nil {
                return err
            }
            table(hf.Right, code + "1")
            if err != nil {
                return err
            }
            return nil
        }
    }

    err := table(hf, "")
    if err != nil {
        return nil, err
    }

    return f, nil
}

func Encode(codes map[byte]uint32, data []byte) (*EncodedData, error) {

    var padding uint8
    var bitPosition int
    var byteResult uint8

    bitPosition = 7
    byteResult = 0
    payload := make([]byte, 0)

    for _, ch := range data {
        _, prs := codes[ch]
        if !prs {
            return nil, fmt.Errorf("Did not find byte = %v in codes\n", ch)
        }

        code := codes[ch] // uint32
        stack := make([]uint8, 0)

        if code != 0 {
            for code != 0 {
                bit := uint8(code & 1) 
                stack = append(stack, bit)
                code >>= 1
            }

            for i := len(stack)-1; i >= 0; i-- {
                bit := stack[i]
                byteResult |= (bit << bitPosition)

                bitPosition -= 1
                if bitPosition == -1 {
                    bitPosition = 7
                    payload = append(payload, byteResult)
                    byteResult = 0
                }
            }
        } else {
            bitPosition -= 1
            if bitPosition == -1 {
                payload = append(payload, byteResult)
                bitPosition = 7
                byteResult = 0
            }
        }
    }

    padding = 0

    if bitPosition < 7 {
        payload = append(payload, byteResult)
        pos := uint8(bitPosition)
        padding = pos + 1
    }


    e := &EncodedData{}
    e.Codes = codes
    e.Data = payload
    e.Padding = padding

    return e, nil
}
