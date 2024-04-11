package email

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ihezebin/oneness/email"
)

func TestEmail(t *testing.T) {
	err := Init(email.Config{
		Host:     email.HostQQMail,
		Port:     email.PortQQMail,
		Username: "ihezebin@qq.com",
		Password: "yrt*********xebsqzcabi***bug",
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	err = Client().Send(ctx, email.NewMessage().
		WithTitle("test").
		WithReceiver("86744316@qq.com").
		WithDate(time.Now()).
		WithSender("hezebin").
		WithHtml(`
			<html>
			<body>
				<h3 style="color:white;background-color:skyblue">
				"Hello World！This is a test mail！"
				</h3>
			</body>
			</html>
		`).
		WithAttach(email.NewAttach("test.txt", strings.NewReader("dsadsad"))),
	)

	if err != nil {
		t.Fatal(err)
	}

	t.Log("send success")
}
