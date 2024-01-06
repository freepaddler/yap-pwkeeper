package main

import (
	"fmt"
	"testing"
	"time"
)

type Note struct {
	note string
}

type Card struct {
	card string
}

func receiver(a, b chan struct{}) {
	fmt.Println("started")
	for {
		select {
		case <-a:
			fmt.Println("a closed")
		case <-b:
			fmt.Println("b closed")
		}
	}
	fmt.Println("finished")
}

func receiver2(ch chan interface{}) {
	fmt.Println("started")
	for {
		got, ok := <-ch
		if !ok {
			fmt.Println("channel closed. exiting")
			return
		}
		switch got.(type) {
		case Note:
			q := got.(Note)
			fmt.Println("note: ", q.note)
		case Card:
			q := got.(Card)
			fmt.Println("card: ", q.card)
		default:
			fmt.Println("unknown type")
		}
	}
}

func Test_just(t *testing.T) {
	s := make(chan interface{})
	go receiver2(s)
	i := 4
	s <- i
	n := Note{note: "this is a note"}
	c := Card{card: "this is a card"}
	s <- n
	s <- c
	close(s)
	time.Sleep(time.Second)
}
