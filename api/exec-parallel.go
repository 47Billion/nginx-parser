package nginx_parser

import (
	"sync"
	"time"
	"github.com/pkg/errors"
	"github.com/apex/log"
)

func ExecMultipleParallelFnWithError(workFns ...func() error) (err error) {
	count := len(workFns)

	wg := sync.WaitGroup{}
	wg.Add(count)

	errs := make(chan error)
	go func() {
		for e := range errs {
			err = errors.Wrap(e, "\n")
		}
	}()

	for _, workFn := range workFns {
		go func(fn func() error) {
			e := fn()
			if e != nil {
				errs <- e
			}
			wg.Done()
		}(workFn)
	}

	wg.Wait()

	close(errs)
	return
}

type Flushable struct {
	M interface{}
}

func ExecLoopAndFlush(flushInterval time.Duration, flushLength int, scan chan *Flushable, flushFn func([]*Flushable) error) (chan bool) {
	done := make(chan bool)
	fn := func() {
		var batch []*Flushable // Batch
		batch = []*Flushable{} // Init

		copy := make(chan *Flushable) // copy of scan; just so that we avoid infinite loop on select
		ticker := time.NewTicker(flushInterval) // ticker - flushInterval
		quit := make(chan bool)

		// Flusher
		flush := func() {
			if len(batch) == 0 {
				return
			}

			if err := flushFn(batch); err == nil {
				log.Debugf("flushFn() err=%+v", err)
				// reset batch
				batch = []*Flushable{}
			}
		}

		// Loop
		loop := func() {
			for {
				select {
				case <-ticker.C:
					flush() // force flush
				case <-quit:
					ticker.Stop() // stop ticker
					flush() // force flush
					close(done) // mark done
					return
				case c := <-copy:
				// read from copy; beware - if copy closes, you will face an infinite loop here!
					batch = append(batch, c) // append
					if len(batch) >= flushLength {
						flush()           // flush only if batch size is reached
					}
				}
			}
		}

		log.Debugf("ExecLoopAndFlush - starting...")

		go loop() // Loop unit quit

		for c := range scan {
			// range
			copy <- c // copy
		}

		log.Debugf("ExecLoopAndFlush - terminating...")
		quit <- true // quit; it mean scan channel is closed
	}

	// Run
	go fn()

	return done
}