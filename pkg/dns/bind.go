package dns

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Bind represents the interface to the bind dns server
type Bind struct {
	Enabled bool
	Host    string
	Sync    bool
	Debug   bool
}

// Update adds a dns record to the bind server
func (b *Bind) Update(d Record) Result {
	cmd := []string{
		fmt.Sprintf("server %s", b.Host),
		fmt.Sprintf("zone %s", getZone(d.Hostname)),
		fmt.Sprintf("update delete %s", d.Hostname),
		fmt.Sprintf("update add %s %s %s %s", d.Hostname, d.TTL, d.Type, d.Data),
		fmt.Sprintf("send\n"),
	}

	output, err := b.invokeCommand(strings.Join(cmd, "\n"))
	if b.Sync {
		// TODO implement rndc sync -clean
	}
	return Result{Error: err, Message: output}
}

// Delete removes the dns record from the bind server
func (b *Bind) Delete(d Record) Result {
	cmd := []string{
		fmt.Sprintf("server %s", b.Host),
		fmt.Sprintf("zone %s", getZone(d.Hostname)),
		fmt.Sprintf("update delete %s", d.Hostname),
		fmt.Sprintf("send\n"),
	}

	output, err := b.invokeCommand(strings.Join(cmd, "\n"))
	if b.Sync {
		// TODO implement rndc sync -clean
	}
	return Result{Error: err, Message: output}
}

func (b *Bind) invokeCommand(t string) (string, error) {
	var f *os.File
	var err error
	if f, err = ioutil.TempFile("", "nsupd"); err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	if _, err = f.WriteString(t); err != nil {
		return "", err
	}
	if b.Debug == false {
		out, err := exec.Command("nsupdate", "-d", fmt.Sprintf("%s", f.Name())).Output()
		return string(out), err
	}

	return fmt.Sprintf("DEBUG MODE\n\n would invoke:\n\n%s", t), nil
}
