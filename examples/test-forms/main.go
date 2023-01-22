package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/lightsaid/gotk/forms"
)

const templateForm = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    <form action="/form" method="post">
        nickname: <input type="text" name="nickname"> <br>
        email: <input type="email" name="email"> <br>
        phone: <input type="text" name="phone"> <br>
        age: <input type="number" name="age"> <br>
        password: <input type="password" name="password"> <br>
        hobby: <br>
            <input type="checkbox" value="foo" name="hobby"> foo
            <input type="checkbox" value="bar" name="hobby"> boo
            <input type="checkbox" value="xyz" name="hobby"> xyz <br>
		action: <br>
            <input type="radio" value="codeing" name="action">codeing <br>
            <input type="radio" value="playing" name="action">playing <br>
        <br>
        选择：
        <select name="yn" id="yn">
            <option value="Y">Yes</option>
            <option value="N">No</option>
        </select>
        <br>
        介绍：<textarea name="introduce" cols="30" rows="3"></textarea>
        <br>
        <button type="submit">提交</button>

        <div style="padding: 20px; color: red;">
			<textarea id="code" readonly cols="100" rows="30" style="color: red;"></textarea>
        </div>
    </form>
    <script>
        let str = "{{.}}"
        document.querySelector("#code").innerHTML = str
    </script>
</body>
</html>
`

func main() {
	var addr = "0.0.0.0:9527"
	mux := http.NewServeMux()

	mux.Handle("/form", http.HandlerFunc(HandleForm))

	fmt.Println("server starting on: ", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func HandleForm(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("formTpl").Parse(templateForm)
	if r.Method == http.MethodGet {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.ExecuteTemplate(w, "formTpl", nil)
	} else {
		r.ParseForm()
		v := forms.New(r.Form)
		fmt.Println(r.Form)
		v.Required("nickname", "email", "password", "hobby", "action", "yn")

		// 是否邮箱
		v.IsEmail("email", "不是邮箱")

		// 是否手机
		v.IsPhone("phone", "不是手机号")

		// 最大最小值
		v.MinLength("nickname", 2)
		v.MaxLength("nickname", 8)

		// check
		age, _ := strconv.Atoi(v.Get("age"))
		// 条件满足 18 岁，如果不满足 则 添加 ‘未成年’
		v.Check(age >= 18, "age", "未成年")

		fmt.Println(">>>> ", v.Matches(v.Get("phone"), regexp.MustCompile(`^1[3-9]\d{9}$`)))

		v.RequiredForMsg("hobby", "hobby 必须选择")

		if !v.Valid() {
			b, _ := json.MarshalIndent(v.Errors, "", "\t")
			t.ExecuteTemplate(w, "formTpl", string(b))
		} else {
			t.ExecuteTemplate(w, "formTpl", "Success")
		}
	}
}
