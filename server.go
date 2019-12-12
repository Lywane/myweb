package myhttp

func RecoveryHandler(c *Context) {
	defer func() {
		if err := recover(); err != nil {
			c.DieWithHttpStatus(500)
		}
	}()
	c.Next()
}