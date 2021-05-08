package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	bufKb        = 25    // http buffer size in KB
	maxKb        = 200   // stop speedtest after downloading maxKb KB
	fullDownload = false // ignore maxKB and download full url
)

func main() {
	log.SetFlags(log.Lshortfile)
	flag.BoolVar(&fullDownload, "f", false, "download url completely instead of stopping at maximum size(-m) flag")
	flag.IntVar(&maxKb, "m", 200, "maximum size in KB to download, not used if -f flag is set")
	url := flag.String("u", "http://speedtest-blr1.digitalocean.com/10mb.test", "url to download")
	flag.Parse()
	fmt.Println("Using url: ", *url)
	start := time.Now()

	resp, err := http.Get(*url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Invalid response ", resp.Status)
	}
	fmt.Printf("Connected to %s in %v\n", resp.Request.Host, time.Since(start))

	d := newDownloader(resp.Body)
	d.downSpeed()
}

type downloader struct {
	buf   []byte
	r     io.Reader
	times int
	start time.Time
	prev  time.Time
	// speeds
	avgSpd, maxSpd float64
}

func newDownloader(r io.Reader) *downloader {
	now := time.Now()
	return &downloader{
		buf:   make([]byte, 1024*bufKb),
		r:     r,
		start: now,
		prev:  now,
	}
}

func (d *downloader) downSpeed() {
	for {
		n, err := io.ReadFull(d.r, d.buf)
		_ = n
		d.times++
		fmt.Printf("%s       ", d.speedstr(true))
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			log.Fatal(err)
		}
		if !fullDownload && d.times*bufKb >= maxKb {
			break
		}
	}
	fmt.Printf("%s in %v\n", d.speedstr(false), time.Since(d.start))
}

func (d *downloader) speeds() {
	elapsed := time.Since(d.start).Seconds()
	now := time.Now()
	d.avgSpd = float64(d.times*bufKb) / elapsed // in KB/s
	curr := float64(bufKb) / now.Sub(d.prev).Seconds()
	if curr > d.maxSpd {
		d.maxSpd = curr
	}
	d.prev = now
}

func (d *downloader) speedstr(calc bool) string {
	if calc {
		d.speeds()
	}
	return fmt.Sprintf("\rGot %s, A: %s/s, M: %s/s",
		kbMb(float64(d.times*bufKb)),
		kbMb(d.avgSpd),
		kbMb(d.maxSpd),
	)
}

// converts bytes in KB to KB or MB as string.
// 2048 -> 2.00MB
// 100 -> 100.00KB
func kbMb(b float64) string {
	if b > 1024 {
		return fmt.Sprintf("%.2fMB", b/1024)
	}
	return fmt.Sprintf("%.2fKB", b)
}
