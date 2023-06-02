package main

import (
	"bufio"
	"container/heap"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Item Struct for the Priority Queue elements
// Priority is the LONG Value
type Item struct {
	url      string
	priority int64
	index    int
}

type ItemQueue []*Item

// A waitgroup to handle all the go-routines
var wg sync.WaitGroup

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ideally this would be args
	const numWorkers = 16
	const queueMaxSize = 10

	// read Filename from stdin
	var filenameReader = bufio.NewReader(os.Stdin)
	filename, err := filenameReader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	filename = strings.TrimSuffix(filename, "\n")

	// This channel is used to send every read line in various go-routines.
	msgChannel := make(chan string)

	// This Channel is used to send Items into the priority queue serializer
	itemChannel := make(chan *Item)
	defer close(itemChannel)

	// open the Input file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create the scanner for the open file
	scanner := bufio.NewScanner(f)

	// Read all lines from file and send them into channel
	// XXX: Only one routine to read the whole file, this assumes we are not IO bound.
	// TODO: For multiple routines i need to control the file seek for each routine
	go func() {
		for scanner.Scan() {
			row := scanner.Text()
			msgChannel <- row
		}

		// Signal the main thread that all lines have been read
		close(msgChannel)
	}()

	// Create the PriorityQueue
	iq := make(ItemQueue, 0)
	heap.Init(&iq)

	// Serialize access to the priority queue
	go func() {
		for i := range itemChannel {

			//fmt.Printf("ITEM: %#v\n", i)

			heap.Push(&iq, i)
			iq.update(i, i.url, i.priority)

			// If the Queue is bigger than maxSize let's pop the `lowest priority` element out
			// To keep in the list only maxSize elements with the highest priority
			if iq.Len() > queueMaxSize {
				heap.Pop(&iq)
			}

		}
	}()

	// Start the workers to process the records
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processRecord(ctx, i, msgChannel, itemChannel)
	}

	wg.Wait()

	//fmt.Printf("\nElements in queue: %d\n\n", iq.Len())

	// extract all the elements in the queue and print in reverse order
	var urls []string
	for iq.Len() > 0 {
		popItem := heap.Pop(&iq).(*Item)
		urls = append(urls, popItem.url)
	}

	for i := 1; i <= len(urls); i++ {
		fmt.Printf("%s\n", urls[len(urls)-i])
	}

}

func processRecord(ctx context.Context, id int, msgChannel chan string, itemChannel chan *Item) {
	defer wg.Done()
	localReadLines := 0
	log.Printf("%d Worker thread has started", id)

	for row := range msgChannel {
		select {
		case <-ctx.Done():
			return
		default:
			localReadLines++

			split := strings.Fields(row)
			url := split[0]

			priority, err := strconv.ParseInt(split[1], 10, 64)
			if err != nil {
				log.Printf("Failed to convert string %s to int64\n", split[1])
				continue
			}

			// Create a new Item and push it in the itemChannel
			item := &Item{
				url:      url,
				priority: priority,
			}

			itemChannel <- item

		}

	}

	log.Printf("%d Worker thread has been completed", id)
	log.Printf("%d Worker thread has processed %d lines", id, localReadLines)
}

// https://pkg.go.dev/container/heap#example-package-PriorityQueue
func (iq ItemQueue) Len() int { return len(iq) }

func (iq ItemQueue) Less(i, j int) bool {
	// We Will pop elements when the queue is too big to keep the queue at a consistent size, so we want lower than here
	return iq[i].priority < iq[j].priority
}

func (iq ItemQueue) Swap(i, j int) {
	iq[i], iq[j] = iq[j], iq[i]
	iq[i].index = i
	iq[j].index = j
}

func (iq *ItemQueue) Push(x any) {
	n := len(*iq)
	item := x.(*Item)
	item.index = n
	*iq = append(*iq, item)
}

func (iq *ItemQueue) Pop() any {
	old := *iq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*iq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (iq *ItemQueue) update(item *Item, url string, priority int64) {
	item.url = url
	item.priority = priority
	heap.Fix(iq, item.index)
}
