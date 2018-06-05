package bugzilla

import (
	"encoding/xml"
	"io"
)

/*
<?xml version="1.0" encoding="UTF-8"?>
<!--  -->
<RDF xmlns="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
     xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
     xmlns:bz="http://www.bugzilla.org/rdf#"
     xmlns:nc="http://home.netscape.com/NC-rdf#">

<bz:result rdf:about="https://bugzilla.mozilla.org/buglist.cgi?ctype=rdf&amp;limit=1&amp;list_id=1074511&amp;query_format=advanced&amp;offset=0">
  <bz:installation rdf:resource="https://bugzilla.mozilla.org/" />
  <bz:query_timestamp>2015-06-30 15:33:15 PDT</bz:query_timestamp>
  <bz:bugs>
    <Seq>
      <li>
        <bz:bug rdf:about="https://bugzilla.mozilla.org/show_bug.cgi?id=25666">
          <bz:id nc:parseType="Integer">25666</bz:id>
          <bz:product>Test Product1</bz:product>
          <bz:component>Ordering</bz:component>
          <bz:assigned_to>iosdev1&#64;mozilla.com</bz:assigned_to>
          <bz:bug_status>CONFIRMED</bz:bug_status>
          <bz:resolution></bz:resolution>
          <bz:short_desc>Bla-bla text</bz:short_desc>
          <bz:changeddate>15:32:49</bz:changeddate>
        </bz:bug>
      </li>
    </Seq>
  </bz:bugs>
</bz:result>

</RDF>
*/

type bzRDF struct {
	XMLName xml.Name `xml:"RDF"`
	Result  bzResult `xml:"result"`
}

type bzResult struct {
	XMLName      xml.Name `xml:"result"`
	Installation string   `xml:"installation"`
	Timestamp    string   `xml:"query_timestamp"`
	Bugs         bzBugs   `xml:"bugs"`
}

type bzBugs struct {
	XMLName xml.Name `xml:"bugs"`
	Seq     bzSeq    `xml:"Seq"`
}

type bzSeq struct {
	XMLName xml.Name `xml:"Seq"`
	Items   []bzLi   `xml:"li"`
}

type bzLi struct {
	XMLName xml.Name `xml:"li"`
	Bug     bzBug    `xml:"bug"`
}

// bzBug summary information from bugzilla ticket
type bzBug struct {
	XMLName     xml.Name `xml:"bug"`
	ID          int      `xml:"id"`
	URL         string   `xml:"about,attr"`
	Product     string   `xml:"product"`
	Component   string   `xml:"component"`
	Assignee    string   `xml:"assigned_to"`
	Status      string   `xml:"bug_status"`
	Resolution  string   `xml:"resolution"`
	Description string   `xml:"short_desc"`
	Changed     string   `xml:"changeddate"`
}

// parseBugzRDF returns list of bugs and their hashes
func parseBugzRDF(reader io.Reader) (results []bzBug, err error) {
	var rdf bzRDF
	err = xml.NewDecoder(reader).Decode(&rdf)
	if err != nil {
		return nil, err
	}

	results = make([]bzBug, len(rdf.Result.Bugs.Seq.Items), len(rdf.Result.Bugs.Seq.Items))
	for i, container := range rdf.Result.Bugs.Seq.Items {
		results[i] = container.Bug
	}
	return results, nil
}
