package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dimonomid/clock"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func main() {
	now := clock.New().Now()
	from := mail.NewEmail("Example User", "sazajun125@juntendou1390.uk")
	to := mail.NewEmail("Example User", "sazajun125@gmail.com")
	content := mail.NewContent("text/html", " ")

	m := mail.NewV3MailInit(from, "test_template_sazajun125", to, content)
	m.SetTemplateID(os.Getenv("SENDGRID_TEMPLATE_ID"))
	m.Personalizations[0].SetDynamicTemplateData("app_name", "test_apps")
	m.Personalizations[0].SetDynamicTemplateData("code", "code")
	m.Personalizations[0].SetDynamicTemplateData("user_name", "sazajun125@gmail.com")
	m.Personalizations[0].SetDynamicTemplateData("expires_minutes", "10")
	m.Personalizations[0].SetDynamicTemplateData("support_email", "sazajun125@gmail.com")
	m.Personalizations[0].SetDynamicTemplateData("year", strconv.Itoa(now.Year()))

	//client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)

	client := &rest.Client{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				MaxIdleConns:          2,
				MaxIdleConnsPerHost:   2,
				IdleConnTimeout:       90 * time.Millisecond,
			},
			Timeout: 5 * time.Second,
		},
	}

	response, err := client.Send(request)
	if err != nil {
		log.Fatal("error sending email:", err.Error())
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		log.Fatal("an unexpected error occurred:", response.StatusCode, response.Body, response.Headers)
	}
	log.Print("email successfully sent...")
	/*
			{
		    "app_name":"test_apps",
		    "code":"code",
		    "user_name":"example@gamil.com",
		    "expires_minutes":"10",
		    "support_email":"sazajun125@gmail.com",
		    "year":"2025"
		}

		response, err := client.Send(message)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(response.StatusCode)
			fmt.Println(response.Body)
			fmt.Println(response.Headers)
		}
	*/
}
