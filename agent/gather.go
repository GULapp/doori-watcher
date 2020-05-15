package main

type Gather interface {
	Gathering() error
	Display()
}
