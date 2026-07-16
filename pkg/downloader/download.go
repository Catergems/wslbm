package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Download fetches url to destDir/<filename> and returns the local path.
func Download(url, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	filename := filepath.Base(url)
	dest := filepath.Join(destDir, filename)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
	}

	total, _ := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)

	f, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if total > 0 {
		fmt.Printf("Downloading %s (%.1f MB)\n", filename, float64(total)/1e6)
	} else {
		fmt.Printf("Downloading %s\n", filename)
	}

	pr := &progressReader{r: resp.Body, total: total, start: time.Now()}
	if _, err := io.Copy(f, pr); err != nil {
		return "", err
	}

	// Final line
	elapsed := time.Since(pr.start).Seconds()
	fmt.Printf("\r[%s] 100%% / %.1fs                    \n", filled(32), elapsed)

	return dest, nil
}

const barWidth = 32

func filled(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = '#'
	}
	return string(s)
}

type progressReader struct {
	r       io.Reader
	total   int64
	written int64
	start   time.Time
}

func (p *progressReader) Read(buf []byte) (int, error) {
	n, err := p.r.Read(buf)
	p.written += int64(n)

	pct := 0.0
	if p.total > 0 {
		pct = float64(p.written) / float64(p.total)
	}

	filled := int(pct * barWidth)
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "#"
	}
	for i := 0; i < empty; i++ {
		bar += " "
	}

	elapsed := time.Since(p.start).Seconds()
	var est string
	if pct > 0.01 && pct < 1.0 {
		remaining := (elapsed / pct) - elapsed
		est = fmt.Sprintf("EST: %.0fs", remaining)
	} else {
		est = "EST: --"
	}

	fmt.Printf("\r[%s] %.0f%% / %s   ", bar, pct*100, est)

	return n, err
}
