package ipa

import (
	"bufio"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/italia/developers-italia-backend/crawler/elastic"
	es "github.com/olivere/elastic"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

func localIPAFile() string {
	return path.Join(viper.GetString("CRAWLER_DATADIR"), "indicepa.csv")
}

// UpdateFromIndicePAIfNeeded downloads the amministrazioni.txt file if it's older than 20 days
// and loads it into Elasticsearch.
func UpdateFromIndicePAIfNeeded(elasticClient *es.Client) error {
	file := localIPAFile()

	needUpdate := true

	// we don't need to update if file does not exist and it's not older than 20 days
	info, err := os.Stat(file)
	if !os.IsNotExist(err) {
		if err != nil {
			log.Fatal(err)
			return err
		}
		if info.ModTime().After(time.Now().AddDate(0, 0, -20)) {
			needUpdate = false
		}
	}

	if needUpdate {
		return UpdateFromIndicePA(elasticClient)
	}

	return nil
}

// UpdateFromIndicePA downloads the amministrazioni.txt file and loads it into Elasticsearch.
func UpdateFromIndicePA(elasticClient *es.Client) error {
	file := localIPAFile()

	url := viper.GetString("INDICEPA_URL")
	log.Infof("Updating our cached copy from IndicePA from %v...", url)

	err := downloadFile(file, url)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Successfully updated from IndicePA")

	return saveIPAToElasticsearch(elasticClient)
}

func saveIPAToElasticsearch(elasticClient *es.Client) error {
	type officeES struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		PEC         string `json:"pec"`
	}
	type amministrazioneES struct {
		IPA         string     `json:"ipa"`
		Description string     `json:"description"`
		Type        string     `json:"type"`
		PEC         string     `json:"pec"`
		Offices     []officeES `json:"office"`
	}

	// Open file for parsing
	f, err := os.Open(localIPAFile())
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the CSV file
	reader := csv.NewReader(f)
	reader.Comma = '\t'
	reader.ReuseRecord = true
	reader.LazyQuotes = true
	lines, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Parse lines
	amms := make(map[string]amministrazioneES)
	for _, line := range lines {
		ipaCode := strings.ToLower(line[0])
		amm := amministrazioneES{
			IPA:         ipaCode,
			Description: line[1],
			Type:        line[12],
		}
		if line[17] == "pec" {
			amm.PEC = line[16]
		} else if line[19] == "pec" {
			amm.PEC = line[18]
		} else if line[21] == "pec" {
			amm.PEC = line[20]
		} else if line[23] == "pec" {
			amm.PEC = line[22]
		} else if line[25] == "pec" {
			amm.PEC = line[24]
		}
		amms[amm.IPA] = amm
	}

	// Download the list of AOO offices
	resp, err := http.Get(viper.GetString("INDICEPA_AOO_URL"))
	if err != nil {
		return err
	}

	// Read the AOO offices CSV file
	reader = csv.NewReader(resp.Body)
	reader.Comma = '\t'
	reader.ReuseRecord = true
	reader.LazyQuotes = true
	lines, err = reader.ReadAll()
	if err != nil {
		return err
	}
	for _, line := range lines {
		ipaCode := strings.ToLower(line[0])
		amm, ok := amms[ipaCode]
		if !ok {
			log.Debugf("skipping non-existing IPA code found in AOO: %s", line[0])
			continue
		}

		off := officeES{
			Code:        line[1],
			Description: line[2],
		}
		if line[16] == "pec" {
			off.PEC = line[15]
		} else if line[18] == "pec" {
			off.PEC = line[17]
		} else if line[20] == "pec" {
			off.PEC = line[19]
		}

		amm.Offices = append(amm.Offices, off)
		amms[ipaCode] = amm

		if len(amm.Offices) == 0 {
			log.Debugf("%s has no offices", ipaCode)
		}
	}

	// Download the list of OU offices
	resp, err = http.Get(viper.GetString("INDICEPA_OU_URL"))
	if err != nil {
		return err
	}

	// Read the OU offices CSV file
	reader = csv.NewReader(resp.Body)
	reader.Comma = '\t'
	reader.ReuseRecord = true
	reader.LazyQuotes = true
	lines, err = reader.ReadAll()
	if err != nil {
		return err
	}
	for _, line := range lines {
		ipaCode := strings.ToLower(line[13])
		amm, ok := amms[ipaCode]
		if !ok {
			log.Debugf("skipping non-existing IPA code found in OU: %s", line[0])
			continue
		}

		off := officeES{
			Code:        line[0],
			Description: line[2],
		}
		if line[18] == "pec" {
			off.PEC = line[17]
		} else if line[20] == "pec" {
			off.PEC = line[19]
		} else if line[22] == "pec" {
			off.PEC = line[21]
		}

		amm.Offices = append(amm.Offices, off)
		amms[ipaCode] = amm

		if len(amm.Offices) == 0 {
			log.Debugf("%s has no offices", ipaCode)
		}
	}

	if len(amms) == 0 {
		return fmt.Errorf("0 entities read from IndicePA; aborting")
	}

	// Delete existing index
	// TODO: use an alias for atomic updates!
	ctx := context.Background()
	_, err = elasticClient.DeleteIndex(viper.GetString("ELASTIC_INDICEPA_INDEX")).Do(ctx)
	if err != nil {
		return err
	}

	// Create mapping if it does not exist
	err = elastic.CreateIndexMapping(viper.GetString("ELASTIC_INDICEPA_INDEX"), elastic.IPAMapping, elasticClient)
	if err != nil {
		return err
	}

	// Perform a bulk request to Elasticsearch
	bulkRequest := elasticClient.Bulk()
	for _, amm := range amms {
		req := es.NewBulkIndexRequest().
			Index(viper.GetString("ELASTIC_INDICEPA_INDEX")).
			Type("pa").
			Id(amm.IPA).
			Doc(amm)
		bulkRequest.Add(req)
	}
	bulkResponse, err := bulkRequest.Do(ctx)
	if err != nil {
		return err
	}

	log.Debugf("%d records indexed", len(bulkResponse.Indexed()))

	return nil
}

// GetAdministrationName return the administration name associated to the "codice iPA" asssociated.
// TODO: load this mappings in memory instead of scanning the file every time
func GetAdministrationName(codiceiPA string) string {
	dataFile, err := ioutil.ReadFile(localIPAFile())
	if err != nil {
		log.Error(err)
		return ""
	}
	input := string(dataFile)

	// Scan the file, line by line.
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		amm := parseLine(scanner.Text())
		if amm.CodAmm == codiceiPA {
			return amm.DesAmm
		}
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("error reading standard input %v:", err)
	}

	return ""
}

// parseLine populate an Amministrazione with the values read.
func parseLine(line string) Amministrazione {
	data := strings.Split(line, "	")
	if len(data) < 31 {
		os.Remove(localIPAFile())
		panic(fmt.Sprintf("Line is shorter than expected [%s] - Removing the local CSV file as it might be corrupt; run this crawler again in order to download it again.", line))
	}
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
