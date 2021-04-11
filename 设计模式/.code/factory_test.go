package factory

import "testing"

func TestFactory(t *testing.T) {
	cpf := &ConfigParserFactory{}
	xcp := cpf.createParser("xml")
	xcp.parse()
	jcp := cpf.createParser("json")
	jcp.parse()
}
