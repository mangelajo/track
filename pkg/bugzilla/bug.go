package bugzilla

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strings"
	"time"
)

// Bug bugzilla bug data structure
type Bug struct {
	ID         int
	URL        string
	Product    string
	Component  string
	Assignee   string
	Status     string
	Resolution string
	Subject    string
	PMScore    int
	Severity   string
	Changed    time.Time
}

func (b *Bug) String() string {

	return fmt.Sprintf("%d (%8s)\t%s\t%s\t%20s\t%s", b.ID, b.Status,
		b.Assignee, strings.Replace(b.URL, "show_bug.cgi?id=", "", 1), b.Component, b.Subject)
}

//NewBugFromBzBug constructor for Bug
func NewBugFromBzBug(protoBug bzBug) (bug *Bug, err error) {
	bug = &Bug{
		ID:         protoBug.ID,
		URL:        protoBug.URL,
		Product:    protoBug.Product,
		Component:  protoBug.Component,
		Assignee:   protoBug.Assignee,
		Status:     protoBug.Status,
		Resolution: protoBug.Resolution,
		Subject:    protoBug.Description,
		PMScore:    protoBug.PMScore,
		Severity:   protoBug.Severity,
	}

	parser := &combinedParser{}
	t, err := parser.parse(protoBug.Changed)
	if err != nil {
		return nil, err
	}
	bug.Changed = t
	return bug, nil
}

// GetSHA1 computes SHA1 of particular bug
func (bug *Bug) GetSHA1() string {
	h := sha1.New()
	io.WriteString(h, fmt.Sprintf("%v-%v", bug.ID, bug.Changed))
	return fmt.Sprintf("%x", h.Sum(nil))
}
