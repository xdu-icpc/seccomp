package run

func (c *Cmd) kill() error {
	return c.Process.Kill()
}
