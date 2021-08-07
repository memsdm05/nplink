package setup

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/memsdm05/nplink/provider"
	"github.com/memsdm05/nplink/utils"
	"github.com/zellyn/kooky"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/zellyn/kooky/allbrowsers"
)

func Auth() {
	prov := SelectedProvider
	session, success := utils.GetCred(prov.Name())

	if !success {
		session = authFlow(prov)
	}

	err := prov.Session(session)

	if err == provider.BadSessionErr {
		if !success {
			fmt.Println("Something happened, please redo the authorization")
		}
		session = authFlow(prov)
		utils.Must(prov.Session(session))
	}

	utils.SetCred(prov.Name(), session)
}

func authFlow(prov provider.Provider) string {
	p := NewPage(prov)
	p.Display()

	fmt.Printf("Do you authorize %s to access your account (%s)? [Y/n]\n",
		p.app,
		p.user)

	var s string

outer:
	for !Config.SkipAuth {
		fmt.Print("> ")
		fmt.Scanln(&s)

		switch strings.ToLower(s) {
		case "yes", "y", "":
			break outer
		case "no", "n":
			fmt.Println("Because you selected \"no\" the application will simply close in 5 seconds\n")
			fmt.Println("goodbye.")
			time.Sleep(5 * time.Second)
			os.Exit(1)
		default:
			fmt.Println("Invalid response, please try again")
		}
	}

	return p.Authorize(prov)
}

type page struct {
	user string
	app  string

	scopes []scope

	width int

	form   url.Values
	client *utils.Client
}
type severity int

const (
	high = severity(iota)
	medium
	low
)

type scope struct {
	severity
	text string
}

func collectCookies() (*url.URL, []*http.Cookie) {
	base, _ := url.ParseRequestURI("https://id.twitch.tv")

	httpCookies := make([]*http.Cookie, 0)

	for _, cookie := range kooky.ReadCookies(kooky.Valid, kooky.Domain(".twitch.tv")) {
		httpCookie := cookie.HTTPCookie()
		httpCookies = append(httpCookies, &httpCookie)
	}
	//fmt.Println(base, httpCookies)
	return base, httpCookies
}

func (p *page) scrape(r io.Reader) {
	doc, _ := goquery.NewDocumentFromReader(r)
	//fmt.Println(doc.Html())
	p.app = strings.TrimSpace(
		doc.Find("div.authorize_prompt h1 strong").Text())

	p.user = strings.TrimSpace(
		doc.Find("p.user-info__username strong").Text())

	doc.Find("ul.high_severity li").Each(func(_ int, s *goquery.Selection) {
		p.scopes = append(p.scopes, scope{
			high,
			s.Find("span").Text(),
		})
	})

	doc.Find("ul.medium_severity li").Each(func(_ int, s *goquery.Selection) {
		p.scopes = append(p.scopes, scope{
			medium,
			s.Find("span").Text(),
		})
	})

	doc.Find("ul.low_severity li").Each(func(_ int, s *goquery.Selection) {
		p.scopes = append(p.scopes, scope{
			low,
			s.Text(),
		})
	})

	doc.Find("form#authorize_form :input[type='hidden']").Each(func(_ int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		value, _ := s.Attr("value")
		p.form.Add(name, value)
	})
	p.form.Set("cancel", "false")
}

func NewPage(prov provider.Provider) *page {
	p := new(page)
	p.form = url.Values{}

	p.client = utils.NewClient()
	p.client.Jar.SetCookies(collectCookies())
	p.client.CheckRedirect = utils.StopRedirect

	resp, err := p.client.Get(prov.URL())
	//b, _ := io.ReadAll(resp.Body)
	//fmt.Println(string(b))

	defer resp.Body.Close()
	utils.Must(err)

	p.scrape(resp.Body)

	for _, scope := range p.scopes {
		if len(scope.text) > p.width {
			p.width = len(scope.text)
		}
	}

	return p
}

func (p *page) Display() {
	fmt.Println(p.app, "wants to access your account,", p.user)
	fmt.Println("Authorizing will allow", p.app, "to: ")
	for _, scope := range p.scopes {
		switch scope.severity {
		case high:
			fmt.Print("⚠\t")
			color.Set(color.Bold)
		case medium:
			fmt.Print("⚠\t")
		case low:
			fmt.Print(" \t")
			color.Set(color.Italic)
		}
		fmt.Println(scope.text)
		color.Set(color.Reset)
	}
}

func (p *page) Authorize(prov provider.Provider) string {
	resp, _ := p.client.PostForm("https://id.twitch.tv/oauth2/authorize", p.form)
	defer resp.Body.Close()

	params, _ := url.Parse(resp.Header.Get("location"))
	sess, _ := prov.ResolveSession(params.Query())
	return sess
}
