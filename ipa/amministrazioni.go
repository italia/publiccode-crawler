package ipa

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// Amministrazione is an Administration from amministrazoni.txt
// Retrieved from: http://www.indicepa.gov.it/documentale/n-opendata.php
type Amministrazione struct {
	CodAmm            string
	DesAmm            string
	Comune            string
	NomeResp          string
	CognResp          string
	Cap               string
	Provincia         string
	Regione           string
	SitoIstituzionale string
	Indirizzo         string
	TitoloResp        string
	TipologiaIstat    string
	TipologiaAmm      string
	Acronimo          string
	CFValidato        string
	CF                string
	Mail1             string
	TipoMail1         string
	Mail2             string
	TipoMail2         string
	Mail3             string
	TipoMail3         string
	Mail4             string
	TipoMail4         string
	Mail5             string
	TipoMail5         string
	URLFacebook       string
	URLTwitter        string
	URLGoogleplus     string
	URLYoutube        string
	LivAccessibili    string
}

// UpdateFile download the amministrazioni.txt file if it's older than 2 days.
func UpdateFile(fileName, fileURL string) error {
	info, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
		return err
	}
	today := time.Now()
	older := today.AddDate(0, 0, -2)

	downloadTime := info.ModTime()

	// If amministrazioni.txt is older that 2 days.
	if downloadTime.Before(older) {
		log.Info("download a new amministrazioni.txt ...")

		err := downloadFile(fileName, fileURL)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return err
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

	// Scan the file, line by line.
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		amm := manageLine(scanner.Text())
		if amm.CodAmm == codiceiPA {
			return amm.DesAmm
		}
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("error reading standard input %v:", err)
	}

	return ""
}

// manageLine populate an Amministrazione with the values read.
func manageLine(line string) Amministrazione {
	data := strings.Split(line, "	")
	amm := Amministrazione{
		CodAmm:            data[0],
		DesAmm:            data[1],
		Comune:            data[2],
		NomeResp:          data[3],
		CognResp:          data[4],
		Cap:               data[5],
		Provincia:         data[6],
		Regione:           data[7],
		SitoIstituzionale: data[8],
		Indirizzo:         data[9],
		TitoloResp:        data[10],
		TipologiaIstat:    data[11],
		TipologiaAmm:      data[12],
		Acronimo:          data[13],
		CFValidato:        data[14],
		CF:                data[15],
		Mail1:             data[16],
		TipoMail1:         data[17],
		Mail2:             data[18],
		TipoMail2:         data[19],
		Mail3:             data[20],
		TipoMail3:         data[21],
		Mail4:             data[22],
		TipoMail4:         data[23],
		Mail5:             data[24],
		TipoMail5:         data[25],
		URLFacebook:       data[26],
		URLTwitter:        data[27],
		URLGoogleplus:     data[28],
		URLYoutube:        data[29],
		LivAccessibili:    data[30],
	}

	return amm
}

func downloadFile(filepath string, url string) error {

	// Create the file.
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			log.Error(err)
		}
	}()

	// Get the data from the url.
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error(err)
		}
	}()

	// Write the body to file.
	_, err = io.Copy(out, resp.Body)

	return err
}
