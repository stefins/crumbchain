package crumb

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/francoispqt/gojay"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"io"
	"github.com/cheggaaa/pb/v3"
)

type Crumb struct {
	index   int
	name    string
	hash    string
	prev_hash string
	content string
}

func (c *Crumb) MarshalJSONObject(enc *gojay.Encoder) {
	enc.IntKey("index", c.index)
	enc.StringKey("name", c.name)
	enc.StringKey("hash", c.hash)
	enc.StringKey("prev_hash",c.prev_hash)
	enc.StringKey("content", c.content)
}

func (c *Crumb) IsNil() bool {
	return c == nil
}

func (c *Crumb) info() string {
	return c.name + c.hash + c.content + "[" + strconv.Itoa(c.index) + "]"
}

func Crumber(name string,size int) {
	//To chuck the file into 1 MB
	chuck := 1024 * 1024 * size
	var content strings.Builder
	cmb := []*Crumb{}
	f, err := os.Open("./" + name)
	if err != nil {
		fmt.Println("File Not Found!")
		return
	}
	reader := bufio.NewReader(f)
	buf := make([]byte,256)
  for {
        n, err := reader.Read(buf)
        if err != nil {
           if err != io.EOF {
               log.Fatal(err)
           }
           break
        }
        //fmt.Print(string(buf[0:n]))
        content.WriteString(string(buf[0:n]))
  }
	fmt.Println("Encoding to base64....")
	encoded := base64.StdEncoding.EncodeToString([]byte(content.String()))

	count := len(encoded)/chuck +1
	fmt.Println("Reading The File")
	bar := pb.StartNew(count)
	for i := 0; ; i++ {

		if i < len(encoded)/chuck {
			h := sha256.New()
			h.Write([]byte(encoded[i*chuck : (i+1)*chuck]))
			loc := new(Crumb)
			if i==0 {
				loc.prev_hash = "0000000"
			}else{
				loc.prev_hash = cmb[i-1].hash
			}
			h.Write([]byte(loc.prev_hash))
			loc.index = i
			loc.name = name
			loc.hash = string(h.Sum(nil))
			loc.content = encoded[i*chuck : (i+1)*chuck]
			cmb = append(cmb, loc)
			bar.Increment()
		} else {
			h := sha256.New()
			h.Write([]byte(encoded[i*chuck:]))
			loc := new(Crumb)
			if i==0 {
				loc.prev_hash = "0000000"
			}else{
				loc.prev_hash = cmb[i-1].hash
			}
			h.Write([]byte(loc.prev_hash))
			loc.index = i
			loc.name = name
			loc.hash = string(h.Sum(nil))
			loc.content = encoded[i*chuck:]
			cmb = append(cmb, loc)
			bar.Increment()
			break
		}
	}
	bar.Finish()
	fmt.Println("Writing Crumbs ....")
	os.Mkdir(name+"-crumb", 0755)
	count = len(cmb)
	bar = pb.StartNew(count)
	for i, cont := range cmb {
		data := &Crumb{cont.index, cont.name, cont.hash,cont.prev_hash,cont.content}
		b, err := gojay.MarshalJSONObject(data)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(name+"-crumb/"+name+strconv.Itoa(i)+".cbc", b, 0644)
		bar.Increment()
	}
	bar.Finish()
}
