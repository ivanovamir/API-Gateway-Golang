package config

type Option func(c *config)

func WithPath(path string) Option {
	return func(c *config) {
		c.Path = path
	}
}
