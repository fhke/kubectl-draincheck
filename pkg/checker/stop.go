package checker

func (c *Checker) Stop() {
	c.pdbLocator.Stop()
}
