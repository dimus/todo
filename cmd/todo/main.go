package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dimus/todo"
	"github.com/gogo/protobuf/proto"
)

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "need add or list cmd")
		os.Exit(1)
	}
	var err error
	switch cmd := flag.Arg(0); cmd {
	case "list":
		err = list()
	case "add":
		err = add(strings.Join(flag.Args()[1:], " "))
	default:
		err = fmt.Errorf("unknown subcomand %s", cmd)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const dbPath = "mydb.pb"

func add(text string) error {
	task := &todo.Task{
		Text: text,
		Done: false,
	}
	b, err := proto.Marshal(task)
	if err != nil {
		return fmt.Errorf("Cannot marshal task: %v", err)
	}

	f, err := os.OpenFile(dbPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("Cannot open filel %f: %v", dbPath, err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("Cannot write to file %s: %v", dbPath, err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("Cannot close file %s: %v", dbPath, err)
	}
	fmt.Println(proto.MarshalTextString(task))
	return nil
}

func list() error {
	return nil
}
