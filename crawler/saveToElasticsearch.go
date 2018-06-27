package crawler

import (
	"context"
	"net/url"

	"github.com/italia/developers-italia-backend/metrics"
	"github.com/olivere/elastic"
	pcode "github.com/publiccodenet/publiccode.yml-parser-go"
	log "github.com/sirupsen/logrus"
)

// SaveToES save the chosen data []byte in elasticsearch
func SaveToES(domain Domain, name string, activityIndex float64, data []byte, index string, elasticClient *elastic.Client) error {
	const (
		// Elasticsearch mapping for publiccode. Check elasticsearch/mappings/.
		mapping = `{
	  "mappings": {
	    "software": {
	      "dynamic_templates": [
	        {
	          "description": {
	            "path_match": "description.*",
	            "mapping": {
	              "type": "object",
	              "properties": {
	                "localisedName": {
	                  "type": "text"
	                },
	                "genericName": {
	                  "type": "text",
	                  "fields": {
	                    "keyword": { "type": "keyword", "ignore_above": 256 }
	                  }
	                },
	                "shortDescription": {
	                  "type": "text"
	                },
	                "longDescription": {
	                  "type": "text"
	                },
	                "documentation": {
	                  "type": "text",
	                  "index": false
	                },
	                "apiDocumentation": {
	                  "type": "text",
	                  "index": false
	                },
	                "featureList": {
	                  "type": "keyword"
	                },
	                "freeTags": {
	                  "type": "keyword"
	                },
	                "screenshots": {
	                  "type": "keyword",
	                  "index": false
	                },
	                "videos": {
	                  "type": "keyword",
	                  "index": false
	                },
	                "awards": {
	                  "type": "keyword"
	                }
	              }
	            }
	          }
	        }
	      ],
	      "properties": {
	        "publiccode-yaml-version": {
	          "type": "text",
	          "index": false
	        },
	        "name": {
	          "type": "text"
	        },
	        "applicationSuite": {
	          "type": "text",
	          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
	        },
	        "url": {
	          "type": "text",
	          "index": false,
	          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
	        },
	        "landingURL": {
	          "type": "text",
	          "index": false
	        },
	        "isBasedOn": {
	          "type": "text",
	          "index": false
	        },
	        "softwareVersion": {
	          "type": "keyword"
	        },
	        "releaseDate": {
	          "type": "date",
	          "format": "strict_date"
	        },
	        "logo": {
	          "type": "text",
	          "index": false
	        },
	        "monochromeLogo": {
	          "type": "text",
	          "index": false
	        },
	        "inputTypes": {
	          "type": "keyword"
	        },
	        "outputTypes": {
	          "type": "keyword"
	        },
	        "platforms": {
	          "type": "keyword"
	        },
	        "tags": {
	          "type": "keyword"
	        },
	        "usedBy": {
	          "type": "text",
	          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
	        },
	        "roadmap": {
	          "type": "text",
	          "index": false
	        },
	        "developmentStatus": {
	          "type": "keyword"
	        },
	        "softwareType": {
	          "type": "keyword"
	        },
	        "intendedAudience-onlyFor": {
	          "type": "keyword"
	        },
	        "intendedAudience-countries": {
	          "type": "keyword"
	        },
	        "intendedAudience-unsupportedCountries": {
	          "type": "keyword"
	        },
	        "legal-license": {
	          "type": "text",
	          "fields": { "keyword": { "type": "keyword", "ignore_above": 256 } }
	        },
	        "legal-mainCopyrightOwner": {
	          "type": "text"
	        },
	        "legal-repoOwner": {
	          "type": "text"
	        },
	        "legal-authorsFile": {
	          "type": "text",
	          "index": false
	        },
	        "maintenance-type": {
	          "type": "keyword"
	        },
	        "maintenance-contractors": {
	          "type": "nested",
	          "properties": {
	            "name": {
	              "type": "text"
	            },
	            "until": {
	              "type": "date",
	              "format": "strict_date"
	            },
	            "website": {
	              "type": "text",
	              "index": false
	            }
	          }
	        },
	        "maintenance-contacts": {
	          "type": "nested",
	          "properties": {
	            "name": {
	              "type": "text"
	            },
	            "email": {
	              "type": "text"
	            },
	            "phone": {
	              "type": "text",
	              "index": false
	            },
	            "affiliation": {
	              "type": "text"
	            }
	          }
	        },
	        "localisation-localisationReady": {
	          "type": "boolean"
	        },
	        "localisation-availableLanguages": {
	          "type": "keyword"
	        },
	        "dependsOn-open": {
	          "type": "nested",
	          "properties": {
	            "name": {
	              "type": "text"
	            },
	            "version-min": {
	              "type": "text",
	              "index": false
	            },
	            "version-max": {
	              "type": "text",
	              "index": false
	            },
	            "version": {
	              "type": "text",
	              "index": false
	            },
	            "optional": {
	              "type": "boolean"
	            }
	          }
	        },
	        "dependsOn-proprietary": {
	          "type": "nested",
	          "properties": {
	            "name": {
	              "type": "text"
	            },
	            "version-min": {
	              "type": "text",
	              "index": false
	            },
	            "version-max": {
	              "type": "text",
	              "index": false
	            },
	            "version": {
	              "type": "text",
	              "index": false
	            },
	            "optional": {
	              "type": "boolean"
	            }
	          }
	        },
	        "dependsOn-hardware": {
	          "type": "nested",
	          "properties": {
	            "name": {
	              "type": "text"
	            },
	            "version-min": {
	              "type": "text",
	              "index": false
	            },
	            "version-max": {
	              "type": "text",
	              "index": false
	            },
	            "version": {
	              "type": "text",
	              "index": false
	            },
	            "optional": {
	              "type": "boolean"
	            }
	          }
	        },
	        "it-conforme-accessibile": {
	          "type": "boolean"
	        },
	        "it-conforme-interoperabile": {
	          "type": "boolean"
	        },
	        "it-conforme-sicuro": {
	          "type": "boolean"
	        },
	        "it-conforme-privacy": {
	          "type": "boolean"
	        },
	        "it-spid": {
	          "type": "boolean"
	        },
	        "it-cie": {
	          "type": "boolean"
	        },
	        "it-anpr": {
	          "type": "boolean"
	        },
	        "it-pagopa": {
	          "type": "boolean"
	        },
	        "it-riuso-codiceIPA": {
	          "type": "keyword"
	        },
	        "it-ecosistemi": {
	          "type": "keyword"
	        },
	        "it-designKit-seo": {
	          "type": "boolean"
	        },
	        "it-designKit-ui": {
	          "type": "boolean"
	        },
	        "it-designKit-web": {
	          "type": "boolean"
	        },
	        "it-designKit-content": {
	          "type": "boolean"
	        },
	        "suggest-name": {
	          "type": "completion"
	        },
	        "vitality-score": {
	          "type": "text",
	          "index": false
	        },
	        "vitality-dataChart": {
	          "type": "integer"
	        },
	        "related-software": {
	          "properties": {
	            "name": {
	              "type": "text",
	              "index": false
	            },
	            "image": {
	              "type": "text",
	              "index": false
	            },
	            "eng": {
	              "properties": {
	                "localised-name": {
	                  "type": "text",
	                  "index": false
	                },
	                "url": {
	                  "type": "text",
	                  "index": false
	                }
	              }
	            },
	            "ita": {
	              "properties": {
	                "localised-name": {
	                  "type": "text",
	                  "index": false
	                },
	                "url": {
	                  "type": "text",
	                  "index": false
	                }
	              }
	            }
	          }
	        },
	        "tags-related": {
	          "type": "keyword"
	        },
	        "popular-tags": {
	          "type": "keyword"
	        },
	        "share-tags": {
	          "type": "keyword"
	        },
	        "old-variant": {
	          "properties": {
	            "name": {
	              "type": "text",
	              "index": false
	            },
	            "eng": {
	              "properties": {
	                "localised-name": {
	                  "type": "text",
	                  "index": false
	                },
	                "url": {
	                  "type": "text",
	                  "index": false
	                },
	                "feature-list": {
	                  "type": "keyword",
	                  "index": false
	                },
	                "vitality-score": {
	                  "type": "integer",
	                  "index": false
	                },
	                "legal-repo-owner": {
	                  "type": "text",
	                  "index": false
	                }
	              }
	            },
	            "ita": {
	              "properties": {
	                "localised-name": {
	                  "type": "text",
	                  "index": false
	                },
	                "url": {
	                  "type": "text",
	                  "index": false
	                },
	                "feature-list": {
	                  "type": "keyword",
	                  "index": false
	                },
	                "vitality-score": {
	                  "type": "integer",
	                  "index": false
	                },
	                "legal-repo-owner": {
	                  "type": "text",
	                  "index": false
	                }
	              }
	            }
	          }
	        },
	        "old-feature-list": {
	          "properties": {
	            "ita": {
	              "type": "keyword",
	              "index": false
	            },
	            "eng": {
	              "type": "keyword",
	              "index": false
	            }
	          }
	        }
	      }
	    }
	  }
	}
` // nolint: misspell
	)

	// Starting with elastic.v5, you must pass a context to execute each service.
	ctx := context.Background()

	// Generate publiccode data using the parser.
	pc := pcode.PublicCode{}
	err := pcode.Parse(data, &pc)
	//	yaml.Unmarshal([]byte(data), &pc)
	if err != nil {
		log.Errorf("Error in publiccode.yml for %s: %v", name, err)
	}

	// Add a document to the index.
	file := PublicCodeES{
		PubliccodeYamlVersion: pc.PubliccodeYamlVersion,

		Name:             pc.Name,
		ApplicationSuite: pc.ApplicationSuite,
		URL:              pc.URL.String(),
		LandingURL:       pc.LandingURL.String(),

		IsBasedOn:       pc.IsBasedOn,
		SoftwareVersion: pc.SoftwareVersion,
		ReleaseDate:     pc.ReleaseDate.Format("2006-01-02"),
		Logo:            pc.Logo,
		MonochromeLogo:  pc.MonochromeLogo,
		InputTypes:      pc.InputTypes,
		OutputTypes:     pc.OutputTypes,

		Platforms: pc.Platforms,

		Tags: pc.Tags,

		FreeTags: pc.FreeTags,

		UsedBy: pc.UsedBy,

		Roadmap: pc.Roadmap.String(),

		DevelopmentStatus: pc.DevelopmentStatus,

		VitalityScore:     activityIndex,
		VitalityDataChart: []int{12, 24, 36, 48, 60, 72, 84, 96, 48},

		RelatedSoftware: nil,

		SoftwareType: pc.SoftwareType,

		IntendedAudienceOnlyFor:              pc.IntendedAudience.OnlyFor,
		IntendedAudienceCountries:            pc.IntendedAudience.Countries,
		IntendedAudienceUnsupportedCountries: pc.IntendedAudience.UnsupportedCountries,

		Description: map[string]Desc{},
		//OldVariants: oldVariant will be added in the search function.

		LegalLicense:            pc.Legal.License,
		LegalMainCopyrightOwner: pc.Legal.MainCopyrightOwner,
		LegalRepoOwner:          pc.Legal.RepoOwner,
		LegalAuthorsFile:        pc.Legal.AuthorsFile,

		MaintenanceType:        pc.Maintenance.Type,
		MaintenanceContractors: []Contractor{},
		MaintenanceContacts:    []Contact{},

		LocalisationLocalisationReady:  pc.Localisation.LocalisationReady,
		LocalisationAvailableLanguages: pc.Localisation.AvailableLanguages,

		DependenciesOpen:        []Dependency{},
		DependenciesProprietary: []Dependency{},
		DependenciesHardware:    []Dependency{},

		ItConformeAccessibile:    pc.It.Conforme.Accessibile,
		ItConformeInteroperabile: pc.It.Conforme.Interoperabile,
		ItConformeSicuro:         pc.It.Conforme.Sicuro,
		ItConformePrivacy:        pc.It.Conforme.Privacy,

		ItRiusoCodiceIPA: pc.It.Riuso.CodiceIPA,

		ItSpid:   pc.It.Spid,
		ItPagopa: pc.It.Pagopa,
		ItCie:    pc.It.Cie,
		ItAnpr:   pc.It.Anpr,

		ItEcosistemi: pc.It.Ecosistemi,

		ItDesignKitSeo:     pc.It.DesignKit.Seo,
		ItDesignKitUI:      pc.It.DesignKit.UI,
		ItDesignKitWeb:     pc.It.DesignKit.Web,
		ItDesignKitContent: pc.It.DesignKit.Content,
	}
	// Populate complex fields.
	for _, contractor := range pc.Maintenance.Contractors {
		file.MaintenanceContractors = append(file.MaintenanceContractors, Contractor{
			Name:    contractor.Name,
			Website: contractor.Website.String(),
			Until:   contractor.Until.String(),
		})
	}
	for _, contact := range pc.Maintenance.Contacts {
		file.MaintenanceContacts = append(file.MaintenanceContacts, Contact{
			Name:        contact.Name,
			Email:       contact.Email,
			Affiliation: contact.Affiliation,
			Phone:       contact.Phone,
		})
	}
	for lang := range pc.Description {
		file.Description[lang] = Desc{
			LocalisedName:    pc.Description[lang].LocalisedName,
			GenericName:      pc.Description[lang].GenericName,
			ShortDescription: pc.Description[lang].ShortDescription,
			LongDescription:  pc.Description[lang].LongDescription,
			Documentation:    pc.Description[lang].Documentation.String(),
			APIDocumentation: pc.Description[lang].APIDocumentation.String(),
			FeatureList:      pc.Description[lang].FeatureList,
			Screenshots: func(screenshots []string) []string {
				var s []string
				s = append(s, screenshots...)
				return s
			}(pc.Description[lang].Screenshots),
			Videos: func(videos []*url.URL) []string {
				var v []string
				for _, video := range videos {
					v = append(v, video.String())
				}
				return v
			}(pc.Description[lang].Videos),
			Awards: pc.Description[lang].Awards,
		}

	}
	for _, dependency := range pc.Dependencies.Open {
		file.DependenciesOpen = append(file.DependenciesOpen, Dependency{
			Name:       dependency.Name,
			VersionMin: dependency.VersionMin,
			VersionMax: dependency.VersionMax,
			Optional:   dependency.Optional,
			Version:    dependency.Version,
		})
	}
	for _, dependency := range pc.Dependencies.Proprietary {
		file.DependenciesProprietary = append(file.DependenciesProprietary, Dependency{
			Name:       dependency.Name,
			VersionMin: dependency.VersionMin,
			VersionMax: dependency.VersionMax,
			Optional:   dependency.Optional,
			Version:    dependency.Version,
		})
	}
	for _, dependency := range pc.Dependencies.Hardware {
		file.DependenciesHardware = append(file.DependenciesHardware, Dependency{
			Name:       dependency.Name,
			VersionMin: dependency.VersionMin,
			VersionMax: dependency.VersionMax,
			Optional:   dependency.Optional,
			Version:    dependency.Version,
		})
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := elasticClient.IndexExists(index).Do(context.Background())
	if err != nil {
		return err
	}
	if !exists {
		_, err := elasticClient.CreateIndex(index).Body(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			return err
		}
	}

	// Put publiccode data in ES.
	_, err = elasticClient.Index().
		Index(index).
		Type("software").
		Id(domain.Host + "/" + name + "_" + index).
		BodyJson(file).
		Do(ctx)
	if err != nil {
		return err
	}

	metrics.GetCounter("repository_file_indexed", index).Inc()

	return nil
}
