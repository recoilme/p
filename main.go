package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	baseUri = "https://api.vk.com/method/"
	members = "groups.getMembers"
	gid     = "mailru"
	fields  = "sex,bdate,city,country,photo_50,photo_200_orig,photo_200,photo_400_orig,photo_max,photo_max_orig,has_mobile,contacts,connections,site,education,can_post,last_seen,relation"
	version = "&version=5.64"
	//тут 1000 поставь как отладишь
	count = "&count=1000" //max = 1000
)

var (
	defHeaders = make(map[string]string)
)

type Config struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

type Profiles struct {
	Response struct {
		Count int `json:"count"`
		Users []struct {
			UID          int    `json:"uid"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			Sex          int    `json:"sex"`
			Bdate        string `json:"bdate,omitempty"`
			City         int    `json:"city,omitempty"`
			Country      int    `json:"country,omitempty"`
			Photo50      string `json:"photo_50"`
			Photo200     string `json:"photo_200,omitempty"`
			PhotoMax     string `json:"photo_max"`
			Photo200Orig string `json:"photo_200_orig"`
			Photo400Orig string `json:"photo_400_orig,omitempty"`
			PhotoMaxOrig string `json:"photo_max_orig"`
			HasMobile    int    `json:"has_mobile,omitempty"`
			CanPost      int    `json:"can_post"`
			Site         string `json:"site,omitempty"`
			LastSeen     struct {
				Time     int `json:"time"`
				Platform int `json:"platform"`
			} `json:"last_seen,omitempty"`
			Hidden          int    `json:"hidden,omitempty"`
			Deactivated     string `json:"deactivated,omitempty"`
			Skype           string `json:"skype,omitempty"`
			Twitter         string `json:"twitter,omitempty"`
			Instagram       string `json:"instagram,omitempty"`
			Facebook        string `json:"facebook,omitempty"`
			FacebookName    string `json:"facebook_name,omitempty"`
			University      int    `json:"university,omitempty"`
			UniversityName  string `json:"university_name,omitempty"`
			Faculty         int    `json:"faculty,omitempty"`
			FacultyName     string `json:"faculty_name,omitempty"`
			Graduation      int    `json:"graduation,omitempty"`
			EducationForm   string `json:"education_form,omitempty"`
			EducationStatus string `json:"education_status,omitempty"`
			Relation        int    `json:"relation,omitempty"`
			RelationPartner struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
			} `json:"relation_partner,omitempty"`
			HomePhone string `json:"home_phone,omitempty"`
		} `json:"users"`
	} `json:"response"`
}

func init() {
	defHeaders["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:52.0) Gecko/20100101 Firefox/52.0"
	defHeaders["Accept-Language"] = "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3"
	defHeaders["Referer"] = "https://ya.ru/"
	defHeaders["Cookie"] = ""
}

type response struct {
	resp *http.Response
	url  string
}

func get(url string, r chan response) {
	if resp, err := http.Get(url); err == nil {
		r <- response{resp, url}
	}
}

func main() {
	fmt.Printf("user.UID\tuser.FirstName\tuser.LastName\tuser.Photo200Orig\tuser.Sex\tuser.Bdate\tuser.City\tuser.Facebook\tuser.FacebookName\tuser.HomePhone\tuser.Instagram\tuser.Skype\tuser.Twitter\tuser.UniversityName\tuser.EducationStatus")
	//see https://vk.com/dev/groups.getMembers
	urlmask := baseUri + members + "?group_id=" + gid + "&fields=" + url.QueryEscape(fields) + version
	var offset = 0
	//это бесконечн цикл в гоу
	for {
		url := urlmask + count + "&offset=" + strconv.Itoa(offset)
		offset = offset + 1000
		//if offset == 100 {
		//break
		//}

		b := HttpGet(url, nil)
		//log.Println(string(b))
		//use https://mholt.github.io/json-to-go/ Luke!
		if b == nil {
			log.Println("empty")
			break
		}
		var res Profiles
		err := json.Unmarshal(b, &res)
		if err != nil {
			log.Println("Error", url, err)
			break
		}
		items := res.Response.Users
		if items == nil || len(items) == 0 {
			log.Println("no items")
			//больше нет юзеров
			break
		}
		for _, user := range items {
			fmt.Printf("%d\t%s\t%s\t%s\t%d\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n", user.UID, user.FirstName, user.LastName, user.Photo200Orig, user.Sex,
				user.Bdate, user.City, user.Facebook, user.FacebookName, user.HomePhone, user.Instagram, user.Skype, user.Twitter,
				user.UniversityName, user.EducationStatus)
		}
		//закоменть брик как будешь готов по формату
		//

		//break
		//vk не любит запросы больше 3 в секунду
		time.Sleep(900 * time.Millisecond)

	}
}

func TimeoutDialer(config *Config) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(config.ReadWriteTimeout))
		return conn, nil
	}
}

func NewTimeoutClient(args ...interface{}) *http.Client {
	// Default configuration
	config := &Config{
		ConnectTimeout:   7 * time.Second,
		ReadWriteTimeout: 7 * time.Second,
	}

	// merge the default with user input if there is one
	if len(args) == 1 {
		timeout := args[0].(time.Duration)
		config.ConnectTimeout = timeout
		config.ReadWriteTimeout = timeout
	}

	if len(args) == 2 {
		config.ConnectTimeout = args[0].(time.Duration)
		config.ReadWriteTimeout = args[1].(time.Duration)
	}

	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(config),
		},
	}
}

// HttpGet create request with default headers + custom headers
func HttpGet(url string, headers map[string]string) []byte {
	//log.Println("httpGet", url)

	client := NewTimeoutClient()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	for k, v := range defHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	} else {
		return body
	}

	return nil
}
