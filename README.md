Appcache: Manifest Generator
=============================

Supported checksum types:

- GitRevChecksum
- HgIdChecksum
- FileContentChecksum
- TimestampChecksum


Example
--------

```go
import "github.com/c9s/appcache"
var gManifest = appcache.NewManifest()

func init() {
	if Environment == PRODUCTION {
		gManifest.SetChecksumType(appcache.GitRevChecksum)

		gManifest.AddIgnorePattern("\\.html$")
		gManifest.AddIgnorePattern("\\.htm$")
		gManifest.AddIgnorePattern("\\.sass-cache")
		gManifest.AddIgnorePattern(".less$")
		gManifest.AddIgnorePattern(".sass$")
		gManifest.AddIgnorePattern(".scss$")

		gManifest.AddCache("/public/css/font-awesome/font/fontawesome-webfont.woff?v=3.2.0")
		gManifest.AddCache("/public/css/font-awesome/font/fontawesome-webfont.ttf?v=3.2.0")

		// gManifest.AddCacheFromDirectory("public/src", "public", "/public")
		gManifest.AddCacheFromDirectory("public/css/main.css", "public", "/public")
		gManifest.AddCacheFromDirectory("public/css/main.min.css", "public", "/public")
		gManifest.AddCacheFromDirectory("public/css/font-awesome/css/font-awesome.min.css", "public", "/public")
		gManifest.AddCacheFromDirectory("public/built", "public", "/public")

		gManifest.AddNetwork("*")

		// manifest.AddFallback("/", "/offline")
		gManifest.SetComment("comment here")
	}
}

func manifestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/cache-manifest")
	fmt.Fprint(w, gManifest.CacheString())
}
```
