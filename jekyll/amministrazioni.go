package jekyll

import (
	"bufio"
	"context"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/italia/developers-italia-backend/crawler"
	"github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// Administration is a simple description of an Administration.
type Administration struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	CodiceIPA string `json:"ipa"`
}

// AmministrazioniYML generate a yml file with all the amministrazioni in es.
func AmministrazioniYML(filename string, elasticClient *elastic.Client) error {
	log.Debug("Generating amministrazioni.yml")

	// Create file if not exists.
	if _, err := os.Stat(filename); os.IsExist(err) {
		err := os.Remove(filename)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	// Open file.
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close() // nolint: errcheck

	// Administrations data.
	var administrations []Administration

	// Extract all the documents.
	searchResult, err := elasticClient.Search().
		Index("publiccode").               // search in index "publiccode"
		Query(elastic.NewMatchAllQuery()). // specify the query
		Pretty(true).                      // pretty print request and response JSON
		From(0).Size(10000).               // get first 10k elements. It can be changed.
		Do(context.Background())           // execute
	if err != nil {
		log.Error(err)
	}

	// Foreach search result check if codiceIPA is not empty.
	var pctype crawler.PublicCodeES
	for _, item := range searchResult.Each(reflect.TypeOf(pctype)) {
		i := item.(crawler.PublicCodeES)
		if i.ItRiusoCodiceIPA != "" {
			administrations = append(administrations, Administration{
				Name:      getNomeAmministrazione(i.ItRiusoCodiceIPA),
				URL:       i.LandingURL,
				CodiceIPA: i.ItRiusoCodiceIPA,
			})
		}
	}
	// Debug note if file will be empty.
	if len(administrations) == 0 {
		log.Debug("amministrazioni.yml is empty.")
	}

	// Remove duplicates.
	administrations = removeDuplicates(administrations)

	// Marshal yml.
	d, err := yaml.Marshal(&administrations)
	if err != nil {
		return err
	}
	//Append data to file.
	if _, err = f.WriteString(string(d)); err != nil {
		return err
	}

	return err
}

func removeDuplicates(elements []Administration) []Administration {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []Administration{}

	for v := range elements {
		if encountered[elements[v].CodiceIPA] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v].CodiceIPA] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// Amministrazione is an Administration from amministrazoni.txt
// Retrieved from: http://www.indicepa.gov.it/documentale/n-opendata.php
type Amministrazione struct {
	codAmm            string
	desAmm            string
	Comune            string
	nomeResp          string
	cognResp          string
	Cap               string
	Provincia         string
	Regione           string
	sitoIstituzionale string
	Indirizzo         string
	titoloResp        string
	tipologiaIstat    string
	tipologiaAmm      string
	acronimo          string
	cfValidato        string
	Cf                string
	mail1             string
	tipoMail1         string
	mail2             string
	tipoMail2         string
	mail3             string
	tipoMail3         string
	mail4             string
	tipoMail4         string
	mail5             string
	tipoMail5         string
	urlFacebook       string
	urlTwitter        string
	urlGoogleplus     string
	urlYoutube        string
	livAccessibili    string
}

func getNomeAmministrazione(codiceiPA string) string {
	file := "jekyll/amministrazioni.txt"
	dataFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error(err)
		return ""
	}
	input := string(dataFile)
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		amm := manageLine(scanner.Text())
		if amm.codAmm == codiceiPA {
			return amm.desAmm
		}
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("error reading standard input %v:", err)
	}

	return ""
}

func manageLine(line string) Amministrazione {
	data := strings.Split(line, "	")
	amm := Amministrazione{
		codAmm:            data[0],
		desAmm:            data[1],
		Comune:            data[2],
		nomeResp:          data[3],
		cognResp:          data[4],
		Cap:               data[5],
		Provincia:         data[6],
		Regione:           data[7],
		sitoIstituzionale: data[8],
		Indirizzo:         data[9],
		titoloResp:        data[10],
		tipologiaIstat:    data[11],
		tipologiaAmm:      data[12],
		acronimo:          data[13],
		cfValidato:        data[14],
		Cf:                data[15],
		mail1:             data[16],
		tipoMail1:         data[17],
		mail2:             data[18],
		tipoMail2:         data[19],
		mail3:             data[20],
		tipoMail3:         data[21],
		mail4:             data[22],
		tipoMail4:         data[23],
		mail5:             data[24],
		tipoMail5:         data[25],
		urlFacebook:       data[26],
		urlTwitter:        data[27],
		urlGoogleplus:     data[28],
		urlYoutube:        data[29],
		livAccessibili:    data[30],
	}

	return amm
}
