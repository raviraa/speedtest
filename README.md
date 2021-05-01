# speedtest

Simple command line app to test internet speed in console. Can do partial approximate test too downloading small size instead of downloading a big file, if bandwidth is constrained. Prints the connection time and download speed.

```
$./speedtest  -h
Usage of ./speedtest:
  -f	download url completely instead of stopping at maximum size(-m) flag
  -m int
    	maximum size in KB to download, not used if -f flag is set (default 200)
  -u string
    	url to download (default "http://speedtest-blr1.digitalocean.com/10mb.test")

```

```
$./speedtest
Using url:  http://speedtest-blr1.digitalocean.com/10mb.test
Connected to speedtest-blr1.digitalocean.com in 262.927285ms
Downloaded 200.00KB, Speed 633.32KB/s in 315.802636ms
```
