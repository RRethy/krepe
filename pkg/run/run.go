package run

func Run(pkg, pipeline, function string) error {
	r, err := newRunnable(pkg, pipeline, function)
	if err != nil {
		return err
	}
	return r.run()
}
