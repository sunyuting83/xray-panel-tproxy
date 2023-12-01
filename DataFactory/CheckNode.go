package datafactory

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	config "xpanel/Config"

	"github.com/olahol/melody"
	"golang.org/x/net/proxy"
)

var list = []string{"https://www.google.com.hk",
	"https://baidu.com",
	"https://www.youtube.com/img/desktop/yt_1200.png",
	"https://github.com/webgl-globe/data/data.json"}

func TestData(u string) bool {
	dialer, err := proxy.SOCKS5("tcp", "localhost:7891", nil, proxy.Direct)
	if err != nil {
		return false
	}
	client := &http.Client{
		Timeout:   time.Duration(15 * time.Second),
		Transport: &http.Transport{Dial: dialer.Dial},
	}
	reqest, err := http.NewRequest("GET", u, nil)

	reqest.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	reqest.Header.Set("Content-Type", "application/json")
	reqest.Header.Set("X-Requested-With", "XMLHttpRequest")
	reqest.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.114 Safari/537.36")

	if err != nil {
		return false
	}
	response, err := client.Do(reqest)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return true
}

// SyncCheckData
func SyncCheckData(m *melody.Melody, ID string) {
	var wg sync.WaitGroup
	for i := 0; i < len(list); i++ {
		i0 := i
		wg.Add(1)
		go func() {
			speed := "0"
			start := time.Now()
			res := TestData(list[i0])
			timeElapsed := time.Since(start)
			TimeElapsedStr := Float64ToStringWithPrecision(timeElapsed.Seconds(), 2)
			if res {
				speed = TimeElapsedStr
			}
			speeData := &config.Message{
				Type: "testspeed",
				UUID: ID,
				Data: strings.Join([]string{list[i0], speed}, "||||"),
			}
			sedData, _ := json.Marshal(speeData)
			m.Broadcast(sedData)
			wg.Done()
		}()
	}
	wg.Wait()
	return
}

func Float64ToStringWithPrecision(value float64, precision int) string {
	return strconv.FormatFloat(value, 'f', precision, 64)
}
