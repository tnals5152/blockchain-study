package blockchain

import (
	"flag"
	"os"
)

type CLI struct{}

// TODO: 이름으로 메소드 실행
var cmdMatch map[string]string = map[string]string{
	"new":  "newCmd",
	"add":  "addCmd",
	"list": "listCmd",
}

func (cli *CLI) createBlockchain() {
	blockchain := NewBlockchain()
	blockchain.db.Close()
}

func (cli *CLI) addBlock(data string) { //새로운 블록 추가(기존에 있던 체인에 추가)
	blockchain := GetBlockchain()
	defer blockchain.db.Close()

	blockchain.AddBlock(data)
}

func (cli *CLI) list() {
	blockchain := GetBlockchain()
	defer blockchain.db.Close()

	blockchain.List()
}

func (cli *CLI) Run() {
	newCmd := flag.NewFlagSet("new", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	addBlockData := addCmd.String("data", "", "")

	switch os.Args[1] {
	case "new":
		newCmd.Parse(os.Args[2:])
	case "add":
		addCmd.Parse(os.Args[2:])
	case "list":
		listCmd.Parse(os.Args[2:])
	default:
		os.Exit(1)
	}

	if newCmd.Parsed() {
		cli.createBlockchain()
	}

	if addCmd.Parsed() {
		if *addBlockData == "" {
			addCmd.Usage()
			os.Exit(1)
		}

		cli.addBlock(*addBlockData)
	}
	if listCmd.Parsed() {
		cli.list()
	}
}
