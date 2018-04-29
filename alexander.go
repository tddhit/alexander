package alexander

import (
	"bufio"
	"context"
	"net/http"
	neturl "net/url"
	"os"
	"time"

	"github.com/tddhit/tools/log"
)

const (
	MAX_TIME_DURATION = 1 * time.Hour
)

type Stats struct {
	max     int64
	Min     int64
	avg     int64
	success int64
	fail    int64
	elapsed int64
}

func Summarize(stats []*Stats) {
	var (
		max     int64
		min     int64 = int64(MAX_TIME_DURATION)
		avg     int64
		success int64
		fail    int64
		elapsed int64
	)
	for _, s := range stats {
		if s.max > max {
			max = s.max
		} else if s.Min < min {
			min = s.Min
		}
		success += s.success
		fail += s.fail
		elapsed += s.elapsed
	}
	avg = elapsed / success

	log.Info("\t\t\tSummary:")
	log.Infof("\t\t\tmax:\t%dms", max/1000000)
	log.Infof("\t\t\tmin:\t%dms", min/1000000)
	log.Infof("\t\t\tavg:\t%dms", avg/1000000)
	log.Infof("\t\t\tfail:\t%d", fail)
	log.Infof("\t\t\tsucc:\t%d", success)
}

func Request(ctx context.Context, rawurl string, data []string,
	header http.Header, repeatTimes int, opt Option, stats *Stats) error {

	url, err := neturl.Parse(rawurl)
	if err != nil {
		return err
	}
	c := NewClient(opt, url.Host)
	for i := 0; i < repeatTimes; i++ {
		for j, _ := range data {
			select {
			case <-ctx.Done():
				goto exit
			default:
				start := time.Now()
				_, err := c.Request("POST", url.Path, header, []byte(data[j]))
				end := time.Now()
				elapsed := end.Sub(start)

				if int64(elapsed) > stats.max {
					stats.max = int64(elapsed)
				} else if int64(elapsed) < stats.Min {
					stats.Min = int64(elapsed)
				}

				if err != nil {
					log.Error(err, elapsed)
					stats.fail++
				} else {
					log.Info(elapsed)
					stats.elapsed += int64(elapsed)
					stats.success++
				}
			}
		}
	}
exit:
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
