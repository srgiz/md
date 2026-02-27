package kernel

type di struct {
	services map[string]any
}

func (c *di) Service(name string) any {
	return c.services[name]
}

func (c *di) AddService(name string, service any) {
	c.services[name] = service
}
