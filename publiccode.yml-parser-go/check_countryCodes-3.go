package publiccode

import "strings"

// checkCountryCodes3 tells whether the 3-letter country code is valid (ISO 3166-1 alpha-3 code) or not and returns it.
func (p *parser) checkCountryCodes3(key, code string) error {
	// If it's not a valid 3 letters code.
	if len(code) != 3 {
		return newErrorInvalidValue(key, "invalid ISO 3166-1 alpha-3 country code: %s", code)
	}

	code = strings.ToUpper(code)
	for _, c := range countryCodes3 {
		if c == code {
			return nil
		}
	}
	return newErrorInvalidValue(key, "unknown ISO 3166-1 alpha-3 country code: %s", code)
}

// A countryCodes3 is a list of the valid (ISO 3166-1 alpha-3) codes.
// Updated: 2018-05-01
// Reference: https://www.iso.org/glossary-for-iso-3166.html
var countryCodes3 = []string{
	"AFG", // Afghanistan
	"ALA", // Åland Islands
	"ALB", // Albania
	"DZA", // Algeria
	"ASM", // American Samoa
	"AND", // Andorra
	"AGO", // Angola
	"AIA", // Anguilla
	"ATA", // Antarctica
	"ATG", // Antigua and Barbuda
	"ARG", // Argentina
	"ARM", // Armenia
	"ABW", // Aruba
	"AUS", // Australia
	"AUT", // Austria
	"AZE", // Azerbaijan
	"BHS", // Bahamas
	"BHR", // Bahrain
	"BGD", // Bangladesh
	"BRB", // Barbados
	"BLR", // Belarus
	"BEL", // Belgium
	"BLZ", // Belize
	"BEN", // Benin
	"BMU", // Bermuda
	"BTN", // Bhutan
	"BOL", // Bolivia (Plurinational State of)
	"BES", // Bonaire, Sint Eustatius and Saba
	"BIH", // Bosnia and Herzegovina
	"BWA", // Botswana
	"BVT", // Bouvet Island
	"BRA", // Brazil
	"IOT", // British Indian Ocean Territory
	"VGB", // British Virgin Islands
	"BRN", // Brunei Darussalam
	"BGR", // Bulgaria
	"BFA", // Burkina Faso
	"BDI", // Burundi
	"CPV", // Cabo Verde
	"KHM", // Cambodia
	"CMR", // Cameroon
	"CAN", // Canada
	"CYM", // Cayman Islands
	"CAF", // Central African Republic
	"TCD", // Chad
	"CHL", // Chile
	"CHN", // China
	"HKG", // China, Hong Kong Special Administrative Region
	"MAC", // China, Macao Special Administrative Region
	"CXR", // Christmas Island
	"CCK", // Cocos (Keeling) Islands
	"COL", // Colombia
	"COM", // Comoros
	"COG", // Congo
	"COK", // Cook Islands
	"CRI", // Costa Rica
	"CIV", // Côte d'Ivoire
	"HRV", // Croatia
	"CUB", // Cuba
	"CUW", // Curaçao
	"CYP", // Cyprus
	"CZE", // Czechia
	"PRK", // Democratic People's Republic of Korea
	"COD", // Democratic Republic of the Congo
	"DNK", // Denmark
	"DJI", // Djibouti
	"DMA", // Dominica
	"DOM", // Dominican Republic
	"ECU", // Ecuador
	"EGY", // Egypt
	"SLV", // El Salvador
	"GNQ", // Equatorial Guinea
	"ERI", // Eritrea
	"EST", // Estonia
	"ETH", // Ethiopia
	"FLK", // Falkland Islands (Malvinas)
	"FRO", // Faroe Islands
	"FJI", // Fiji
	"FIN", // Finland
	"FRA", // France
	"GUF", // French Guiana
	"PYF", // French Polynesia
	"ATF", // French Southern Territories
	"GAB", // Gabon
	"GMB", // Gambia
	"GEO", // Georgia
	"DEU", // Germany
	"GHA", // Ghana
	"GIB", // Gibraltar
	"GRC", // Greece
	"GRL", // Greenland
	"GRD", // Grenada
	"GLP", // Guadeloupe
	"GUM", // Guam
	"GTM", // Guatemala
	"GGY", // Guernsey
	"GIN", // Guinea
	"GNB", // Guinea-Bissau
	"GUY", // Guyana
	"HTI", // Haiti
	"HMD", // Heard Island and McDonald Islands
	"VAT", // Holy See (Vatican City State)
	"HND", // Honduras
	"HUN", // Hungary
	"ISL", // Iceland
	"IND", // India
	"IDN", // Indonesia
	"IRN", // Iran (Islamic Republic of)
	"IRQ", // Iraq
	"IRL", // Ireland
	"IMN", // Isle of Man
	"ISR", // Israel
	"ITA", // Italy
	"JAM", // Jamaica
	"JPN", // Japan
	"JEY", // Jersey
	"JOR", // Jordan
	"KAZ", // Kazakhstan
	"KEN", // Kenya
	"KIR", // Kiribati
	"KWT", // Kuwait
	"KGZ", // Kyrgyzstan
	"LAO", // Lao People's Democratic Republic
	"LVA", // Latvia
	"LBN", // Lebanon
	"LSO", // Lesotho
	"LBR", // Liberia
	"LBY", // Libya
	"LIE", // Liechtenstein
	"LTU", // Lithuania
	"LUX", // Luxembourg
	"MDG", // Madagascar
	"MWI", // Malawi
	"MYS", // Malaysia
	"MDV", // Maldives
	"MLI", // Mali
	"MLT", // Malta
	"MHL", // Marshall Islands
	"MTQ", // Martinique
	"MRT", // Mauritania
	"MUS", // Mauritius
	"MYT", // Mayotte
	"MEX", // Mexico
	"FSM", // Micronesia (Federated States of)
	"MCO", // Monaco
	"MNG", // Mongolia
	"MNE", // Montenegro
	"MSR", // Montserrat
	"MAR", // Morocco
	"MOZ", // Mozambique
	"MMR", // Myanmar
	"NAM", // Namibia
	"NRU", // Nauru
	"NPL", // Nepal
	"NLD", // Netherlands
	"NCL", // New Caledonia
	"NZL", // New Zealand
	"NIC", // Nicaragua
	"NER", // Niger
	"NGA", // Nigeria
	"NIU", // Niue
	"NFK", // Norfolk Island
	"MNP", // Northern Mariana Islands
	"NOR", // Norway
	"OMN", // Oman
	"PAK", // Pakistan
	"PLW", // Palau
	"PAN", // Panama
	"PNG", // Papua New Guinea
	"PRY", // Paraguay
	"PER", // Peru
	"PHL", // Philippines
	"PCN", // Pitcairn
	"POL", // Poland
	"PRT", // Portugal
	"PRI", // Puerto Rico
	"QAT", // Qatar
	"KOR", // Republic of Korea
	"MDA", // Republic of Moldova
	"REU", // Réunion
	"ROU", // Romania
	"RUS", // Russian Federation
	"RWA", // Rwanda
	"BLM", // Saint Barthélemy
	"SHN", // Saint Helena
	"KNA", // Saint Kitts and Nevis
	"LCA", // Saint Lucia
	"MAF", // Saint Martin (French Part)
	"SPM", // Saint Pierre and Miquelon
	"VCT", // Saint Vincent and the Grenadines
	"WSM", // Samoa
	"SMR", // San Marino
	"STP", // Sao Tome and Principe
	"",    // Sark
	"SAU", // Saudi Arabia
	"SEN", // Senegal
	"SRB", // Serbia
	"SYC", // Seychelles
	"SLE", // Sierra Leone
	"SGP", // Singapore
	"SXM", // Sint Maarten (Dutch part)
	"SVK", // Slovakia
	"SVN", // Slovenia
	"SLB", // Solomon Islands
	"SOM", // Somalia
	"ZAF", // South Africa
	"SGS", // South Georgia and the South Sandwich Islands
	"SSD", // South Sudan
	"ESP", // Spain
	"LKA", // Sri Lanka
	"PSE", // State of Palestine
	"SDN", // Sudan
	"SUR", // Suriname
	"SJM", // Svalbard and Jan Mayen Islands
	"SWZ", // Swaziland
	"SWE", // Sweden
	"CHE", // Switzerland
	"SYR", // Syrian Arab Republic
	"TJK", // Tajikistan
	"THA", // Thailand
	"MKD", // The former Yugoslav Republic of Macedonia
	"TLS", // Timor-Leste
	"TGO", // Togo
	"TKL", // Tokelau
	"TON", // Tonga
	"TTO", // Trinidad and Tobago
	"TUN", // Tunisia
	"TUR", // Turkey
	"TKM", // Turkmenistan
	"TCA", // Turks and Caicos Islands
	"TUV", // Tuvalu
	"UGA", // Uganda
	"UKR", // Ukraine
	"ARE", // United Arab Emirates
	"GBR", // United Kingdom of Great Britain and Northern Ireland
	"TZA", // United Republic of Tanzania
	"UMI", // United States Minor Outlying Islands
	"USA", // United States of America
	"VIR", // United States Virgin Islands
	"URY", // Uruguay
	"UZB", // Uzbekistan
	"VUT", // Vanuatu
	"VEN", // Venezuela (Bolivarian Republic of)
	"VNM", // Viet Nam
	"WLF", // Wallis and Futuna Islands
	"ESH", // Western Sahara
	"YEM", // Yemen
	"ZMB", // Zambia
	"ZWE", // Zimbabwe
}
