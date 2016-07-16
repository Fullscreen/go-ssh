go-ssh
======

go-ssh is a simple wrapper around exec calls to ssh for programmatic initiation
ssh session and commands.

Example
=======

```go
if err := ssh.New(host).Connect(); err != nil {
	return err
}
```
