// +build go1.3,!plan9,!solaris

package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func (vt *ValidatorTest) UpdateEmailFile(t *testing.T, emails []string) {
	var err error
	vt.auth_email_file, err = os.OpenFile(
		vt.auth_email_file.Name(), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		t.Fatal("failed to re-open temp file for updates")
	}
	vt.WriteEmails(t, emails)
}

func (vt *ValidatorTest) UpdateEmailFileViaRenameAndReplace(
	t *testing.T, emails []string) {
	orig_file := vt.auth_email_file
	var err error
	vt.auth_email_file, err = ioutil.TempFile("", "test_auth_emails_")
	if err != nil {
		t.Fatal("failed to create temp file for rename and replace: " +
			err.Error())
	}
	vt.WriteEmails(t, emails)

	moved_name := orig_file.Name() + "-moved"
	err = os.Rename(orig_file.Name(), moved_name)
	err = os.Rename(vt.auth_email_file.Name(), orig_file.Name())
	if err != nil {
		t.Fatal("failed to rename and replace temp file: " +
			err.Error())
	}
	vt.auth_email_file = orig_file
	os.Remove(moved_name)
}

func TestValidatorOverwriteEmailListViaRenameAndReplace(t *testing.T) {
	vt := NewValidatorTest(t)
	defer vt.TearDown()

	vt.WriteEmails(t, []string{"xyzzy@example.com"})
	domains := []string(nil)
	updated := make(chan bool, 1)
	validator := vt.NewValidator(domains, updated)

	if !validator("xyzzy@example.com") {
		t.Error("email in list should validate")
	}

	vt.UpdateEmailFileViaRenameAndReplace(t, []string{"plugh@example.com"})
	<-updated

	if validator("xyzzy@example.com") {
		t.Error("email removed from list should not validate")
	}
}
