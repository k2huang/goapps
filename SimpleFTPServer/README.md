## 简易FTP服务器

### 1.基本功能
`/` : 显示上传页面<br/>
`/list` : 显示服务器上的文件列表<br/>
`/upload` : 上传文件的请求路径<br/>


### 2.说明
net/http包中有一段示例代码：<br/>
```go
// To serve a directory on disk (/tmp) under an alternate URL
// path (/tmpfiles/), use StripPrefix to modify the request
// URL's path before the FileServer sees it:
http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp")))) 
```
当你通过 /tmpfiles 访问服务器的时候可以显示 服务器上 `/tmp` 下的文件列表。<br/>

`本程序通过结合 上述示例代码 和 标准库中上传文件的API，实现了一个 简易FTP服务器。`