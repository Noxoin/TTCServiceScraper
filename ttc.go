package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Prop struct {
	Endpoint         string
	V                string
	S                string
	L                string
	P                int64
	DefaultSortOrder string
	Sig              string
	ItemId           string
	AutoFireSearch   bool
}

func (p *Prop) GetURL() (string, error) {
	u, err := url.Parse(p.Endpoint)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Add("v", p.V)
	q.Add("s", p.S)
	q.Add("l", p.L)
	q.Add("p", strconv.FormatInt(p.P, 10))
	q.Add("defaultSortOrder", p.DefaultSortOrder)
	q.Add("sig", p.Sig)
	q.Add("itemid", p.ItemId)
	q.Add("autoFireSearch", strconv.FormatBool(p.AutoFireSearch))
	u.RawQuery = q.Encode()
	return u.String(), nil
}

const PROPERTIES = `{"endpoint":"https://www.ttc.ca//sxa/search/results/","v":"{23DC07D4-6BAC-4B98-A9CC-07606C5B1322}","s":"{99D7699F-DB47-4BB1-8946-77561CE7B320}","l":"","p":10,"defaultSortOrder":"ContentDateFacet,Descending","sig":"","itemid":"{72CC555F-9128-4581-AD12-3D04AB1C87BA}","autoFireSearch":true}`

func QueryTTCSubwayServiceLandingPageProperties(mock bool) (string, error) {
	if mock {
		return PROPERTIES, nil
	}
	resp, err := http.Get("https://www.ttc.ca/service-advisories/subway-service")
	if err != nil {
		fmt.Printf("Error in querying ttc page: %v", err)
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Error in parsing ttc response body: %v", err)
		return "", err
	}
	properties, exists := doc.Find(".search-results").Attr("data-properties")
	if !exists {
		fmt.Printf("No such attribute 'data-properties' in search-results HTML Node")
		return "", err
	}
	return properties, nil
}

const CACHED_SEARCH_RESULTS = `{"TotalTime":54,"CountTime":0,"QueryTime":18,"Signature":"","Index":"sitecore_sxa_web_index","Count":5,"Results":[{"Id":"32e64882-7314-499e-9c42-162098528e23","Language":"en","Path":"/sitecore/content/TTC/DevProto/Home/service-advisories/subway-service/1 Line 1 Lawrence to St Clair full weekend closure March 12 and 13","Url":"/service-advisories/subway-service/1-Line-1-Lawrence-to-St-Clair-full-weekend-closure-March-12-and-13","Name":null,"Html":"<a title=\"1 Line 1 Lawrence to St Clair full weekend closure March 12 and 13\" href=\"/service-advisories/subway-service/1-Line-1-Lawrence-to-St-Clair-full-weekend-closure-March-12-and-13\" role=\"heading\" aria-level=\"2\"><div class=\"sa-title c-news-results__link u-mb-sm\"><span class=\"field-routename\">line 1</span><span class=\"sa-dash\">&#8211;</span><span class=\"field-satitle\">Lawrence to St Clair full weekend closure March 12 and 13</span></div></a><div class=\"sa-subtitle c-news-results__date field-sasubtitle\"></div><div class=\"sa-effective-date c-news-results__date\"><span class=\"ed-start-date field-starteffectivedate\">March 12, 2022<span>&nbsp;</span></span><span class=\"effective-date-tolabel\">to </span><span class=\"field-endeffectivedate\">March 13, 2022</span></div>"},{"Id":"f775661d-2ce1-4f6a-a9e8-06da85d4ef9d","Language":"en","Path":"/sitecore/content/TTC/DevProto/Home/service-advisories/subway-service/1 Line 1 Finch to Eglinton nightly early closures March 7 to 10","Url":"/service-advisories/subway-service/1-Line-1-Finch-to-Eglinton-nightly-early-closures-March-7-to-10","Name":null,"Html":"<a title=\"1 Line 1 Finch to Eglinton nightly early closures March 7 to 10\" href=\"/service-advisories/subway-service/1-Line-1-Finch-to-Eglinton-nightly-early-closures-March-7-to-10\" role=\"heading\" aria-level=\"2\"><div class=\"sa-title c-news-results__link u-mb-sm\"><span class=\"field-routename\">line 1</span><span class=\"sa-dash\">&#8211;</span><span class=\"field-satitle\">Finch to Eglinton nightly early closures March 7 to 10</span></div></a><div class=\"sa-subtitle c-news-results__date field-sasubtitle\"></div><div class=\"sa-effective-date c-news-results__date\"><span class=\"ed-start-date field-starteffectivedate\">March 7, 2022<span>&nbsp;</span></span><span class=\"effective-date-tolabel\">to </span><span class=\"field-endeffectivedate\">March 10, 2022</span></div>"},{"Id":"38fe9f34-bcff-48d3-a519-1dd8f37a340e","Language":"en","Path":"/sitecore/content/TTC/DevProto/Home/service-advisories/subway-service/1 Line 1 Lawrence to St Clair single day closure March 19","Url":"/service-advisories/subway-service/1-Line-1-Lawrence-to-St-Clair-single-day-closure-March-19","Name":null,"Html":"<a title=\"1 Line 1 Lawrence to St Clair single day closure March 19\" href=\"/service-advisories/subway-service/1-Line-1-Lawrence-to-St-Clair-single-day-closure-March-19\" role=\"heading\" aria-level=\"2\"><div class=\"sa-title c-news-results__link u-mb-sm\"><span class=\"field-routename\">line 1</span><span class=\"sa-dash\">&#8211;</span><span class=\"field-satitle\">Lawrence to St Clair single day closure March 19</span></div></a><div class=\"sa-subtitle c-news-results__date field-sasubtitle\"></div><div class=\"sa-effective-date c-news-results__date\"><span class=\"ed-start-date field-starteffectivedate\">March 19, 2022 - 12:00 AM<span>&nbsp;</span></span><span class=\"effective-date-tolabel\">to </span><span class=\"field-endeffectivedate\">12:00 AM</span></div>"},{"Id":"fed6deda-da7a-4f50-8769-a8cd5f91e909","Language":"en","Path":"/sitecore/content/TTC/DevProto/Home/service-advisories/subway-service/1 Line 1 Yonge University Vaughan Metropolitan Centre to Wilson nightly early closures March 21 to 2","Url":"/service-advisories/subway-service/1-Line-1-Yonge-University-Vaughan-Metropolitan-Centre-to-Wilson-nightly-early-closures-March-21-to-2","Name":null,"Html":"<a title=\"1 Line 1 Yonge University Vaughan Metropolitan Centre to Wilson nightly early closures March 21 to 2\" href=\"/service-advisories/subway-service/1-Line-1-Yonge-University-Vaughan-Metropolitan-Centre-to-Wilson-nightly-early-closures-March-21-to-2\" role=\"heading\" aria-level=\"2\"><div class=\"sa-title c-news-results__link u-mb-sm\"><span class=\"field-routename\">line 1 (yonge-university)</span><span class=\"sa-dash\">&#8211;</span><span class=\"field-satitle\">Vaughan Metropolitan Centre to Wilson nightly early closures March 21 to 24</span></div></a><div class=\"sa-subtitle c-news-results__date field-sasubtitle\"></div><div class=\"sa-effective-date c-news-results__date\"><span class=\"ed-start-date field-starteffectivedate\">March 21, 2022<span>&nbsp;</span></span><span class=\"effective-date-tolabel\">to </span><span class=\"field-endeffectivedate\">March 24, 2022</span></div>"},{"Id":"a2fa6c85-4d8d-4546-a96c-883b75a6f032","Language":"en","Path":"/sitecore/content/TTC/DevProto/Home/service-advisories/subway-service/1 Line 1 Finch to Eglinton nightly early closures March 14 to 17","Url":"/service-advisories/subway-service/1-Line-1-Finch-to-Eglinton-nightly-early-closures-March-14-to-17","Name":null,"Html":"<a title=\"1 Line 1 Finch to Eglinton nightly early closures March 14 to 17\" href=\"/service-advisories/subway-service/1-Line-1-Finch-to-Eglinton-nightly-early-closures-March-14-to-17\" role=\"heading\" aria-level=\"2\"><div class=\"sa-title c-news-results__link u-mb-sm\"><span class=\"field-routename\">line 1</span><span class=\"sa-dash\">&#8211;</span><span class=\"field-satitle\">Finch to Eglinton nightly early closures March 14 to 17</span></div></a><div class=\"sa-subtitle c-news-results__date field-sasubtitle\"></div><div class=\"sa-effective-date c-news-results__date\"><span class=\"ed-start-date field-starteffectivedate\">March 14, 2022<span>&nbsp;</span></span><span class=\"effective-date-tolabel\">to </span><span class=\"field-endeffectivedate\">March 17, 2022</span></div>"}]}`

type Results struct {
	TotalTime int64
	CountTime int64
	QueryTime int64
	Signature string
	Index     string
	Count     int64
	Results   []Result
}

type Result struct {
	Id       string
	Language string
	Path     string
	Url      string
	Name     string
	Html     string
}

func QueryTTCSearchResults(mock bool, uri string) (string, error) {
	if mock {
		return CACHED_SEARCH_RESULTS, nil
	}
	resp, err := http.Get(uri)
	if err != nil {
		fmt.Printf("Error in querying ttc api: %v", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("Error in parsing ttc api response body: %v", err)
		return "", err
	}
	return string(body), nil
}

func runTTCStage(mock bool) ([]Event, error) {
	properties, err := QueryTTCSubwayServiceLandingPageProperties(mock)
	if err != nil {
		return nil, err
	}
	var props Prop
	json.Unmarshal([]byte(properties), &props)
	url, err := props.GetURL()
	if err != nil {
		return nil, err
	}
	resp, err := QueryTTCSearchResults(mock, url)
	if err != nil {
		return nil, err
	}
	var results Results
	json.Unmarshal([]byte(resp), &results)
	var events []Event
	for _, r := range results.Results {
		event, err := extractEventInfo(r.Html)
		if err != nil {
			log.Printf("Skipping event: %v", err)
			continue
		}
		events = append(events, *event)
	}
	return events, nil
}

func extractEventInfo(html string) (*Event, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("unable to parse HTML result: %w", err)
	}
	href, exists := doc.Find("a").Attr("href")
	if !exists {
		return nil, fmt.Errorf("unable to find link")
	}
	title := doc.Find(".field-satitle").Text()
	route := strings.Title(doc.Find(".field-routename").Text())
	t := determineClosureType(title)
	start, end, err := parseEffectiveDatesFromSubtitle(doc.Find(".sa-effective-date").Text())
	if err != nil {
		return nil, fmt.Errorf("unable to parse effective date subtitle: %w", err)
	}
	return &Event{
		Summary: fmt.Sprintf("%s - %s", route, title),
		Uri:     fmt.Sprintf("https://ttc.ca%s", href),
		Type:    t,
		Start:   start,
		End:     end,
	}, nil
}

func parseEffectiveDatesFromSubtitle(subtitle string) (time.Time, time.Time, error) {
	nilTime := time.Now()
	re := regexp.MustCompile(`[A-Z][a-z]+ [\d]{1,2}, [\d]{4}`)
	matches := re.FindAllString(subtitle, -1)
	if len(matches) < 1 || len(matches) > 2 {
		return nilTime, nilTime, fmt.Errorf("unexpected amount of dates in subtitle: %s", subtitle)
	}
	start, err := time.Parse("January _2, 2006", strings.TrimSpace(matches[0]))
	if err != nil {
		return nilTime, nilTime, fmt.Errorf("unable to parse start date: %v", err)
	}
	if len(matches) == 1 {
		return start, start, nil
	}
	end, err := time.Parse("January _2, 2006", strings.TrimSpace(matches[1]))
	if err != nil {
		return nilTime, nilTime, fmt.Errorf("unable to parse end date: %v", err)
	}
	return start, end, nil
}

func determineClosureType(title string) ClosureType {
	if strings.Contains(title, "weekend closure") || isSingleDayClosure(title) {
		return FullDay
	} else if strings.Contains(title, "nightly early closures") {
		return NightOnly
	}
	return Undefined
}

func isSingleDayClosure(title string) bool {
	return strings.Contains(title, "single day closure")
}
