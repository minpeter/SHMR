package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/minpeter/SHMR/pkg/de"
)

func main() {
	// 서브 커맨드를 구현하기 위한 플래그
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// add 커맨드의 플래그
	addURL := addCmd.String("url", "", "URL to add")
	addToken := addCmd.String("token", "", "Token to add")

	// remove 커맨드의 플래그
	removeID := removeCmd.String("id", "", "ID to remove")
	removeToken := removeCmd.String("token", "", "Token to remove")

	// list 커맨드의 플래그
	// 없음

	// 적어도 하나의 서브 커맨드가 제공되어야 함
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli-tool [command]")
		fmt.Println("Available commands: add, remove, list")
		os.Exit(1)
	}

	// 각 서브 커맨드에 따른 처리
	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addURL == "" || *addToken == "" {
			fmt.Println("Please provide URL and token to add")
			addCmd.PrintDefaults()
			os.Exit(1)
		}
		add(*addURL, *addToken)
	case "remove":
		removeCmd.Parse(os.Args[2:])
		if *removeID == "" && *removeToken == "" {
			fmt.Println("Please provide an ID to remove")
			removeCmd.PrintDefaults()
			os.Exit(1)
		}
		remove(*removeID, *removeToken)
	case "list":
		listCmd.Parse(os.Args[2:])
		list()
	default:
		fmt.Println("Invalid command")
		fmt.Println("Available commands: add, remove, list")
		os.Exit(1)
	}
}

func add(url, token string) {
	// add 커맨드 로직 구현
	fmt.Println("Adding:", url, token)

	id, err := de.New(url, token)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Added successfully", id)
}

func remove(id, token string) {
	// remove 커맨드 로직 구현

	err := de.Remove(id, token)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Removing ID:", id)

}

func list() {
	// list 커맨드 로직 구현
	fmt.Println("Listing items")
	// map[string]string de.List()
	list, err := de.List()
	if err != nil {
		fmt.Println(err)
		return
	}

	for id, state := range list {
		fmt.Println(id, state)
	}

}
