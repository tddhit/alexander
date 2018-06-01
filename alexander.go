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
	max     time.Duration
	Min     time.Duration
	avg     time.Duration
	elapsed time.Duration
	success int64
	fail    int64
}

func Summarize(stats []*Stats) {
	var (
		max     time.Duration
		min     time.Duration = MAX_TIME_DURATION
		avg     time.Duration
		elapsed time.Duration
		success int64
		fail    int64
	)
	for _, s := range stats {
		if s.max > max {
			max = s.max
		}
		if s.Min < min {
			min = s.Min
		}
		success += s.success
		fail += s.fail
		elapsed += s.elapsed
	}
	if (success + fail) > 0 {
		avg = elapsed / time.Duration(success+fail)
	}

	log.Info("\t\t\tSummary:")
	log.Info("\t\t\tmax:\t", max)
	log.Info("\t\t\tmin:\t", min)
	log.Info("\t\t\tavg:\t", avg)
	log.Info("\t\t\ttotal:\t", elapsed)
	log.Info("\t\t\tsucc:\t", success)
	log.Info("\t\t\tfail:\t", fail)
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

				if elapsed > stats.max {
					stats.max = elapsed
				}
				if elapsed < stats.Min {
					stats.Min = elapsed
				}
				stats.elapsed += elapsed

				if err != nil {
					log.Error(err, elapsed)
					stats.fail++
				} else {
					log.Info(elapsed)
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
