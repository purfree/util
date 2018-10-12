package net

type Route struct {
	handlers map[string]func(cmd string, data []byte, session *Session)
}

func (p *Route) Register(cmd string, handler func(cmd string, data []byte, session *Session)) {
	p.handlers[cmd] = handler
}
