package appcache

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Manifest struct {
	CacheItems       []string
	NetworkItems     []string
	FallbackPatterns [][]string
	Comment          string
	buildCache       string
	IgnorePatterns   []string
	Verbose          bool
	ChecksumType     int
	addedFiles       []string
}

const (
	GitRevChecksum = iota
	TimestampChecksum
	FileContentChecksum
	HgIdChecksum
)

const MIMEType = "text/cache-manifest"

func chomp(str string) string {
	return strings.TrimRight(str, "\n\r")
}

func HgId() string {
	var out bytes.Buffer
	cmd := exec.Command("hg", "identify")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return chomp(out.String())
}

func GitParseAbbrRev() string {
	var out bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return chomp(out.String())
}

func GitParseRev() string {
	var out bytes.Buffer
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return chomp(out.String())
}

func (self *Manifest) SetChecksumType(t int) {
	self.ChecksumType = t
}

func (self *Manifest) AddCache(item string) {
	self.CacheItems = append(self.CacheItems, item)
}

func (self *Manifest) AddIgnorePattern(item string) {
	self.IgnorePatterns = append(self.IgnorePatterns, item)
}

func (self *Manifest) AddNetwork(item string) {
	self.NetworkItems = append(self.NetworkItems, item)
}

func (self *Manifest) AddFallback(pattern string, target string) {
	self.FallbackPatterns = append(self.FallbackPatterns, []string{pattern, target})
}

func (self *Manifest) SetComment(comment string) {
	self.Comment = comment
}

func (self *Manifest) AddComment(comment string) {
	self.Comment += comment
}

func (self *Manifest) AddCacheFromDirectory(root string, publicRoot string, prefix string) error {
	var err = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		for _, pattern := range self.IgnorePatterns {
			if matched, _ := regexp.MatchString(pattern, p); matched {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
				// return nil
			}
			/*
				if strings.Contains(p, pattern) {
					return nil
				}
			*/
			/*
				matched, err := path.Match(pattern, p)
				if err != nil {
					panic(err)
					return nil
				}
				if matched {
					return nil
				}
			*/
		}

		var itemPath = p[len(publicRoot):]
		if self.Verbose {
			log.Println("Adding cache item: ", itemPath)
		}
		self.AddCache(prefix + itemPath)
		self.addedFiles = append(self.addedFiles, p)
		return nil
	})
	if err != nil {
		return err
	}
	return err
}

func (self *Manifest) BuildTimestampChecksum() string {
	var now = time.Now()
	var unix = now.Unix()
	return fmt.Sprintf("%d", unix) + ", " + now.Format(time.RFC3339Nano)
}

func (self *Manifest) BuildChecksum() string {
	if self.ChecksumType != 0 {
		if self.ChecksumType == TimestampChecksum {
			return self.BuildTimestampChecksum()
		} else if self.ChecksumType == GitRevChecksum {
			return "git:" + GitParseAbbrRev() + ":" + GitParseRev()
		} else if self.ChecksumType == HgIdChecksum {
			return "hg:" + HgId()
		} else if self.ChecksumType == FileContentChecksum {
			h := md5.New()
			for _, file := range self.addedFiles {
				contents, err := ioutil.ReadFile(file)
				if err != nil {
					log.Println(err)
				}
				h.Write([]byte(contents))
			}
			return fmt.Sprintf("md5:%x", h.Sum(nil))
		}
	}
	return self.BuildTimestampChecksum()
}

func (self *Manifest) CacheString() string {
	if self.buildCache != "" {
		return self.buildCache
	}
	self.buildCache = self.String()
	return self.buildCache
}

func (self *Manifest) String() string {
	var output string = "CACHE MANIFEST"
	output += "\n# " + self.BuildChecksum()

	if len(self.CacheItems) > 0 {
		output += "\nCACHE:\n" + strings.Join(self.CacheItems, "\n")
	}
	if len(self.NetworkItems) > 0 {
		output += "\nNETWORK:\n" + strings.Join(self.NetworkItems, "\n")
	}

	if len(self.FallbackPatterns) > 0 {
		output += "\nFALLBACK:"
		for _, item := range self.FallbackPatterns {
			output += "\n" + item[0] + " " + item[1]
		}
	}
	if self.Comment != "" {
		output += "\n# " + strings.Replace(self.Comment, "\n", "\n# ", 0)
	}
	return output
}

func (self *Manifest) Write(w io.Writer) {
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.Header().Set("Content-Type", MIMEType)
	}
	fmt.Fprint(w, self.CacheString())
}

func NewManifest() *Manifest {
	return &Manifest{}
}
