package main

import (
	"log"
	"net"
	rfc3164 "github.com/influxdata/go-syslog/v3/rfc3164"
	//rfc5424 "github.com/influxdata/go-syslog/v3/rfc5424"
	"gopkg.in/yaml.v2"
	"fmt"
	"io/ioutil"
)

type Hostcfg struct {
	IP string `ip`
	Type string `type`	
}

type Rulecfg struct {
	Type string 
	MatchType string 
	MatchPattern string
	Action string 
}

type YAMLT struct {
	Hosts []Hostcfg `hosts`
	Rules []Rulecfg `rules`
}

const (
	ParserType_Default int 	= 0
	ParserType_Regex     	= 1
	ParserType_Regexi     	= 2
	ParserType_Match     	= 3
	ParserType_Matchi     	= 4
)

const (
	ActionType_Discard int 	= 0
	ActionType_Webhook     	= 1
	ActionType_Alert     	= 2
	ActionType_Retain     	= 3
)

type RULE struct {
	ParserType 	int
	ActionType 	int
	Pattern		string
}

type PROCESSOR struct {
	RulesSlice[] RULE
}

var PRules map[string][]RULE

func read_cfg() {
	var CFG YAMLT
	PRules = make(map[string][]RULE, 0)
	
	cfgdata, err := ioutil.ReadFile("syslog-processor.yaml") // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	
	CFG = YAMLT{}
        err = yaml.Unmarshal([]byte(cfgdata), &CFG)
        if err != nil {
                log.Fatalf("error: %v", err)
        }
        fmt.Printf("--- t:\n%v\n\n", CFG)
	log.Println(CFG)

	

	for _, rule := range CFG.Rules {
		//var ptype int
		//var atype int
		var rulekey string
		// Identify rule namespace (key), device type
		rulekey = rule.Type
		newrule := RULE{} 
		switch {
			case rule.MatchType == "regex":
				newrule.ParserType = ParserType_Regex
			case rule.MatchType == "regexi":
				newrule.ParserType = ParserType_Regexi
			case rule.MatchType == "match":
				newrule.ParserType = ParserType_Match
			case rule.MatchType == "matchi":
				newrule.ParserType = ParserType_Matchi
			default:
				newrule.ParserType = ParserType_Default
		}
		//if (PRules[rulekey] == nil) {
		//	PRules[rulekey] := []RULE{}
		//}
		PRules[rulekey] = append(PRules[rulekey], newrule)
		log.Println(rule)
	}
	fmt.Printf("LOL:%v\n", PRules)
}

func main() {
	read_cfg()
	
	// listen to incoming udp packets
	pc, err := net.ListenPacket("udp", ":1514")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	buf := make([]byte, 65536)
	for {
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}
		go process_rfc3164(pc, addr, buf[:n])
	}
}

func process_rfc3164(pc net.PacketConn, addr net.Addr, buf []byte) {
	p := rfc3164.NewParser()
	m, e := p.Parse(buf)
	if (e != nil) {
		return
	} else {
		log.Println(m)
	}
}


func fetchValue(value interface{}) {
    switch value.(type) {
    case string:
        fmt.Printf("%v is an interface \n ", value)
    case bool:
        fmt.Printf("%v is bool \n ", value)
    case float64:
        fmt.Printf("%v is float64 \n ", value)
    case []interface{}:
        fmt.Printf("%v is a slice of interface \n ", value)
        for _, v := range value.([]interface{}) { // use type assertion to loop over []interface{}
            fetchValue(v)
        }
    case map[string]interface{}:
        fmt.Printf("%v is a map \n ", value)
        for _, v := range value.(map[string]interface{}) { // use type assertion to loop over map[string]interface{}
            fetchValue(v)
        }
    default:
        fmt.Printf("%v is unknown \n ", value)
    }
}
