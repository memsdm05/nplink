package provider

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
	"github.com/memsdm05/nplink/util"
	"github.com/zellyn/kooky"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	_ "github.com/memsdm05/nplink/provider/nightbot"
)

var chatbots = map[string]Provider{}

type Provider interface {
	AuthURL() string
	SetSession(session string)
	CheckSession() bool
	ResolveSession(vals url.Values) (string, error)

	SetCommand(name, msg string, extra string)
	DeleteCommand()
}

func Register(name string, provider Provider)  {
	chatbots[name] = provider
}


type Page struct {
	user string
	app  string

	scopes []Scope

	width int

	form   url.Values
	client http.Client
}
type severity int

const (
	high = severity(iota)
	medium
	low
)

type Scope struct {
	severity
	text string
}

func must(e error) {
	if e != nil {
		panic(e)
	}
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

func scrape(r io.Reader, p *Page) {
	doc, _ := goquery.NewDocumentFromReader(r)

	p.app = strings.TrimSpace(
		doc.Find("div.authorize_prompt h1 strong").Text())

	p.user = strings.TrimSpace(
		doc.Find("p.user-info__username strong").Text())

	doc.Find("ul.high_severity li").Each(func(_ int, s *goquery.Selection) {
		p.scopes = append(p.scopes, Scope{
			high,
			s.Find("span").Text(),
		})
	})

	doc.Find("ul.medium_severity li").Each(func(_ int, s *goquery.Selection) {
		p.scopes = append(p.scopes, Scope{
			medium,
			s.Find("span").Text(),
		})
	})

	doc.Find("ul.low_severity li").Each(func(_ int, s *goquery.Selection) {
		p.scopes = append(p.scopes, Scope{
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

func NewPage(prov Provider) *Page {
	p := new(Page)
	p.form = url.Values{}

	c, _ := cookiejar.New(nil)
	p.client.Jar = c
	p.client.Jar.SetCookies(collectCookies())
	p.client.CheckRedirect = util.StopRedirect

	resp, err := p.client.Get(prov.AuthURL())
	must(err)
	defer resp.Body.Close()

	scrape(resp.Body, p)

	for _, scope := range p.scopes {
		if len(scope.text) > p.width {
			p.width = len(scope.text)
		}
	}

	return p
}

func (p *Page) Display() {
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

func (p *Page) Authorize(prov Provider) (string, error) {
	resp, _ := p.client.PostForm("https://id.twitch.tv/oauth2/authorize", p.form)
	defer resp.Body.Close()

	params, parseErr := url.Parse(resp.Header.Get(	"location"))
	if parseErr != nil {
		return "", parseErr
	}

	sess, resolveErr := prov.ResolveSession(params.Query())
	return sess, resolveErr
}