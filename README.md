# 邮件发送服务

### 功能
能发邮件、抄送、HTML内容、附件等

### 起因
项目需要发送邮件给公司人员

### 参考
- github.com/open-falcon/mail-provider

### 依赖
- Linux下的sendmail

### 原理
使用SMTP服务发送。

### 用法

- Protocol: HTTP

- Method: POST

- 收件人: 必填, tos, 多个用逗号分隔

- 抄送人: 选填, ccs, 多个用逗号分隔

- 主题: 必填, subject

- 内容: 必填, content, 可用HTML代码

- 附件: 选填, attachment, 可单文件可数组

### 实例
##### FORM表单
```html
<form method="post" action="http://****:4000/sender/mail" enctype="multipart/form-data">
    <input type="text" name="tos" value="***@qq.com"><br>
    <input type="text" name="ccs" value="***@163.com"><br>
    <input type="text" name="subject" value="搭建邮件服务器-测试附件发送"><br>
    <input type="file" name="attachment"><br>
    <input type="file" name="attachment"><br>
    <textarea name="content">
<p>测试邮件发送服务</p>
<p>附件、抄送等功能</p>
<img src="http://*****.png">
    </textarea><br>
    <input type="submit">
</form>
```