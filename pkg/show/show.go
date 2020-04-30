package show

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/viper"

	"github.com/mangelajo/track/pkg/bugzilla"
	"github.com/mangelajo/track/pkg/storecache"
)

func OpenBz(bzId int, getClient func() *bugzilla.Client) int {
	html, err := storecache.RetrieveCache(bzId, "", false)

	if err == nil {
		openHTML(bzId, html)
		return 0
	}

	client := getClient()

	html, _, err = client.ShowBugHTML(bzId, "")
	if err == nil {
		openHTML(bzId, html)
		return 0

	} else {
		fmt.Printf("Error: %s", err)
		return 1
	}
}

func openHTML(bzid int, html *[]byte) {
	filename := fmt.Sprintf("/tmp/bz%d.html", bzid)
	fmt.Printf("Wrote %s\n", filename)
	writeHTML(html, filename)
	OpenURL(filename)
}

func OpenURL(url string) {
	err := exec.Command(viper.Get("htmlOpenCommand").(string), url).Run()
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

func writeHTML(html *[]byte, outputFile string) {
	htmlStr := string(*html)

	// This rewrites the links in the html from relative to absolute
	htmlStr = strings.Replace(htmlStr, "src=\"", "src=\""+viper.Get("bzurl").(string)+"/", -1)
	htmlStr = strings.Replace(htmlStr, "href=\"", "href=\""+viper.Get("bzurl").(string)+"/", -1)
	htmlStr = strings.Replace(htmlStr, "action=\"", "action=\""+viper.Get("bzurl").(string)+"/", -1)

	f, err := os.Create(outputFile)
	defer f.Close()

	if err != nil {
		fmt.Printf("Error creating %s : %s", outputFile, err)
		os.Exit(1)
	}

	data := []byte(htmlStr)
	f.Write(data)

}
