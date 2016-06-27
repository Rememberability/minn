package main

import (
	"fmt"
	"strconv"
	"os"
	"database/sql"
	"io"
	"ivartj/args"
)

func init() {
	cmdRegister("back", cmdBack, cmdBackUsage)
}

func cmdBackUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s <deck> back [ <card-id> ]\n", mainProgramName)
}

func cmdBackArgs(cmd *cmdContext) (int, bool) {

	plainArgs := []string{}
	tok := args.NewTokenizer(cmd.Args)

	for tok.Next() {

		if tok.IsOption() {
			switch tok.Arg() {
			case "-h", "--help":
				cmdBackUsage(os.Stdout)
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

	cmdBackUsage(os.Stderr)
	cmd.Exit(1)
	return 0, false
}

func cmdBack(cmd *cmdContext) {

	cardId, current := cmdBackArgs(cmd)

	var err error
	if current {
		cardId, err = sm2CurrentCard(cmd.DB())
		if err != nil {
			panic(err)
		}
	}

	row := cmd.QueryRow("select back from cards where card_id = ?;", cardId)
	var back string
	err = row.Scan(&back)
	if err == sql.ErrNoRows {
		fmt.Fprintf(os.Stderr, "No card by that card ID.\n");
		cmd.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database query error: %s.\n", err.Error())
		cmd.Exit(1)
	}

	fmt.Println(back)

	cmd.Commit()
}
