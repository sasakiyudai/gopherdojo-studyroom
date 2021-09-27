package main

import (
	"testing"
	"os"
)

func errExists(t *testing.T, err error) bool {
	t.Helper()
	if err != nil {
		return true
	}
	return false
}

func TestValidateArgs(t *testing.T) {
	imgDirDataTests := []struct {
		caseName           string
		src          string
		dst string
		dirPath            string
		errRaisedFlag      bool
	}{
		{"nodir1", "png", "jpg", "", true},
		{"nodir2", "jpg", "svg", "", true},
		{"ok1", "png", "jpg", "testdatahoge", false},
		{"ok2", "jpg", "gif", "testdata", false},
		{"ok3", "png", "jpeg", "testdatahoge", false},
		{"invalid ext", "jjj", "ppp", "testdata", true},
	}

	backupArgs := os.Args
	for _, tt := range imgDirDataTests {
		t.Run(tt.caseName, func(t *testing.T) {
			os.Args = []string{"", "-src",tt.src, "-dst", tt.dst, tt.dirPath}
			
			err := parseArgs()
			if errExists(t, err) != tt.errRaisedFlag {
				t.Errorf("error: %#v", err)
			}
		})
		os.Args = backupArgs
	}
}
