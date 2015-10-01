package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jimmykuu/gt-go-sdk"
)

const page = `<!doctype html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html;charset=UTF-8">
	<title>极验行为式验证 Golang 类网站安装测试页面</title>
</head>
<body>
	<style type="text/css">
		.container{
			width: 960px;
			margin: 0 auto;
		}
		.content{
			width: 960px;
			margin: 10 auto;
			border-top: 1px solid #ccc;
		}
		.box{
			width:300px;
			margin: 30px auto;
		}
		.header{
			margin: 80px auto 30px auto;
			text-align: center;
			font-size: 34px;
		}
		input{
			width: 200px;
			padding: 6px 9px;
		}
		button{
			cursor: pointer;
			line-height: 35px;
			width: 110px;
			margin:30px 0 0 90px;
			border: 1px solid #FFFFF0;
			background-color: #31C552;
			border-radius: 4px;
			font-size: 14px;
			color: #FFFFF0;
		}
	</style>

	<div class="container">
		<div class="header">
			极验行为式验证 Golang 类网站安装测试页面
		</div>
		<div class="content">
			<form method="post" action="/login">
				<div class="box">
					<label>邮箱：</label>
					<input type="text" name="email" value="geetest@geetest.com"/>
				</div>
				<div class="box">
					<label>密码：</label>
					<input type="password" name="password" value="geetest"/>
				</div>
                <div class="box">
                    <script type="text/javascript" src="{{.url}}"></script>
                </div>
				<div class="box">
					<button id="submit-button">提交</button>
				</div>
			</form>
		</div>
	</div>
</body>
</html>
`

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		geeTest := geetest.NewGeeTest("", "")

		if r.Method == "POST" {
			challenge := r.FormValue("geetest_challenge")
			validate := r.FormValue("geetest_validate")
			seccode := r.FormValue("geetest_seccode")
			result := geeTest.Validate(challenge, validate, seccode)
			if result {
				fmt.Fprintf(w, "success")
			} else {
				fmt.Fprint(w, "fail")
			}

			return
		}

		t, err := template.New("login").Parse(page)

		if err != nil {
			log.Fatal(err)

			return
		}

		challenge := geeTest.Challenge()
		url := geeTest.EmbedURL(challenge)

		err = t.Execute(w, map[string]interface{}{
			"url": url,
		})

		if err != nil {
			log.Fatal(err)
		}
	})

	log.Println("Server start on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
