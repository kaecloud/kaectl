package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/glamour"
	"github.com/kaecloud/kaectl/internal/run"
	"github.com/kaecloud/kaectl/pkg/browser"
)

// OpenInBrowser opens the url in a web browser based on OS and $BROWSER environment variable
func OpenInBrowser(url string) error {
	browseCmd, err := browser.Command(url)
	if err != nil {
		return err
	}
	return run.PrepareCmd(browseCmd).Run()
}

func RenderMarkdown(text string) (string, error) {
	// Glamour rendering preserves carriage return characters in code blocks, but
	// we need to ensure that no such characters are present in the output.
	text = strings.ReplaceAll(text, "\r\n", "\n")

	renderStyle := glamour.WithStandardStyle("notty")
	// TODO: make color an input parameter
	if isColorEnabled() {
		renderStyle = glamour.WithEnvironmentConfig()
	}

	tr, err := glamour.NewTermRenderer(
		renderStyle,
		// glamour.WithBaseURL(""),  // TODO: make configurable
		// glamour.WithWordWrap(80), // TODO: make configurable
	)
	if err != nil {
		return "", err
	}

	return tr.Render(text)
}

func Pluralize(num int, thing string) string {
	if num == 1 {
		return fmt.Sprintf("%d %s", num, thing)
	}
	return fmt.Sprintf("%d %ss", num, thing)
}

func fmtDuration(amount int, unit string) string {
	return fmt.Sprintf("about %s ago", Pluralize(amount, unit))
}

func FuzzyAgo(ago time.Duration) string {
	if ago < time.Minute {
		return "less than a minute ago"
	}
	if ago < time.Hour {
		return fmtDuration(int(ago.Minutes()), "minute")
	}
	if ago < 24*time.Hour {
		return fmtDuration(int(ago.Hours()), "hour")
	}
	if ago < 30*24*time.Hour {
		return fmtDuration(int(ago.Hours())/24, "day")
	}
	if ago < 365*24*time.Hour {
		return fmtDuration(int(ago.Hours())/24/30, "month")
	}

	return fmtDuration(int(ago.Hours()/24/365), "year")
}

func Humanize(s string) string {
	// Replaces - and _ with spaces.
	replace := "_-"
	h := func(r rune) rune {
		if strings.ContainsRune(replace, r) {
			return ' '
		}
		return r
	}

	return strings.Map(h, s)
}

// We do this so we can stub out the spinner in tests -- it made things really flakey. this is not
// an elegant solution.
var StartSpinner = func(s *spinner.Spinner) {
	s.Start()
}

var StopSpinner = func(s *spinner.Spinner) {
	s.Stop()
}

func Spinner(w io.Writer) *spinner.Spinner {
	return spinner.New(spinner.CharSets[11], 400*time.Millisecond, spinner.WithWriter(w))
}

func IsURL(s string) bool {
	return strings.HasPrefix(s, "http:/") || strings.HasPrefix(s, "https:/")
}

func DisplayURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	return u.Hostname() + u.Path
}

func GreenCheck() string {
	return Green("✓")
}

// FileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ExpandUser(pathname string) string {
	u, _ := user.Current()
	homedir := u.HomeDir
	if pathname == "~" {
		return homedir
	} else if strings.HasPrefix(pathname, "~/") {
		pathname = filepath.Join(homedir, pathname[2:])
	}
	return pathname
}

func CompressDir(w *zip.Writer, dirName string) error {
	err := filepath.Walk(dirName, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			err := CompressFile(w, path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func CompressFile(w *zip.Writer, filename string) error {
	srcF, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer srcF.Close()

	f, err := w.Create(filename)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, srcF); err != nil {
		return err
	}
	return nil
}

func Compress(zipFilename string, filename string) error {
	outFile, err := os.Create(zipFilename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)
	defer w.Close()

	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		// do directory stuff
		err = CompressDir(w, filename)
	case mode.IsRegular():
		// upload file if SrcName is regular file
		err = CompressFile(w, filename)
	}
	return err
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

