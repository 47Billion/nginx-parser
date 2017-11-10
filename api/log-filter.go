package nginx_parser

import (
	"os"
	"bufio"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"path/filepath"
	"sync"
	"reflect"

	"github.com/satyrius/gonx"
	"github.com/pkg/errors"
	"github.com/apex/log"
)

// Known Args
/*
	SubID
	a
	ab
	adwords
	aid
	aif1
	aif5
	aifa
	an
	and1
	andi
	api_key
	append_app_conv_trk_params
	av
	br
	c
	ca
	campaignid
	cid
	city
	citycode
	cl
	click_t
	click_ts
	country
	cr
	creative
	d
	ddl_enabled
	ddl_to
	de
	did
	dk
	dn
	dnt
	e
	event_type
	gclid
	goal_AC
	goal_FP
	goal_ac
	goal_fp
	goal_id
	h
	hop_cnt
	hops
	i
	id
	idfa
	idfv ifa1 ifa5 install_ts ip is_view_thru k keyword lag loc_physical_ms lpurl ma mac1 matchtype mm_uuid mo n network o odin op p pid pl pm pn pr pref product_id r rand re referrer region rt s sc scs sdk seq siteid src st strategy su subid t tt_adv_id u udi1 udi5 udid ut utm_content v vca vcr ve voucherKey voucherkey]

 */


const (
	nginxLogFormat = `$remote_addr - - [$time_local] "$request_method $request_query $request_protocol"`
	nginxLocalTime = "01/Jan/2017:00:59:58 -0800"
	dateRFC1123Z = "Mon, 02 Jan 2006 15:04:05 -0700"
)

func ReadFile(file string, records chan *Record) (err error) {
	fmt.Printf("%s \n", file)
	// File Reader
	fileReader, err := os.Open(file)
	if err != nil {
		err = errors.Wrapf(err, "Failed to read file=%s", file)
	}

	// Nginx
	reader := gonx.NewReader(fileReader, nginxLogFormat)
	// Loop
	for {
		rec, err := reader.Read()
		if err != nil {
			// End
			if err == io.EOF {
				log.Debugf("Done! file=%s", file)
				break;
			}
			log.Errorf("Error in reading; err=%+v", err)
		}


		// Record
		record := &Record{}

		// IP
		remoteAddress, err := rec.Field("remote_addr")
		if err == nil {
			record.RemoteAddress = remoteAddress
		}

		// Date
		date, err := rec.Field("time_local")
		if err == nil {
			record.TimeLocal = date
		}

		// Request
		reqQuery, err := rec.Field("request_query")
		if err == nil {
			path, args, err := ParseReq(reqQuery)
			if err == nil {
				record.Path = path
				record.Args = args
			}
		}

		// To Channel
		records <- record
	}

	return
}

func ParseReq(req string) (path string, args map[string][]string, err error) {
	// Parse
	u, err := url.Parse(req)
	if err != nil {
		err = errors.Wrapf(err, "Failed to parse req=%s", req)
		return
	}

	path = u.Path
	args = u.Query()

	return
}

func Start() {
	dir := "testData"
	files := make(chan string)
	records := make(chan *Record)
	data := make(chan *Entry)

	// Find Files
	go func() {
		FindLogFiles(dir, files)
		close(files)
	}()

	go func() {
		done := ExecParallelFn(5, "reader", func() {
			for log := range files {
				err := ReadFile(log, records)
				if err != nil {
					panic(err)
				}
			}
		})
		<-done
		close(records)
	}()

	e := Entry{}
	go func() {
		for r := range records {
			t := reflect.TypeOf(e)
			reflect.ValueOf(&e).Elem().FieldByName("Path").SetString(r.Path)
			reflect.ValueOf(&e).Elem().FieldByName("RemoteAddress").SetString(r.RemoteAddress)
			reflect.ValueOf(&e).Elem().FieldByName("TimeLocal").SetString(r.TimeLocal)

			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				if val, found := r.Args[strings.ToLower(field.Name)]; found {
					valStr := strings.Join(val, ",")
					reflect.ValueOf(&e).Elem().FieldByName(field.Name).SetString(valStr)
				}
			}
			data <- &e
		}
		close(data)

	}()

	LoopRecords(data)
}

func FindLogFiles(dir string, files chan string) {
	ext := ".log"
	filepath.Walk(dir, func(file string, fileInfo os.FileInfo, err error) error {
		if nil != err {
			return io.EOF
		}

		if fileInfo.IsDir() || filepath.Ext(file) != ext {
			return nil
		}

		// To Channel
		files <- file
		return nil
	})
}

func Sort(in map[string][]string) string {
	keys := make([]string, 0, len(in))
	for key := range in {
		keys = append(keys, key)
	}

	// sort
	sort.Strings(keys) //sort by key

	var out []string
	for _, key := range keys {
		out = append(out, fmt.Sprintf("%s:%+v", key, in[key][0]))
	}

	return strings.Join(out, ", ")
}

func readLine(file string, ch chan string) {
	in, _ := os.Open(file)
	defer in.Close()

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ch <- scanner.Text()
	}
}

// Runs func `workFn` parallely in goroutines
func ExecParallelFn(numWorkers int, name string, workFn func()) <-chan bool {
	done := make(chan bool) // done channel

	// Work - runs func `workFn` parallely in goroutines.
	// Waits for all workFn to finish
	// Mark end to done
	go func() {
		_logger := log.WithFields(log.Fields{// Logger
			"file": "exec-parallel.go",
			"fn": name,
		})
		_logger.Debugf("Preparing [%d] workers", numWorkers)

		wg := &sync.WaitGroup{}
		wg.Add(numWorkers)

		for i := 1; i <= numWorkers; i++ {
			_logger.Debugf("%d/%s: workFn() - start", i, name)

			go func(id int) {
				// work
				workFn()

				// mark done
				wg.Done()

				_logger.Debugf("%d/%s: workFn() - finish", id, name)
			}(i)
		}

		// wait
		wg.Wait()

		_logger.Debug("Finish()")

		// Mark done
		done <- true
	}()

	return done
}

func PrintForGoStruct(keys []string) {
	fmt.Printf("\n")
	fmt.Printf("type Entry struct { \n")
	for _, k := range keys {
		fmt.Printf("\t%s string `csv:\"%s\"` \n", k, k)
	}
	fmt.Printf("} \n")
}

func (e *Entry) From(r *Record) {
}
