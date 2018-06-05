package bugzilla

import (
	"crypto/sha1"
	"fmt"
	"io"
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
	Changed    time.Time
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
