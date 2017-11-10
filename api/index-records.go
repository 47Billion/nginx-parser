package nginx_parser

import (
	"sync"
	"time"
	"github.com/apex/log"
)

var (
	parallelRecordsLooper = 5

	flushInterval = 5 * time.Second
	flushLength = 100

	_indexerLogger = log.WithFields(log.Fields{// Logger
		"name": "mysql",
		"file": "mysql_queries.go",
	})
)

func LoopRecords(records chan *Entry) {
	_indexerLogger.Info("LoopRecords")
	// Flusher
	flush := make(chan *Flushable)
	// Run periodic flush to MySQL & ES
	done := ExecLoopAndFlush(flushInterval, flushLength, flush, func(arr []*Flushable) (err error) {
		// Make a slice of Audit objects
		recordsToFlush := make([]*Entry, 0) // length=0
		for _, f := range arr {
			r := f.M.(*Entry)
			recordsToFlush = append(recordsToFlush, r)
		}
		if err := Index(recordsToFlush); err != nil {
			_indexerLogger.Errorf("Failed to index records; count=%d", len(recordsToFlush))
		}
		_indexerLogger.Debugf("indexer Flush => %d", len(recordsToFlush))
		return
	})

	// Loop records
	// - Persist JSON
	// - Push to flushable
	loop := func() {
		for r := range records {
			// To Flush
			f := &Flushable{}
			f.M = r // Record
			_indexerLogger.Debugf("indexer loop - recordId=%d", r.Id)
			flush <- f
		}

	}

	// Run multiple loopers
	wg := sync.WaitGroup{}
	wg.Add(parallelRecordsLooper)
	for i := 0; i < parallelRecordsLooper; i++ {
		go func() {
			loop()
			wg.Done()
		}()
	}
	wg.Wait() // Wait
	// close
	close(flush)
	<-done // Waiting for flusher

}

