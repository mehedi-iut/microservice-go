# Microservice in Go
## Multipart data handling

### What is multipart data in HTTP?
Multipart data is a type of data format that can be used in HTTP (Hypertext Transfer Protocol) messages to send multiple types of data as a single message. It is typically used to send binary data, such as images or files, along with textual data in the same HTTP request or response.
Multipart data consists of multiple parts, where each part is separated by a boundary string. The boundary string is a unique sequence of characters that marks the beginning and end of each part of the message.
Each part of the multipart data can have its own content type and encoding, allowing different types of data to be included in the same message. For example, a single HTTP request could include a text message, an image file, and a PDF document, all in multipart data format.
Multipart data is commonly used in web applications to upload files or submit form data that includes both textual and binary data. It is specified in the MIME (Multipurpose Internet Mail Extensions) standard, which is used for encoding and exchanging different types of data over the internet.

From previous episode we use simple fileserver by declaring handlerfunc
```go
ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadREST)
```
we change the handlerFunction from **ServeHTTP** to **UploadREST** to distinguish between simple fileserver and multipart
For multipart request we don't need url parameter

```go
ph.HandleFunc("/", fh.UploadMultipart)
```
now because multi-part the request is actually gonna have the file name and it's gonna have the ID and all of that sort of stuff in it we don't need to to necessarily capture that here so we're just gonna have a post which goes to the root path("/")

now in the handler we need to define **UploadMultipart** function
```go
// UploadMultipart data
func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	id, idErr := strconv.Atoi(r.FormValue("id"))
	f.log.Info("Process form for id", "id", id)

	if idErr != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected integer id", http.StatusBadRequest)
		return
	}

	ff, mh, err := r.FormFile("file")
	if err != nil {
		f.log.Error("Bad request", "error", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
		return
	}

	f.saveFile(r.FormValue("id"), mh.Filename, rw, ff)
}
```
here we first call the ```r.parseMultipartForm(128 * 1024)``` which read all the information in the multipart request
and parse it to key value pair.
multi-part request is going to have a body which is containing all of this information so does that mean 
that we have to parse out all of these boundaries ourselves? not too much. there go has a bunch of stuff for doing this
on the the HTTP request what you actually have is two methods **parseForm** for dealing with sort of URL encoded forms and
you have **parseMultipartForm** so we can use these methods so parse multi-part form what is that going to 
do is, it's gonna read all of that multi-part form data for us it's gonna
automatically start breaking up the body on the boundary and it's going to store it for us in to some end to 
a collection that we can later reference

```id, idErr := strconv.Atoi(r.FormValue("id"))``` in this code we used *request* **FormValue** field to retrieve the 
value from the multi part form data [link](https://pkg.go.dev/net/http#Request.FormValue)
and converted into integer using **strconv.Atoi()**

```go
ff, mh, err := r.FormFile("file")
```
again, we use *request* **FormFile** key to get the file [link](https://pkg.go.dev/net/http#Request.FormFile)
here **"id"** and **"file"** is coming from frontend field definition.
```javascript
const data = new FormData()
        data.append('file', this.state.file);
        data.append('id', this.state.id);
```
It is in the path *frontend/src/Admin.js*
then we save our file using ```f.saveFile(r.FormValue("id"), mh.Filename, rw, ff)```
now, in our **saveFile** function, we need to change the input argument, from ```r *http.Request``` to ```r io.ReadCloser```
For this change, we need to change the **UploadREST**  **saveFile** 
```f.saveFile(id, fn, rw, r.Body)```, we just change from **r** to **r.Body** as r.Body is **io.ReadCloser**
