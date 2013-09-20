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
var manifest = appcache.NewManifest()

func init() {
	if Environment == PRODUCTION {
		manifest.SetChecksumType(appcache.GitRevChecksum)

		manifest.AddIgnorePattern("\\.html$")
		manifest.AddIgnorePattern("\\.htm$")
		manifest.AddIgnorePattern("\\.sass-cache")
		manifest.AddIgnorePattern(".less$")
		manifest.AddIgnorePattern(".sass$")
		manifest.AddIgnorePattern(".scss$")

		manifest.AddCache("/public/css/font-awesome/font/fontawesome-webfont.woff?v=3.2.0")
		manifest.AddCache("/public/css/font-awesome/font/fontawesome-webfont.ttf?v=3.2.0")

        // AddCacheFromFile([local file path], [local base path], [base url])
		manifest.AddCacheFromFile("local/path/to/public/css/main.css", "local/path/to/public", "/public")
		manifest.AddCacheFromFile("local/path/to/public/css/main.min.css", "local/path/to/public", "/public")
		manifest.AddCacheFromFile("local/path/to/public/css/font-awesome/css/font-awesome.min.css", "public", "/public")
		manifest.AddCacheFromDirectory("app/built", "app", "/public/app")

		manifest.AddNetwork("*")

		// manifest.AddFallback("/", "/offline")
		manifest.SetComment("comment here")
	}
}

func manifestHandler(w http.ResponseWriter, r *http.Request) {
    manifest.Write(w)

    // same as below
	w.Header().Set("Content-Type", "text/cache-manifest")
	fmt.Fprint(w, manifest.CacheString())
}
```
