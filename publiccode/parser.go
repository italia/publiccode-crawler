package publiccode

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v2"
)

// Parse loads the yaml bytes and tries to parse it. Return an error if fails.
func Parse(in []byte, pc *PublicCode) error {
	var s map[interface{}]interface{}

	d := yaml.NewDecoder(bytes.NewReader(in))
	if err := d.Decode(&s); err != nil {
		return err
	}

	return newParser(pc).parse(s)
}

type parser struct {
	pc      *PublicCode
	missing map[string]bool
}

func newParser(pc *PublicCode) *parser {
	var p parser
	p.pc = pc
	p.missing = make(map[string]bool)
	for _, k := range mandatoryKeys {
		p.missing[k] = true
	}
	return &p
}

func (p *parser) parse(s map[interface{}]interface{}) error {
	if err := p.decoderec("", s); err != nil {
		return err
	}
	if err := p.finalize(); err != nil {
		return err
	}
	return nil
}

func (p *parser) decoderec(prefix string, s map[interface{}]interface{}) (es ErrorParseMulti) {
	for ki, v := range s {
		k, ok := ki.(string)
		if !ok {
			es = append(es, ErrorInvalidKey{Key: fmt.Sprint(ki)})
			continue
		}
		if prefix != "" {
			k = prefix + "/" + k
		}
		delete(p.missing, k)

		switch v := v.(type) {
		case string:
			if err := p.decodeString(k, v); err != nil {
				es = append(es, err)
			}
		case []interface{}:
			sl := []string{}
			sli := make(map[interface{}]interface{})

			for idx, v1 := range v {
				// if array of strings
				if s, ok := v1.(string); ok {
					sl = append(sl, s)
					if len(sl) == len(v) { //the v1.(string) check should be extracted.
						if err := p.decodeArrString(k, sl); err != nil {
							es = append(es, err)
						}
					}
					// if array of objects
				} else if _, ok := v1.(map[interface{}]interface{}); ok {
					sli[k] = v1
					if err := p.decodeArrObj(k, sli); err != nil {
						es = append(es, err)
					}

				} else {
					es = append(es, newErrorInvalidValue(k, "array element %d not a string", idx))
				}
			}
		case map[interface{}]interface{}:

			if errs := p.decoderec(k, v); len(errs) > 0 {
				es = append(es, errs...)
			}
		default:
			if v == nil {
				panic(fmt.Errorf("key \"%s\" is empty. Remove it or fill with valid values.", k))
			}
			panic(fmt.Errorf("key \"%s\" - invalid type: %T", k, v))
		}
	}
	return
}
