package appcache

import "testing"

func TestManifest(t *testing.T) {
	var manifest = NewManifest()
	manifest.AddCacheFromDirectory("public/js", "public")
	manifest.AddFallback("/", "/offline")
	manifest.SetComment("Author")
	t.Log(manifest.String())
}
