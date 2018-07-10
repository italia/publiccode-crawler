package ipa

import (
	"bufio"
	"io/ioutil"
	"strings"

	log "github.com/sirupsen/logrus"
)

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

// GetAdministrationName return the administration name associated to the "codice iPA" asssociated.
func GetAdministrationName(codiceiPA string) string {
	file := "./ipa/amministrazioni.txt"
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
