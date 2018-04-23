package publiccode

var mandatoryKeys = []string{
	"version",
	"url",
	"legal/license",
	"legal/repo-owner",
	"maintenance/type",
	"description/name",
	"description/platforms",
	"description/shortdesc",
	"description/longdesc",
}

func (p *parser) decodeString(key string, value string) (err error) {
	switch key {
	case "version":
		p.pc.Version = value
		if p.pc.Version != Version {
			return newErrorInvalidValue(key, "version %s not supported", p.pc.Version)
		}
	case "url":
		p.pc.Url, err = p.checkUrl(key, value)
	case "upstream-url":
		return p.decodeArrString(key, []string{value})
	case "legal/license":
		p.pc.Legal.License = value
		return p.checkSpdx(key, value)
	case "legal/main-copyright-owner":
		p.pc.Legal.MainCopyrightOwner = value
	case "legal/authors-file":
		p.pc.Legal.AuthorsFile, err = p.checkFile(key, value)
	case "legal/repo-owner":
		p.pc.Legal.RepoOwner = value
	case "maintenance/type":
		for _, v := range []string{"community", "commercial", "none"} {
			if v == value {
				p.pc.Maintenance.Type = value
				return nil
			}
		}
		return newErrorInvalidValue(key, "invalid value: %s", value)
	case "maintenance/until":
		p.pc.Maintenance.Until, err = p.checkDate(key, value)
		return err
	case "maintenance/maintainer":
		return p.decodeArrString(key, []string{value})
	case "description/name":
		p.pc.Description.Name = value
	case "description/logo":
		return p.decodeArrString(key, []string{value})
	case "description/version":
		p.pc.Description.Version = value
	case "description/platforms":
		p.pc.Description.Platforms = value
	case "description/released":
		p.pc.Description.Released, err = p.checkDate(key, value)
	case "meta/scope":
		return p.decodeArrString(key, []string{value})
	case "meta/pa-type":
		return p.decodeArrString(key, []string{value})
	case "meta/category":
		p.pc.Meta.Category = value
	case "meta/tags":
		return p.decodeArrString(key, []string{value})
	case "meta/used-by":
		return p.decodeArrString(key, []string{value})
	case "dependencies/open":
		return p.decodeArrString(key, []string{value})
	case "dependencies/proprietary":
		return p.decodeArrString(key, []string{value})
	case "dependencies/hardware":
		return p.decodeArrString(key, []string{value})
	default:
		return ErrorInvalidKey{key + " : String"}
	}
	return
}

func (p *parser) decodeArrString(key string, value []string) error {
	switch key {
	case "upstream-url":
		for _, v := range value {
			if u, err := p.checkUrl(key, v); err != nil {
				return err
			} else {
				p.pc.UpstreamUrl = append(p.pc.UpstreamUrl, u)
			}
		}
	case "maintenance/maintainer":
		p.pc.Maintenance.Maintainer = value
	case "description/logo":
		for _, v := range value {
			if err := p.checkImage(key, v); err != nil {
				return err
			} else {
				p.pc.Description.Logo = append(p.pc.Description.Logo, v)
			}
		}
	case "description/screenshots":
		for _, v := range value {
			if f, err := p.checkFile(key, v); err != nil {
				return err
			} else {
				p.pc.Description.Screenshots = append(p.pc.Description.Screenshots, f)
			}
		}
	case "description/videos":
		for _, v := range value {
			if u, err := p.checkUrl(key, v); err != nil {
				return err
			} else {
				p.pc.Description.Videos = append(p.pc.Description.Videos, u)
			}
		}
	case "meta/scope":
		for _, v := range value {
			p.pc.Meta.Scope = append(p.pc.Meta.Scope, v)
		}
	case "meta/pa-type":
		for _, v := range value {
			if u, err := p.checkPaTypes(key, v); err != nil {
				return err
			} else {
				p.pc.Meta.PaType = append(p.pc.Meta.PaType, u)
			}
		}

	case "meta/tags":
		for _, v := range value {
			p.pc.Meta.Tags = append(p.pc.Meta.Tags, v)
		}
	case "meta/used-by":
		for _, v := range value {
			p.pc.Meta.UsedBy = append(p.pc.Meta.UsedBy, v)
		}
	case "dependencies/open":
		for _, v := range value {
			if len(v) > 50 {
				return newErrorInvalidValue(key, " %s is too long.  (max 50 chars)", key)
			}
			p.pc.Dependencies.Open = append(p.pc.Dependencies.Open, v)
		}
	case "dependencies/proprietary":
		for _, v := range value {
			if len(v) > 50 {
				return newErrorInvalidValue(key, " %s is too long.  (max 50 chars)", key)
			}
			p.pc.Dependencies.Proprietary = append(p.pc.Dependencies.Proprietary, v)
		}
	case "dependencies/hardware":
		for _, v := range value {
			if len(v) > 50 {
				return newErrorInvalidValue(key, " %s is too long.  (max 50 chars)", key)
			}
			p.pc.Dependencies.Hardware = append(p.pc.Dependencies.Hardware, v)
		}
	default:
		return ErrorInvalidKey{key + " : Array of Strings"}
	}
	return nil
}

func (p *parser) decodeArrObj(key string, value map[interface{}]interface{}) error {
	switch key {
	case "maintenance/technical-contacts":
		for _, v := range value {
			var contact Contact

			for k, val := range v.(map[interface{}]interface{}) {
				if k.(string) == "name" {
					contact.Name = val.(string)
				} else if k.(string) == "email" {
					contact.Email = val.(string)
					if err := p.checkEmail(key, val.(string)); err != nil {
						return err
					}
				} else if k.(string) == "affiliation" {
					contact.Affiliation = val.(string)
				} else {
					return newErrorInvalidValue(key, " %s contains an invalid value", k)
				}
			}
			if contact.Name == "" {
				return newErrorInvalidValue(key, " name is mandatory.")
			}
			if contact.Email == "" {
				return newErrorInvalidValue(key, " email is mandatory.")
			}
			p.pc.Maintenance.TechnicalContacts = append(p.pc.Maintenance.TechnicalContacts, contact)
		}
	case "description/shortdesc":
		for _, v := range value {
			var descriptions Desc
			for k, val := range v.(map[interface{}]interface{}) {
				if len(val.(string)) > 100 {
					return newErrorInvalidValue(key, " %s is too long.  (max 100 chars)", k)
				}

				if k.(string) == "en" {
					descriptions.En = val.(string)
				} else if k.(string) == "it" {
					descriptions.It = val.(string)
				} else {
					return newErrorInvalidValue(key, " %s contains an invalid value", k)
				}
			}
			p.pc.Description.Shortdesc = append(p.pc.Description.Shortdesc, descriptions)
		}
	case "description/longdesc":
		for _, v := range value {
			var descriptions Desc
			for k, val := range v.(map[interface{}]interface{}) {
				if len(val.(string)) < 500 || len(val.(string)) > 10000 {
					return newErrorInvalidValue(key, " \"%s\" has an invalid number of characters: %d.  (min 500 chars, max 10000 chars)", k, len(val.(string)))
				}

				if k.(string) == "en" {
					descriptions.En = val.(string)
				} else if k.(string) == "it" {
					descriptions.It = val.(string)
				} else {
					return newErrorInvalidValue(key, " %s contains an invalid value", k)
				}
			}
			p.pc.Description.LongDesc = append(p.pc.Description.LongDesc, descriptions)
		}
	default:
		return ErrorInvalidKey{key + " : Array of Objects"}
	}
	return nil
}

func (p *parser) finalize() (es ErrorParseMulti) {
	// "maintenance/until" is mandatory (if the software is commercially maintained)
	if p.pc.Maintenance.Type == "commercial" && p.pc.Maintenance.Until.IsZero() {
		es = append(es, newErrorInvalidValue("maintenance/until", "not found, mandatory for a commercial maintenance"))
	}
	// "maintenance/maintainer" is mandatory  (if there is a maintenance)
	if &p.pc.Maintenance != nil {
		if p.pc.Maintenance.Maintainer == nil {
			es = append(es, newErrorInvalidValue("maintenance/maintainer", "not found, mandatory if  (if there is a maintenance)"))
		}
	}
	// "maintenance/technical-contacts" is mandatory  (if there is a maintenance)
	if &p.pc.Maintenance != nil {
		if p.pc.Maintenance.TechnicalContacts == nil {
			es = append(es, newErrorInvalidValue("maintenance/technical-contacts", "not found, mandatory if  (if there is a maintenance)"))
		}
	}

	// mandatory "description/released" if "description/version" is present
	if p.pc.Description.Version != "" && p.pc.Description.Released.IsZero() {
		es = append(es, newErrorInvalidValue("description/released", "not found, mandatory if a description/version is set"))
	}

	// mandatoryKeys check
	for k, v := range p.missing {
		if v {
			es = append(es, newErrorInvalidValue(k, k+" is a mandatory key."))
		}
	}

	return
}
