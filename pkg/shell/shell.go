package shell

import (
	"github.com/mangelajo/ishell"
	"github.com/mangelajo/track/pkg/bugzilla"
	"strconv"
	"github.com/mangelajo/track/pkg/show"
	"fmt"
)

var shellBugs map[int32] *bugzilla.Cbug = make(map[int32] *bugzilla.Cbug)

func Shell(bugs *[]bugzilla.Cbug, getClient func() *bugzilla.Client) {

	// map the bugs to shellBugs
	for _, bug := range *bugs {
		shellBugs[bug.Cbug_id.Number] = &bug
	}

	currentBugN := 0



	shell := ishell.New()
	shell.Println("Track interactive shell")

	(*bugs)[currentBugN].ShortSummary(true)

	shell.AddCmd(&ishell.Cmd{
		Name:    "open",
		Aliases: nil,
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				bzId, err := strconv.Atoi(c.Args[0])
				if err == nil {
					show.OpenBz(bzId, getClient)
				}
			} else if len(c.Args) == 0 {
				bug := (*bugs)[currentBugN]
				show.OpenBz(int(bug.Cbug_id.Number), getClient)
			}

		},
		Help:      "open a bugzilla",
		LongHelp:  "",
		Completer: nil,
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "show",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 1 {
				bzId, err := strconv.Atoi(c.Args[0])
				if err == nil {
					fmt.Println("")
					shellBugs[int32(bzId)].ShortSummary(bugzilla.USE_COLOR)
				}  else if len(c.Args) == 0 {
					bug := (*bugs)[currentBugN]
					fmt.Println("")
					bug.ShortSummary(bugzilla.USE_COLOR)
				}
			}

		},
		Help: "show a bugzilla",
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "next",
		Func: func(c *ishell.Context) {

			currentBugN += 1

			bug := (*bugs)[currentBugN]
			fmt.Println()
			bug.ShortSummary(bugzilla.USE_COLOR)
		},
		Help: "next bugzilla",
	})

	shell.Run()
	shell.Close()
}