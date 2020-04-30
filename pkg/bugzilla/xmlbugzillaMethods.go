package bugzilla

import (
	"fmt"

	"github.com/fatih/color"
)

func (bz *Cbugzilla) Authenticated() bool {
	// when a bug was grabbed as an authenticated user the login will be in the
	// exporter attribute
	return bz.Attrexporter != ""
}

func (extbug *Cexternal_bugs) URL() string {

	var sfmt string

	switch extbug.Attrname {
	case "Red Hat Customer Portal":
		sfmt = "https://access.redhat.com/support/cases/%s"
	case "Red Hat Knowledge Base (Solution)":
		sfmt = "https://access.redhat.com/site/solutions/%s"
	case "Red Hat Engineering Gerrit":
		sfmt = "https://code.engineering.redhat.com/gerrit/#/c/%s"
	case "OpenStack gerrit":
		sfmt = "https://review.openstack.org/#/c/%s/"
	case "OpenStack Storyboard":
		sfmt = "https://storyboard.openstack.org/#!/story/%s"
	case "Launchpad":
		sfmt = "https://bugs.launchpad.net/bugs/%s"
	case "Trello":
		sfmt = "https://trello.com/c/%s"
	default:
		sfmt = "%s"
	}

	id := extbug.Content
	return fmt.Sprintf(sfmt, id)
}

func (bug *Cbug) URL() string {
	return fmt.Sprintf("http://bugzilla.redhat.com/%d", bug.Cbug_id.Number)
}

const USE_COLOR = true
const NO_COLOR = false

func (bi *Cbug) ShortSummary(useColor bool) {
	if useColor {
		color.Set(color.FgWhite, color.Bold)
	}
	fmt.Printf("BZ %d (%8s) %s\n", bi.Cbug_id.Number, bi.Cbug_status.Content, bi.Cshort_desc.Content)
	if useColor {
		color.Unset()
	}
	fmt.Printf("  Product: %s ver: %s target: %s (%s)\n", bi.Cproduct.Content, bi.Cversion.Content,
		bi.Ctarget_release.Content, bi.Ctarget_milestone.Content)
	fmt.Printf("  Keywords: %s\n", bi.Ckeywords.Content)
	if bi.Cassigned_to != nil {
		fmt.Printf("  Assigned to: %s\n", bi.Cassigned_to.Content)
	}
	fmt.Printf("  * bugzilla: %s\n", bi.URL())
	for _, x := range bi.Cexternal_bugs {

		fmt.Printf("  * %s : %s\n", x.Attrname, x.URL())
	}
	fmt.Println("")
}
