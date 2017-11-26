package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	SaveDir = `D:\VMShared\` //文件上传与文件列表显示路径
)

var (
	fileSvr = http.FileServer(http.Dir(SaveDir))
)

//显示上传页面
func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("rootHandler:", r.Method, r.URL.Path)

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tpl := template.Must(template.ParseFiles("tmpl/upload.html"))
	tpl.Execute(w, nil)
}

//上传到服务器的文件的处理
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("uploadHandler:", r.Method, r.URL.Path)

	if r.Method == http.MethodPost {
		r.ParseMultipartForm(2 << 20)

		log.Println("token:", r.FormValue("token"))
		input, info, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println("Error:", err)
			return
		}
		defer input.Close()

		output, err := os.OpenFile(SaveDir+info.Filename, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		defer output.Close()

		io.Copy(output, input)
		http.Redirect(w, r, "/list", http.StatusFound)
	}
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.Handle("/list/", http.StripPrefix("/list/", http.FileServer(http.Dir(SaveDir))))

	log.Fatal(http.ListenAndServe(":1234", nil))
}
