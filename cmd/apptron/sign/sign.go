package sign

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

func Sign(dir, bundleName, exePath string) (error, *bytes.Buffer) {
	entitlementsPlist := filepath.Join(dir, "entitlements.plist")
	if err := ioutil.WriteFile(entitlementsPlist, entitlements([]string{
		"com.apple.security.inherit",
		"com.apple.security.automation.apple-events",
		"com.apple.security.network.client",
		"com.apple.security.network.server",
		"com.apple.security.files.user-selected.executable",
		"com.apple.security.files.user-selected.read-write",
		"com.apple.security.cs.allow-jit",
		"com.apple.security.cs.allow-unsigned-executable-memory",
		"com.apple.security.cs.allow-dyld-environment-variables",
		"com.apple.security.cs.disable-library-validation",
		"com.apple.security.cs.disable-executable-page-protection",
		"com.apple.security.application-groups",
	}), 0644); err != nil {
		return err, nil
	}
	cmd := exec.Command("codesign", "-s", "Developer ID Application", "-f", "--timestamp", "-i", bundleName, "-o", "runtime", "--entitlements", entitlementsPlist, exePath)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	return cmd.Run(), &buf
}

func entitlements(keys []string) []byte {
	tmpl := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
%s
</dict>
</plist>
`
	for k := range keys {
		keys[k] = fmt.Sprintf("\t<key>%s</key><true/>", keys[k])
	}
	return []byte(fmt.Sprintf(tmpl, strings.Join(keys, "\n")))
}
