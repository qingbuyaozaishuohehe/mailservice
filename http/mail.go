package http

import (
	"strings"
	"io/ioutil"
	"net/http"
	"mailservice/email"
	"mailservice/config"
	"net/mail"
	"fmt"
)

func HttpMail(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	fmt.Println(r.MultipartForm.File)
	//fmt.Println(r.MultipartForm.File["attachment"])
	//for _, v := range r.MultipartForm.File["attachment"]{
	//	fmt.Println(v.Filename)
	//}
	return


	if r.MultipartForm == nil {
		Output(w, nil, http.StatusBadRequest, "获取POST数据失败")
		return
	}
	token, err := GetPostValue(w, r, "token", false)
	if err != nil{
		return
	}
	globalConfig := config.Get()
	if globalConfig.Http.Token != token{
		Output(w, nil, http.StatusForbidden, "Token验证失败")
		return
	}
	fromAddress, err := GetPostValue(w, r, "fromaddress", true)
	if err != nil{
		return
	}
	fromName, err := GetPostValue(w, r, "fromname", true)
	if err != nil{
		return
	}
	from := mail.Address{Name:fromName, Address:fromAddress}
	tos, err := GetPostValue(w, r, "tos", true)
	if err != nil{
		return
	}
	tos = strings.Replace(tos, ",", ";", -1)
	toList := strings.Split(tos, ";")
	if len(toList) <= 0{
		Output(w, nil, http.StatusBadRequest, "无效的收件人")
		return
	}
	ccs, _ := GetPostValue(w, r, "ccs", false)
	ccs = strings.Replace(ccs, ",", ";", -1)
	var ccList []string
	if len(ccs) > 0{
		ccList = strings.Split(ccs, ";")
		if len(ccList) <= 0{
			Output(w, nil, http.StatusBadRequest, "无效的抄送人")
			return
		}
	}
	subject, err := GetPostValue(w, r, "subject", true)
	if err != nil{
		return
	}
	content, err := GetPostValue(w, r, "content", true)
	if err != nil{
		return
	}
	bodyContentType, err := GetPostValue(w, r, "body_content_type", false)
	if err != nil{
		bodyContentType = "text/plain"
	}
	//获取附件
	fileList := r.MultipartForm.File["attachment"]
	fmt.Println(r.MultipartForm.File)
	return
	attachmentList := email.NewAttachmentList()
	if len(fileList) > 0{
		for _, file := range fileList{
			if len(file.Filename) <= 0{
				Output(w, nil, http.StatusBadRequest, "附件"+file.Filename+"没有名字")
				return
			}
			f, err := file.Open()
			if err != nil{
				Output(w, nil, http.StatusBadRequest, "打开附件"+file.Filename+"失败, err: " + err.Error())
				return
			}
			fileContent, err := ioutil.ReadAll(f)
			if err != nil{
				Output(w, nil, http.StatusBadRequest, "读取附件"+file.Filename+"失败, err: " + err.Error())
				return
			}
			attachment := email.NewAttachment(file.Filename, fileContent, false)
			attachmentList[file.Filename] = attachment
			if err != nil{
				Output(w, nil, http.StatusBadRequest, "发送附件"+file.Filename+"失败, err: " + err.Error())
				return
			}
		}
	}
	//准备发送
	m := email.New()
	m.From = from
	m.To = toList
	m.Cc = ccList
	m.Subject = subject
	m.Body = content
	m.BodyContentType = bodyContentType
	m.Attachments = attachmentList
	err = m.Send()
	if err != nil {
		Output(w, nil, http.StatusInternalServerError, err.Error())
		return
	}else{
		Output(w, "success", http.StatusOK, "")
		return
	}
}