# bugzilla
Golang client for bugzilla API


### list bugs
```go
client, err := bugzilla.NewClient(bugzillaURL, bugzillaUser, bugzillaPassword)
if err != nil {
	return err
}
client.BugList(limit, offset)
```

### bug details for #444
```go
client, err := bugzilla.NewClient(bugzillaURL, bugzillaUser, bugzillaPassword)
if err != nil {
	return err
}
client.BugInfo(444)
```

### add comment to #444
```go
client, err := bugzilla.NewClient(bugzillaURL, bugzillaUser, bugzillaPassword)
if err != nil {
	return err
}
client.AddComment(444, "Hello word!"))
```
