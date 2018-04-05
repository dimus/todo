package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
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
		return fmt.Errorf("Cannot open filel %s: %v", dbPath, err)
	}

	if err := gob.NewEncoder(f).Encode(int64(len(b))); err != nil {
		return fmt.Errorf("Cannot encode length of a task to file: %v", err)
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("Cannot write to file %s: %v", dbPath, err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("Cannot close file %s: %v", dbPath, err)
	}
	fmt.Println("Creating new record")
	fmt.Println(proto.MarshalTextString(task))
	return nil
}

func list() error {
	b, err := ioutil.ReadFile(dbPath)
	if err != nil {
		return fmt.Errorf("Cannot read file %s: %v", dbPath, err)
	}
	for {
		if len(b) == 0 {
			return nil
		} else if len(b) < 4 {
			return fmt.Errorf("db file is too small: %d bytes", len(b))
		}

		var length int64
		err = gob.NewDecoder(bytes.NewReader(b[:4])).Decode(&length)
		if err != nil {
			return fmt.Errorf("Cannot decode the legth of next message: %v", err)
		}
		b = b[4:]

		var task todo.Task
		if err := proto.Unmarshal(b[:length], &task); err != nil {
			return fmt.Errorf("Cannot read task: %v", err)
		}
		b = b[length:]
		if task.Done {
			fmt.Print("👌 : ")
		} else {
			fmt.Print("⏰ : ")
		}
		fmt.Printf("%s\n", task.Text)
	}
	return nil
}
