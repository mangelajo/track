package shell

import (
	"github.com/abiosoft/ishell"
	"github.com/mangelajo/track/pkg/bugzilla"
	"strconv"
	"github.com/mangelajo/track/pkg/show"
	"fmt"
	"time"
)

var shellBugs map[int32] *bugzilla.Cbug = make(map[int32] *bugzilla.Cbug)

func Shell(bugs *[]bugzilla.Cbug, getClient func() *bugzilla.Client) {

	bzNames := []string {}

	if len(*bugs) == 0 {
		fmt.Println("\nNo bugs for the shell, bye! :)\n")
		return
	}

	// map the bugs to shellBugs
	for _, bug := range *bugs {
		shellBugs[bug.Cbug_id.Number] = &bug
		bzNames = append(bzNames, fmt.Sprintf("%d", bug.Cbug_id.Number))
	}

	currentBugN := 0

	shell := ishell.New()
	shell.Println("Track interactive shell")

	(*bugs)[currentBugN].ShortSummary(true)

	shell.AddCmd(&ishell.Cmd{
		Name:    "open",
		Aliases: []string{"o"},
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
		Help:      "open a bugzilla from cache",
		LongHelp:  "",
		Completer: func([]string) []string { return bzNames },
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "show",
		Aliases: []string{"s"},
		Func: func(c *ishell.Context) {

			if len(c.Args) == 1 {
				bzId, err := strconv.Atoi(c.Args[0])
				if err == nil {
					fmt.Println("")
					bug, _ := shellBugs[int32(bzId)]
					if bug != nil {
						bug.ShortSummary(bugzilla.USE_COLOR)
					}
				}
			}  else if len(c.Args) == 0 {
				bug := (*bugs)[currentBugN]
				fmt.Println("")
				bug.ShortSummary(bugzilla.USE_COLOR)
			}
		},
		Help: "show a bugzilla",
		Completer: func([]string) []string { return bzNames },
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "next",
		Aliases: []string{"n"},
		Func: func(c *ishell.Context) {

			currentBugN += 1

			if currentBugN >= len(*bugs) {
				currentBugN = len(*bugs) - 1
			}

			bug := (*bugs)[currentBugN]
			fmt.Println()
			bug.ShortSummary(bugzilla.USE_COLOR)
		},
		Help: "next bugzilla",
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "prev",
		Aliases: []string{"p"},
		Func: func(c *ishell.Context) {

			currentBugN -= 1

			if currentBugN <= 0 {
				currentBugN = 0
			}

			bug := (*bugs)[currentBugN]
			fmt.Println()
			bug.ShortSummary(bugzilla.USE_COLOR)
		},
		Help: "previous bugzilla",
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "go",
		Aliases: []string{"g"},
		Func: func(c *ishell.Context) {

			if len(c.Args) == 1 {
				bzId, err := strconv.Atoi(c.Args[0])
				if err == nil {
					fmt.Println("")
					bug, _ := shellBugs[int32(bzId)]
					if bug != nil {
						fmt.Println("Opening ", bug.URL())
						show.OpenURL(bug.URL())
					}
				}
			}  else if len(c.Args) == 0 {

				bug := (*bugs)[currentBugN]
				fmt.Println("Opening ", bug.URL())
				show.OpenURL(bug.URL())
			}
		},
		Help: "open bugzilla from server url",
		Completer: func([]string) []string { return bzNames },
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "links",
		Aliases: []string{"l"},
		Func: func(c *ishell.Context) {

			links := []string{}
			bug := (*bugs)[currentBugN]

			if len(c.Args) == 1 {
				bzId, err := strconv.Atoi(c.Args[0])
				if err == nil {
					fmt.Println("")
					bugp, _ := shellBugs[int32(bzId)]
					bug = *bugp
				}
			}

			if len(bug.Cexternal_bugs)<1 {
				return
			}

			for _, link := range bug.Cexternal_bugs {
				links = append(links,
					fmt.Sprintf("%s %s", link.Attrname, link.URL()))
			}

			var choices []int

			if len(links)>1 {
				choices = c.Checklist(links,
					"Please select the links you want to open",
					nil)
			} else {
				choices = []int {0}
			}

			urls := func() (c []string) {
				for _, v := range choices {
					c = append(c, bug.Cexternal_bugs[v].URL())
				}
				return
			}
			for _, url := range urls() {
				c.Println("Opening ", url)
				show.OpenURL(url)
				// Provide some time to avoid Firefox/browser choking :)
				time.Sleep(500 * time.Millisecond)
			}
		},
		Help: "open links from bugzilla",
		Completer: func([]string) []string { return bzNames },
	})

	shell.Run()
	shell.Close()
}