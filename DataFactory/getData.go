package datafactory

import (
	"io"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/proxy"
)

func GetData(u string, p bool) (s []byte, err error) {
	client := &http.Client{
		Timeout: time.Duration(15 * time.Second),
	}
	if p {
		dialer, err := proxy.SOCKS5("tcp", "localhost:7891", nil, proxy.Direct)
		if err != nil {
			return []byte(""), err
		}
		client = &http.Client{
			Timeout:   time.Duration(15 * time.Second),
			Transport: &http.Transport{Dial: dialer.Dial},
		}
	}
	reqest, err := http.NewRequest("GET", u, nil)

	reqest.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	reqest.Header.Set("Content-Type", "application/json")
	reqest.Header.Set("X-Requested-With", "XMLHttpRequest")
	reqest.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	if err != nil {
		return []byte(""), err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return []byte(""), err
	}
	defer response.Body.Close()
	d, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}
	return d, nil
}

// SyncGetData get data
func SyncGetData(list []string, porxy bool) (d []string) {
	var wg sync.WaitGroup
	d = []string{}
	for i := 0; i < len(list); i++ {
		i0 := i
		wg.Add(1)
		go func() {
			res, errers := GetData(list[i0], porxy)
			if errers != nil {
				d = append(d, "")
			}
			d = append(d, string(res))
			wg.Done()
		}()
	}
	wg.Wait()
	return
}
