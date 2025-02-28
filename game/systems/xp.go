package systems

import "github.com/yohamta/donburi/ecs"

type xpsystem struct{}

var XPSystem = &xpsystem{}

func (x *xpsystem) Update(e *ecs.ECS) {
	
}