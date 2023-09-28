package application

import "github.com/obnahsgnaw/application/pkg/utils"

type Cluster struct {
	id   string
	name string
}

func NewCluster(id, name string) *Cluster {
	return &Cluster{
		id:   id,
		name: name,
	}
}

func (c *Cluster) Id() string {
	return c.id
}

func (c *Cluster) Name() string {
	return c.name
}

func (c *Cluster) String() string {
	return utils.ToStr(c.name, "[", c.id, "]")
}
