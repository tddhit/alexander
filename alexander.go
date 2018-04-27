package alexander

import (
	"bufio"
	"net/http"
	neturl "net/url"
	"os"
)

func Request(rawurl string, data []string, header http.Header, requestNum int, opt Option) error {
	url, err := neturl.Parse(rawurl)
	if err != nil {
		return err
	}
	c := NewClient(opt, url.Host)
	for i := 0; i < requestNum; i++ {
		for j, _ := range data {
			c.Request("POST", url.Path, header, []byte(data[j]))
		}
	}
	return nil
}

func Load(postFile string) ([]string, error) {
	f, err := os.Open(postFile)
	if err != nil {
		return nil, err
	}
	data := make([]string, 0)
	buf := make([]byte, 4096)
	scanner := bufio.NewScanner(f)
	scanner.Buffer(buf, cap(buf))
	for scanner.Scan() {
		line := scanner.Text()
		data = append(data, line)
	}
	return data, nil
}
