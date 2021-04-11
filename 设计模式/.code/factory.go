package factory

import "fmt"

type ConfigParser interface {
	parse()
}

type XmlConfigParser struct {
}

func (xcp XmlConfigParser) parse() {
	fmt.Println("This is a xml parser.")
}

type JsonConfigParser struct {
}

func (jcp JsonConfigParser) parse() {
	fmt.Println("This is a json parser.")
}

type YamlConfigParser struct {
}

func (ycp YamlConfigParser) parse() {
	fmt.Println("This is a yaml parser.")
}

// 工厂类
type ConfigParserFactory struct {
}

// 工厂方法
func (cpf ConfigParserFactory) createParser(format string) ConfigParser {
	switch format {
	case "xml":
		return XmlConfigParser{}
	case "json":
		return JsonConfigParser{}
	case "yaml":
		return YamlConfigParser{}
	default:
		return nil
	}
}
