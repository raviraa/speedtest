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
	bufKb        = 50    // http buffer size in KB
	maxKb        = 400   // stop speedtest after downloading maxKb KB
	fullDownload = false // ignore maxKB and download full url
)

func main() {
	flag.BoolVar(&fullDownload, "f", false, "download url completely instead of stopping at maximum size(-m) flag")
	flag.IntVar(&maxKb, "m", 400, "maximum size in KB to download, not used if -f flag is set")
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
	buf       []byte
	r         io.Reader
	iterNum   int
	startTime time.Time
	prevTime  time.Time
	// speeds
	avgSpd, maxSpd float64
}

func newDownloader(r io.Reader) *downloader {
	now := time.Now()
	return &downloader{
		buf:       make([]byte, 1024*bufKb),
		r:         r,
		startTime: now,
		prevTime:  now,
	}
}

func (d *downloader) downSpeed() {
	for {
		n, err := io.ReadFull(d.r, d.buf)
		_ = n
		d.iterNum++
		fmt.Printf("%s       ", d.speedstr(true))
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				break
			}
			log.Fatal(err)
		}
		if !fullDownload && d.iterNum*bufKb >= maxKb {
			break
		}
	}
	fmt.Printf("%s in %v\n", d.speedstr(false), time.Since(d.startTime))
}

func (d *downloader) speeds() {
	elapsed := time.Since(d.startTime).Seconds()
	now := time.Now()
	d.avgSpd = float64(d.iterNum*bufKb) / elapsed // in KB/s
	currSpd := float64(bufKb) / now.Sub(d.prevTime).Seconds()
	if currSpd > d.maxSpd {
		d.maxSpd = currSpd
	}
	d.prevTime = now
}

func (d *downloader) speedstr(notFinalRun bool) string {
	if notFinalRun {
		d.speeds()
	}
	return fmt.Sprintf("\rGot %s, A: %s/s, M: %s/s",
		kbOrMb(float64(d.iterNum*bufKb)),
		kbOrMb(d.avgSpd),
		kbOrMb(d.maxSpd),
	)
}

// converts bytes in KB to KB or MB as string.
// 2048 -> 2.00MB
// 100 -> 100.00KB
func kbOrMb(b float64) string {
	if b > 1024 {
		return fmt.Sprintf("%.2fMB", b/1024)
	}
	return fmt.Sprintf("%.2fKB", b)
}
