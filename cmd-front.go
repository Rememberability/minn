package main

import (
	"fmt"
	"strconv"
	"os"
	"database/sql"
	"ivartj/args"
	"io"
)

func init() {
	cmdRegister("front", cmdFront, cmdFrontUsage)
}

func cmdFrontUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> front [ <card-id> ]\n", mainProgramName)
}

func cmdFrontArgs(cmd *cmdContext) (int, bool) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				cmdFrontUsage(os.Stdout)
				cmd.Exit(0)
			default:
				fmt.Fprintf(os.Stderr, "Unrecognized option: %s.\n", tok.Arg())
				cmd.Exit(1)
			}
		} else {
			plainArgs = append(plainArgs, tok.Arg())
		}
	}

	if tok.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error on processing command-line arguments: %s.\n", tok.Err().Error())
		cmd.Exit(1)
	}

	switch(len(plainArgs)) {
	case 0:
		return 0, true
	case 1:
		cardId, err := strconv.Atoi(plainArgs[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse Card ID: %s.\n", err.Error())
			cmd.Exit(1)
		}
		return cardId, false
	}

	cmdFrontUsage(os.Stderr)
	cmd.Exit(1)
	return 0, false
}

func cmdFront(cmd *cmdContext) {

	cardId, current := cmdFrontArgs(cmd)

	var err error
	if current {
		cardId, err = sm2CurrentCard(cmd.DB())
		if err != nil {
			panic(err)
		}
	}

	row := cmd.QueryRow("select front from cards where card_id = ?;", cardId)
	var front string
	err = row.Scan(&front)
	if err == sql.ErrNoRows {
		fmt.Fprintf(os.Stderr, "No card by that card ID.\n");
		cmd.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database query error: %s.\n", err.Error())
		cmd.Exit(1)
	}

	fmt.Println(front)

	cmd.Commit()
}
