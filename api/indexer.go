package nginx_parser

func Index(rs []*Entry) error {
	var err error
	// index into ES and SQL in parallel
	err = ExecMultipleParallelFnWithError(
		func() error {
			return SqlInsertRecords(rs)
		},
	)
	return err
}
