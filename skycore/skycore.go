package main

import (
	//"bufio"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/wgerlach/Skycore/skycore/go-etcd-0.4/etcd"
	"github.com/wgerlach/Skycore/skycore/vendor/github.com/MG-RAST/AWE/lib/shock"
	"github.com/wgerlach/Skycore/skycore/vendor/github.com/MG-RAST/golib/go-uuid/uuid"
	"github.com/wgerlach/Skycore/skycore/vendor/github.com/fsouza/go-dockerclient"
	"io"
	"mime/multipart"
	"net/http" // should all be done by shock lib
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var shock_url_default = os.Getenv("SKYCORE_SHOCK")
var shock_token_default = os.Getenv("SKYCORE_SHOCK_TOKEN")
var etcd_urls_string_default = "http://127.0.0.1:4001"
var docker_socket_default = "unix:///var/run/docker.sock"

var flags *flag.FlagSet

type Skycore struct {
	Shock_client  shock.ShockClient
	Docker_client *docker.Client
	Etcd_client   *etcd.Client
}

type Dockerimage_attributes struct {
	//Temporary  string `json:"temporary"`
	Type       string `json:"type"`       // required= "dockerimage"
	Id         string `json:"id"`         // required
	Name       string `json:"name"`       // required
	Repository string `json:"repository"` // required
	Tag        string `json:"tag"`        // required
	//Docker_version interface{}   `json:"docker_version"` // deprecataed, better use complete Image
	Base_image_tag string        `json:"base_image_tag"`
	Base_image_id  string        `json:"base_image_id"`
	Dockerfile     string        `json:"dockerfile"`
	Image          *docker.Image `json:"image"` // warning: size cannot be stored in AND retrieved from JSON (int64 vs float64)
}

// shock.metagenomics.anl.gov

// git pull ; ./compile.sh && sudo ./skyc load --shock shock.metagenomics.anl.gov 7c10f6dd-5291-45fe-a938-e2ae2027482a

// on wolfgang
// scp -i ~/.ssh/wo_magellan_pubkey.pem ubuntu@140.221.67.161:/home/ubuntu/gopath/bin/skycore . ; scp -i ~/.ssh/wo_magellan_pubkey.pem ./skycore core@140.221.67.208:

func NumberToString(n int64, sep rune) string {
	s := strconv.FormatInt(n, 10)
	//s := strconv.Itoa(n)

	startOffset := 0
	var buff bytes.Buffer

	if n < 0 {
		startOffset = 1
		buff.WriteByte('-')
	}

	l := len(s)

	commaIndex := 3 - ((l - startOffset) % 3)

	if commaIndex == 3 {
		commaIndex = 0
	}

	for i := startOffset; i < l; i++ {

		if commaIndex == 3 {
			buff.WriteRune(sep)
			commaIndex = 0
		}
		commaIndex++

		buff.WriteByte(s[i])
	}

	return buff.String()
}

func CopyTicker(writer io.Writer, reader io.Reader, name string) (written int64, err error) {

	var read int64
	read = 0

	ticker := time.NewTicker(time.Second * 3)
	go func() {
		var previous int64
		previous = 0
		for t := range ticker.C {
			//fmt.Println("Tick at", t)
			rate := (read - previous) / 3
			fmt.Printf("%s: Bytes submitted: %s    Bytes per second: %s   (%s)\n", name, NumberToString(read, ','), NumberToString(rate, ','), t) // \r or \n
			previous = read
		}
	}()

	//var p float32
	for {
		buffer := make([]byte, 1024)
		cBytes, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		read = read + int64(cBytes)
		//fmt.Printf("read: %v \n", read)
		//p = float32(read) / float32(length) * 100
		//fmt.Printf("progress: %v \n", p)
		writer.Write(buffer[0:cBytes])
	}
	written = read
	ticker.Stop()
	fmt.Printf("\n")
	//read, err := io.Copy(part, reader)

	fmt.Fprintf(os.Stdout, fmt.Sprintf("%s: bytes submitted: %s\n", name, NumberToString(read, ',')))
	return
}

// Creates a new file upload http request with optional extra params
func newStreamUploadRequest(uri string, params map[string]string, paramName, path string, reader io.Reader) (req *http.Request, err error) {
	// with modifications taken from:
	// http://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/

	//file, err := os.Open(path)
	//if err != nil {
	//	return nil, err
	//}
	//defer file.Close()

	//body := &bytes.Buffer{}

	bodyReader, bodyWriter := io.Pipe()

	multiWriter := multipart.NewWriter(bodyWriter)

	errChan := make(chan error, 1) //TODO read channel

	go func() {
		defer bodyWriter.Close()

		part, err := multiWriter.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			errChan <- err
			return
		}
		fmt.Fprintf(os.Stdout, "Start copying into form\n")

		// this stuff below inlcuding for loop is only for progres information
		_, err = CopyTicker(part, reader, "form")
		if err != nil {
			errChan <- err
			return
		}
		for k, v := range params {
			if err := multiWriter.WriteField(k, v); err != nil {
				errChan <- err
				return
			}
		}
		errChan <- multiWriter.Close()
	}()

	req, err = http.NewRequest("POST", uri, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", multiWriter.FormDataContentType())

	return
}

func (skyc *Skycore) ExportImageNonBlocking(name string) (*io.PipeReader, error) {

	image_raw_reader, image_raw_writer := io.Pipe()

	go func() {
		defer image_raw_writer.Close()

		fmt.Fprintf(os.Stdout, "Start export from docker engine\n")
		// version that works: echo -e "GET /images/ubuntu:14.04/get HTTP/1.0\r\n" | nc -U /var/run/docker.sock // with some http header lines

		opts := docker.ExportImageOptions{Name: name, OutputStream: image_raw_writer} // this is an io.Writer
		err := skyc.Docker_client.ExportImage(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("error: ", err.Error(), "\n")) // TODO need channel to stop here or close reader

		} else {
			fmt.Fprintf(os.Stdout, "finished exporting image from docker engine\n")
		}
		return
	}()

	return image_raw_reader, nil
}

func (skyc *Skycore) GzipNonBlocking(reader io.Reader) (*io.PipeReader, error) {

	// go routine return values are written to image_gzipped_reader
	image_gzipped_reader, image_gzipped_writer := io.Pipe() //  final gzip result can be read from reader, gzip writes into writer
	fmt.Fprintf(os.Stdout, "Start compressing stream\n")
	gzip_writer := gzip.NewWriter(image_gzipped_writer) // return a *Writer ; A compress/gzip.Writer is an io.WriteCloser. ; expects io.Writer
	//defer gzip_writer.Close()
	go func() {
		//defer image_gzipped_writer.Close()
		//defer gzip_writer.Close()

		//	 is a writeCloser and needs to be closed.
		fmt.Fprintf(os.Stdout, "Start gzipping\n")

		//n, err := CopyTicker(gzip_writer, reader, "gzip")
		n, err := io.Copy(gzip_writer, reader) // target, source // image_raw_reader
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("error writing to gzip: %s\n", err.Error()))
		}

		gzip_writer.Close()
		image_gzipped_writer.Close()

		fmt.Fprintf(os.Stdout, "gzipped %d bytes, closing gzip go routine\n", n)
	}()

	return image_gzipped_reader, nil
}

func (skyc *Skycore) save_image_to_shock(name string, private_image bool) (node string, err error) {

	name_array := strings.Split(name, ":")
	if len(name_array) != 2 {
		return "", errors.New("error: image name has to have format \"repository:tag\"")
	}
	repository := name_array[0]
	tag := name_array[1]

	if skyc.Docker_client == nil {
		return "", errors.New("error: skyc.Docker_client==nil")
	}

	if skyc.Shock_client.Host == "" {
		return "", errors.New("error: skyc.Shock_client.Host not defined")
	}


	// inspect
	image_obj, err := skyc.Docker_client.InspectImage(name)
	image_id := ""
	if err != nil {
		// image not found
		return
	} else {
		image_id = image_obj.ID
		fmt.Fprintf(os.Stdout, fmt.Sprintf("found image_id: %s\n", image_id))
	}


	if skyc.Shock_client.Token == "" {
		fmt.Fprintf(os.Stdout, "Please provide Shock token (or use option --token):\n")
		var user_token string
		_, err = fmt.Scanf("%s\n", &user_token)
		if err != nil {
			return
		}
		skyc.Shock_client.Token = user_token
	}


	// *** export (save) image from docker engine
	image_raw_reader, err := skyc.ExportImageNonBlocking(image_id)
	if err != nil {
		return
	}

	// *** gzip the image stream
	image_gzipped_reader, err := skyc.GzipNonBlocking(image_raw_reader)
	if err != nil {
		return
	}

	image_obj.Size = 0 // ugly workaround !
	attr_struct := Dockerimage_attributes{

		Id:         image_id,
		Type:       "dockerimage",
		Name:       name,
		Repository: repository,
		Tag:        tag,
		//Docker_version: version.Map(),
		Image: image_obj,
	}

	// other keys used before
	//Base_image_tag 	string
	//Base_image_id 	string
	//Dockerfile 		string

	attr_string, _ := json.Marshal(attr_struct)
	fmt.Fprintf(os.Stdout, "attr_string: %s\n", attr_string)
	//os.Exit(1)

	// make it it mulitpart stream
	params := map[string]string{
		"attributes_str": string(attr_string[:]),
	}

	fmt.Fprintf(os.Stdout, "Mulitpart stuff\n")

	// *** create the http request object
	upload_url := skyc.Shock_client.Host + "/node"
	fmt.Fprintf(os.Stdout, "upload_url: %s\n", upload_url)

	valid_char := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '_' || r == '-' || r == '.':
			return r
		}
		return '_'
	}
	shock_filename := strings.Map(valid_char, image_id+"_"+repository+"."+tag) + ".tar.gz"

	request, err := newStreamUploadRequest(upload_url, params, "upload", shock_filename, image_gzipped_reader)
	if err != nil {
		return
	}

	

	if skyc.Shock_client.Token != "" {
		fmt.Fprintf(os.Stdout, "using token\n")
		request.Header.Add("Authorization", "OAuth "+skyc.Shock_client.Token)
	}
	if err != nil {
		return
	}

	client := &http.Client{}
	fmt.Fprintf(os.Stdout, "Do request\n")
	resp, err := client.Do(request)
	if err != nil {
		return
	}

	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header)
	fmt.Println(body)

	response := new(shock.ShockResponse)
	if err := json.Unmarshal(body.Bytes(), response); err != nil {
		return "", err
	}
	if len(response.Errs) > 0 {
		return "", errors.New(strings.Join(response.Errs, ","))
	}
	new_node_obj := &response.Data
	if new_node_obj == nil {
		err = errors.New("empty node got from Shock")
		return
	}

	node = new_node_obj.Id

	if !private_image {
		_, err = skyc.Shock_client.Make_public(node)
		if err != nil {
			fmt.Fprintf(os.Stderr, "There was an error making the node public: "+err.Error()+"\n")
		}
	}

	fmt.Fprintf(os.Stdout, "New shock node: "+upload_url+"/"+node+"\n")
	fmt.Fprintf(os.Stdout, "example: curl -L http://127.0.0.1:4001/v2/keys/service_images/<servicename>/shock -XPUT -d value=\"%s\"\n", upload_url+"/"+node)
	return
}

func dockerLoadImage(client *docker.Client, download_url string, datatoken string) (err error) {

	if client == nil {
		return errors.New(fmt.Sprintf("Error: docker client not initialized "))
	}

	image_stream, err := shock.FetchShockStream(download_url, datatoken) // token empty here, assume that images are public
	if err != nil {
		return errors.New(fmt.Sprintf("Error getting Shock stream, err=%s, url=%s", err.Error(), download_url))
	}

	gr, err := gzip.NewReader(image_stream) //returns (*Reader, error) // TODO not sure if I have to close gr later ?
	if err != nil {
		return
	}

	fmt.Fprintf(os.Stdout, "loading image... (%s)\n", download_url)

	//var buf bytes.Buffer

	opts := docker.LoadImageOptions{InputStream: gr}
	err = client.LoadImage(opts)
	//err = client.LoadImage(gr, &buf) // in io.Reader, w io.Writer

	if err != nil {
		return errors.New(fmt.Sprintf("Error loading image, err=%s", err.Error()))
	}
	//logger.Debug(1, fmt.Sprintf("load image returned: %v", &buf))

	return
}

func get_attribute_string(attr_map map[string]interface{}, key string) (value string, err error) {

	value_interface, ok := attr_map[key]
	if !ok {
		return "", errors.New(fmt.Sprintf("error: did not find key %s in attributes", key))
	}
	value, ok = value_interface.(string)
	if !ok || value == "" {
		return "", errors.New(fmt.Sprintf("error: could not parse value of key %s in attributes", key))
	}
	return
}

func (skyc *Skycore) get_dockerimage_shocknode_attributes(node_id string) (image_repository string, image_tag string, image_id string, err error) {

	node_response, err := skyc.Shock_client.Get_node(node_id)

	if err != nil {
		err = errors.New(fmt.Sprintf("error getting shock node: ", err.Error()))
		return

	}

	if len(node_response.Errs) > 0 {
		err = errors.New(fmt.Sprintf("error getting shock node: %s", strings.Join(node_response.Errs, ",")))
	}

	//docker_attr := node_response.Data.Attributes.(Dockerimage_attributes)

	attr_json, err := json.Marshal(node_response.Data.Attributes) // ugly hack to get attributes (type: map[string]interface{}) into struct

	var docker_attr Dockerimage_attributes
	err = json.Unmarshal(attr_json, &docker_attr)
	if err != nil {
		return
	}

	image_repository = docker_attr.Repository
	image_tag = docker_attr.Tag
	image_id = docker_attr.Id

	//attr_map, ok := node_response.Data.Attributes.(map[string]interface{}) // is of type map[string]interface{}

	//if !ok {
	//	err = errors.New(fmt.Sprintf("error: could not acces node attributes"))
	//	return
	//}
	//image_repository, err = get_attribute_string(attr_map, "repository")
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "error reading repository from shock node\n")
	//		return
	//	}
	//	image_tag, err = get_attribute_string(attr_map, "tag")
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "error reading tag from shock node\n")
	//		return
	//	}

	//	image_id, err = get_attribute_string(attr_map, "id")
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "error reading id from shock node\n")
	//		return
	//	}

	if len(image_id) != 64 {
		err = errors.New(fmt.Sprintf("error: image_id is not 64 characters long"))
	}

	_, err = hex.DecodeString(image_id)
	if err != nil {
		err = errors.New(fmt.Sprintf("error: image_id is not a valid 64 hexadecimal digit string"))
		return
	}

	fmt.Fprintf(os.Stdout, "image_repository: "+image_repository+"\n")
	fmt.Fprintf(os.Stdout, "image_tag: "+image_tag+"\n")
	fmt.Fprintf(os.Stdout, "image_id: "+image_id+"\n")

	return
}

func (skyc *Skycore) Get_etcd_value(path string) (value string, err error) {

	fmt.Fprintf(os.Stdout, "reading etcd "+path+"\n")
	response, err := skyc.Etcd_client.Get(path, false, false)
	if err != nil {
		return "", err
	}

	value = response.Node.Value

	return
}

func (skyc *Skycore) Get_etcd_shock2image(node_id string) (image_id string) {

	etcd_key := "/skycore/shock2image/" + node_id
	image_id, _ = skyc.Get_etcd_value(etcd_key)

	if image_id != "" {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("found imageid in etcd: %s=%s\n", etcd_key, image_id))
	} else {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("did not find imageid in etcd (%s) \n", etcd_key))
	}
	return
}

func (skyc *Skycore) Set_etcd_shock2image(node_id string, image_id string) {

	etcd_key := "/skycore/shock2image/" + node_id

	_, err := skyc.Etcd_client.Set(etcd_key, image_id, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("warning: could not write image_id to etcd, error: %s\n", err.Error()))
	} else {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("wrote imageid to etcd: %s=%s\n", etcd_key, image_id))
	}

	return
}

func (skyc *Skycore) Set_etcd_image(image_id string, etcd_repository string, etcd_tag string, etcd_shock_node string) {

	if image_id == "" {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("error: image_id empty, should not happen\n"))
		return
	}

	//c := etcd.NewClient(etcd_urls)
	c := skyc.Etcd_client
	etcd_key := "/skycore/image/" + image_id

	_, err := c.Set(etcd_key+"/repository", etcd_repository, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("waring: could not write repository to etcd, error: %s\n", err.Error()))
	}
	_, err = c.Set(etcd_key+"/tag", etcd_tag, 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("waring: could not write tag to etcd, error: %s\n", err.Error()))
	}
	if etcd_shock_node != "" {
		_, err = c.Set(etcd_key+"/shock_node", etcd_repository, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("waring: could not write shock_node to etcd, error: %s\n", err.Error()))
		}
	}
}

func (skyc *Skycore) Get_etcd_image(image_id string) (etcd_repository string, etcd_tag string, etcd_shock_node string) {

	if image_id == "" {
		return "", "", ""
	}

	etcd_key := "/skycore/image/" + image_id

	_, err := skyc.Get_etcd_value(etcd_key)

	if err != nil {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("did not find imageid in etcd (%s) error: %s \n", etcd_key, err.Error()))
		return "", "", ""
	}

	etcd_repository, _ = skyc.Get_etcd_value(etcd_key + "/repository")
	etcd_tag, _ = skyc.Get_etcd_value(etcd_key + "/tag")
	etcd_shock_node, _ = skyc.Get_etcd_value(etcd_key + "/shock_node")

	fmt.Fprintf(os.Stdout, fmt.Sprintf("found imageid in etcd: %s=%s\n", etcd_key, image_id))

	return
}

func (skyc *Skycore) skycore_load(command_arg string, request_tag string) (err error) {

	//skyc.Shock_client
	fmt.Fprint(os.Stdout, "skycore_load\n")

	command_arg_is_hex := false
	_, err = hex.DecodeString(command_arg)
	if err == nil {
		command_arg_is_hex = true
	}

	node_id := ""  // shock node id
	image_id := "" // docker image id, requires search of image in shock

	image_repository := ""
	image_tag := ""

	service_name := ""
	etcd_repository := ""
	etcd_tag := ""
	etcd_shock_node := ""

	if strings.HasPrefix(command_arg, "etcd:") {

		service_name = strings.TrimPrefix(command_arg, "etcd:")

		service_path := "/service_images/" + service_name
		shock_path := service_path + "/shock"

		image_url, err := skyc.Get_etcd_value(shock_path)
		if err != nil {
			return err
		}

		if !strings.Contains(image_url, "://") { // otherwise url.Parse will not parse correctly !
			image_url = "http://" + image_url
		}

		shock_node_url_obj, err := url.Parse(image_url)
		if err != nil {
			return err
		}

		if shock_node_url_obj.Host == "" {
			return errors.New("error: Host not found in url " + image_url)
		}

		fmt.Fprintf(os.Stdout, "image_url: "+image_url+"\n")
		fmt.Fprintf(os.Stdout, "shock_node_url_obj.Host: "+shock_node_url_obj.Host+"\n")
		fmt.Fprintf(os.Stdout, "shock_node_url_obj.Path: "+shock_node_url_obj.Path+"\n")

		if skyc.Shock_client.Host != "" {
			fmt.Fprintf(os.Stdout, fmt.Sprintf("warning: ignoring --shock argument, will use etcd value\n"))
		}

		shock_node_url_obj.RawQuery = "" // removes the shock node "?download" suffix

		node_id = strings.TrimPrefix(shock_node_url_obj.Path, "/node/")

		if node_id == shock_node_url_obj.Path {
			return errors.New("error: could not extract node_id from url, path=" + shock_node_url_obj.Path)
		}

		skyc.Shock_client.Host = shock_node_url_obj.Scheme + "://" + shock_node_url_obj.Host

		if uuid.Parse(node_id) == nil {
			return errors.New("error: node_id not a uuid \"" + node_id + "\"")
		}
		fmt.Fprintf(os.Stdout, "found node_id "+node_id+"\n")
		fmt.Fprintf(os.Stdout, "skyc.Shock_client.Host: "+skyc.Shock_client.Host+"\n")
	} else if uuid.Parse(command_arg) != nil { // Shock node
		node_id = command_arg

	} else if command_arg_is_hex {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("I'll guess that the argument %s is an image_id\n", command_arg))
		image_id = command_arg

		etcd_repository, etcd_tag, etcd_shock_node = skyc.Get_etcd_image(image_id)

	} else {
		return errors.New(fmt.Sprintf("error: argument is not a valid shock node uuid or docker image id\n"))
	}

	if node_id != "" { // we have a Shock node id

		// check if we have seen shock node before
		image_id = skyc.Get_etcd_shock2image(node_id)

		// get metadata for that image
		etcd_repository, etcd_tag, etcd_shock_node = skyc.Get_etcd_image(image_id)
	}

	image_attributes_known := false
	if image_id != "" && etcd_repository != "" && etcd_tag != "" {
		image_attributes_known = true

		image_repository = etcd_repository
		image_tag = etcd_tag
	}

	if node_id == "" {
		node_id = etcd_shock_node
	}

	if !image_attributes_known { // when only shock node id is know

		if node_id == "" {
			return errors.New(fmt.Sprintf("error: not enough information to find docker image\n"))
		}

		if skyc.Shock_client.Host == "" {
			return errors.New("error: shock host not defined")
		}

		image_repository, image_tag, image_id, err = skyc.get_dockerimage_shocknode_attributes(node_id)

		if err != nil {
			return

			//image_repository = ""
			//image_tag = ""
			//image_id = ""
		}

		if image_repository == "" || image_tag == "" || image_id == "" {
			return errors.New(fmt.Sprintf("error: attributes not read"))
		}

		if image_tag == "latest" {
			return errors.New(fmt.Sprintf("error: sorry, tag \"latest\" is not accepted by skycore. Please use some version number, date or timestamp.\n"))

		}

		skyc.Set_etcd_shock2image(node_id, image_id)
		skyc.Set_etcd_image(image_id, image_repository, image_tag, node_id)

	}

	if image_id == "" {
		return errors.New(fmt.Sprintf("error: no image_id found (should not happen)\n"))
	}

	fmt.Fprintf(os.Stdout, fmt.Sprintf("repository: %s\n", image_repository))
	fmt.Fprintf(os.Stdout, fmt.Sprintf("tag: %s\n", image_tag))
	fmt.Fprintf(os.Stdout, fmt.Sprintf("image_id: %s\n", image_id))

	if image_repository == "" || image_tag == "" || image_id == "" {
		return errors.New(fmt.Sprintf("error: attributes not available"))
	}
	// now check if this image_id is already available

	image_obj, err := skyc.Docker_client.InspectImage(image_id)

	if err != nil {
		// image not found
		image_obj = nil
	} else {

		fmt.Fprintf(os.Stdout, fmt.Sprintf("found image %s.\n", image_id))
		// do not exit here, we still might have to tag it
	}

	if image_obj == nil {
		// load image from shock and tag the image
		download_url := skyc.Shock_client.Host + "/node/" + node_id + "?download"

		fmt.Fprintf(os.Stdout, fmt.Sprintf("download_url: %s\n", download_url))
		err = dockerLoadImage(skyc.Docker_client, download_url, skyc.Shock_client.Token)
		if err != nil {
			return errors.New(fmt.Sprintf("Error loading docker image from Shock, err=%s", err.Error()))
		}
		time.Sleep(1 * time.Second)

		image_obj, err = skyc.Docker_client.InspectImage(image_id)
		if err != nil {
			return errors.New(fmt.Sprintf("Error: docker image %s has not been loaded, err=%s", image_id, err.Error()))
		} else {
			fmt.Fprintf(os.Stdout, "found image %s in docker engine\n", image_id)
		}
		tag_opts := docker.TagImageOptions{Repo: image_repository, Tag: image_tag}
		err = skyc.Docker_client.TagImage(image_id, tag_opts) // ignore error...
		if err != nil {
			//fmt.Fprintf("Error: docker image %s has not been loaded, err=%s", image_id, err.Error())
			fmt.Fprintf(os.Stdout, "error tagging image %s: %s\n", image_id, err.Error())
		} else {
			fmt.Fprintf(os.Stdout, "image %s tagged with %s:%s\n", image_id, image_repository, image_tag)
		}

	}

	if request_tag != "" {

		new_image_name := image_repository + ":" + request_tag
		fmt.Fprintf(os.Stdout, "trying to set additional tag %s", new_image_name)

		do_tag := false
		some_old_image_obj, err := skyc.Docker_client.InspectImage(new_image_name)
		if err != nil {
			// image not found, good, will tag
			do_tag = true
		} else {
			if some_old_image_obj.ID != image_id {
				fmt.Fprintf(os.Stdout, "trying to delete old image %s %s", some_old_image_obj.ID, new_image_name)
				err = skyc.Docker_client.RemoveImage(some_old_image_obj.ID)
				if err != nil {
					return errors.New(fmt.Sprintf("error deleting old image, maybe container is running?: %s", err.Error()))
				}

				do_tag = true
			}
		}

		if do_tag {
			tag_opts := docker.TagImageOptions{Repo: image_repository, Tag: request_tag}
			err = skyc.Docker_client.TagImage(image_id, tag_opts)
			if err != nil {

				return errors.New(fmt.Sprintf("error tagging image: %s", err.Error()))
			}
		}

	}

	return nil
}

func usage() {
	fmt.Fprintf(os.Stdout, fmt.Sprintf("\nUsage: %s command [options] [args]\n\n", os.Args[0]))

	fmt.Fprintf(os.Stdout,
		"\n"+
			"Commands:\n"+
			"\n"+
			"  pull [--shock=<url>] <arg>   \n"+
			"           Load image into docker engine, if not exists.\n"+
			"           Argument is shock node, image id or etcd service name\n"+
			"  push [--shock=<url>] <image_name>   \n"+
			"           Export image from Docker engine and upload to Shock.\n"+
			"           Argument is image name including tag: \"repository:tag\"\n",
	)

	fmt.Fprintf(os.Stdout, fmt.Sprintf("\nOptions:\n\n"))

	flags.PrintDefaults()
	fmt.Fprintf(os.Stdout, fmt.Sprintf("\nEnvironment variables that can be used: SKYCORE_SHOCK and SKYCORE_SHOCK_TOKEN\n\n"))

}

func main() {

	var shock_url string
	var shocktoken string
	var no_etcd bool
	var etcd_urls_string string
	var docker_socket string
	var private_image bool
	var request_tag string
	var help bool

	var etcd_urls []string

	flags = flag.NewFlagSet("name", flag.ContinueOnError)

	flags.StringVar(&shock_url, "shock", shock_url_default, "url of Shock server")
	flags.StringVar(&shocktoken, "token", shock_token_default, "OAuth token for Shock")
	flags.BoolVar(&no_etcd, "no_etcd", false, "Disable use of etcd")
	flags.StringVar(&etcd_urls_string, "etcd_urls", etcd_urls_string_default, "Comma separated list of etcd urls (default: "+etcd_urls_string_default+")")
	flags.StringVar(&docker_socket, "docker_socket", "", "docker socket")
	flags.StringVar(&request_tag, "tag", "", "tag image with tag, e.g. --tag=latest, this will tag image twice and untag existing image")
	flags.BoolVar(&private_image, "private", false, "Do not make image public")

	flags.BoolVar(&help, "help", false, "")

	if len(os.Args) < 2 {
		flags.Parse(os.Args)
		usage()
		os.Exit(1)
	}
	flags.Parse(os.Args[2:])

	if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "-help" {
		help = true
	}

	if help {
		usage()
		os.Exit(0)
	}

	if etcd_urls_string != "" {
		etcd_urls = strings.Split(etcd_urls_string, ",")
	}

	skyc := &Skycore{
		Shock_client:  shock.ShockClient{Token: shocktoken, Debug: true},
		Etcd_client:   nil,
		Docker_client: nil,
	}

	if !no_etcd {
		skyc.Etcd_client = etcd.NewClient(etcd_urls)
	}

	// docker client
	var err error

	skyc.Docker_client, err = docker.NewClient(docker_socket_default)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error creating docker client: ", err.Error())
		os.Exit(1)
	}

	version, err := skyc.Docker_client.Version()
	if err != nil {
		fmt.Fprintf(os.Stdout, "error using docker client: ", err.Error())
		os.Exit(1)
	}
	infoSlice := []string(*version)
	for _, line := range infoSlice {
		fmt.Fprintf(os.Stdout, fmt.Sprintf("%s\n", line))
	}

	if false {
		imgs, _ := skyc.Docker_client.ListImages(docker.ListImagesOptions{All: false})
		for _, img := range imgs {
			fmt.Println("ID: ", img.ID)
			fmt.Println("RepoTags: ", img.RepoTags)
			fmt.Println("Created: ", img.Created)
			fmt.Println("Size: ", img.Size)
			fmt.Println("VirtualSize: ", img.VirtualSize)
			fmt.Println("ParentId: ", img.ParentID)
		}
	}

	if shock_url != "" {

		shock_host_url, err := url.Parse(shock_url)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("shock_url %s cannot be parsed: %s \n", shock_url, err.Error()))
			os.Exit(1)
		}
		if shock_host_url.Scheme == "" {
			shock_host_url.Scheme = "http"
		}
		fmt.Fprintf(os.Stdout, fmt.Sprintf("shock host url: %s\n", shock_host_url.String()))

		skyc.Shock_client.Host = shock_host_url.String()

	}

	command := os.Args[1]

	command_arg := ""
	if command == "pull" || command == "push" {
		if len(flags.Args()) == 0 {

			fmt.Fprintf(os.Stderr, fmt.Sprintf("error: argument not found\n"))
			os.Exit(1)
		}

		if len(flags.Args()) > 1 {

			fmt.Fprintf(os.Stderr, fmt.Sprintf("error: more than one argument found. \n"))
			os.Exit(1)
		}

		command_arg = flags.Arg(0)
	}
	switch command {
	case "pull":
		err := skyc.skycore_load(command_arg, request_tag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}
		break
	case "push":

		node_id, err := skyc.save_image_to_shock(command_arg, private_image)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "New Shock node id: "+node_id+"\n")

		break
	case "help":
		usage()
		os.Exit(1)
	case "helloworld":
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%s is a useless command \n", flag.Arg(0)))
		break
	default:
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\"%s\" unknown command \n", flag.Arg(0)))
		usage()
		os.Exit(1)
	}

}
