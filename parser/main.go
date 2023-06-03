package main

import (
	"bufio"
	"container/heap"
	"context"
	"errors"
	"fmt"
	"io"
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

// Read the filename from stdio and return as a string
func readFilePathFromStdio() (path string, err error) {
	// read Filename from stdin
	var filenameReader = bufio.NewReader(os.Stdin)
	filename, err := filenameReader.ReadString('\n')

	if err != nil {
		return "", err
	}

	filename = strings.TrimSuffix(filename, "\n")

	if _, err := os.Stat(filename); err != nil {
		return "", errors.New(fmt.Sprintf("Filename %s is not a valid path or the file does not exist", filename))
	}

	return filename, nil
}

// Read all lines from file and send them into channel
// XXX: Only one routine to read the whole file, this assumes we are not CPU bound.
// TODO: For multiple routines i need to control the file seek for each routine
func fileReader(ctx context.Context, f io.Reader, channel chan string) {
	// Create the scanner for the open file
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		row := scanner.Text()
		channel <- row
	}

	close(channel)

}

// Worker Function
// Process each line , create an Item and push it into the itemChannel for the queue manager to handle
func processRecord(ctx context.Context, id int, readChannel chan string, itemChannel chan *Item) {
	defer wg.Done()
	localReadLines := 0

	if os.Getenv("LOG_DEBUG") != "" {
		log.Printf("%d Worker thread has started", id)
	}

	for row := range readChannel {
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

	if os.Getenv("LOG_DEBUG") != "" {
		log.Printf("%d Worker thread has been completed", id)
		log.Printf("%d Worker thread has processed %d lines", id, localReadLines)
	}
}

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ideally this would be args
	const numWorkers = 4
	const queueMaxSize = 10
	const batchSize = 1000

	// read Filename from stdin
	filename, err := readFilePathFromStdio()
	if err != nil {
		log.Fatal(err)
	}

	// This channel is used to send every read line from the read routine to the worker routines.
	readChannel := make(chan string)

	// This Channel is used to send Items into the priority queue serializer
	itemChannel := make(chan *Item)

	// open the Input file
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Start the Reader
	go fileReader(ctx, f, readChannel)

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
		go processRecord(ctx, i, readChannel, itemChannel)
	}

	// Wait for all worker threads to finish
	wg.Wait()
	// Close the items Channel now that the sender threads have completed
	close(itemChannel)

	// extract all the elements in the queue and print in reverse order
	var urls []string

	// Ensure the queue is of correct size, pop any low priority element above the MaxSize
	for iq.Len() > queueMaxSize {
		heap.Pop(&iq)
	}

	// Extract the Items into the urls slice
	for iq.Len() > 0 {
		popItem := heap.Pop(&iq).(*Item)
		urls = append(urls, popItem.url)
	}

	// Print URLs in reverse order
	for i := 1; i <= len(urls); i++ {
		fmt.Println(urls[len(urls)-i])
	}

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
