package crumb

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"encoding/base64"
	_ "strings"
	"io"
	"github.com/golang/protobuf/proto"
)

func Crumber(name string,size int) {
	//To chuck the file into 1 MB
	chuck := 1024 * 1024 * size
	//var content strings.Builder
	f, err := os.Open("./" + name)
	if err != nil {
		fmt.Println("File Not Found!")
		return
	}
	os.Mkdir(name+"-crumb", 0755)
	reader := bufio.NewReader(f)
	buf := make([]byte,chuck)
	var fuck int
	var prev_crum Crumb
  for i := 0; ; i++ {
        n, err := reader.Read(buf)
        if err != nil {
           if err != io.EOF {
               log.Fatal(err)
           }
           break
        }
				fuck++
				fmt.Println(fuck)
        //fmt.Print(string(buf[0:n]))
        //content.WriteString(string(buf[0:n]))
				loc := new(Crumb)
				if i==0 {
					loc.PrevHash = []byte("0000000")
				}else{
					loc.PrevHash = prev_crum.Hash
				}
				h := sha256.New()
				h.Write([]byte(buf[0:n]))
				h.Write([]byte(loc.PrevHash))
				loc.Index = int64(i)
				loc.Name = name
				loc.Hash = h.Sum(nil)
				loc.Content = base64.StdEncoding.EncodeToString([]byte(buf[0:n]))
				data := &Crumb{Index:loc.Index, Name:loc.Name, Hash:loc.Hash, PrevHash:loc.PrevHash, Content:loc.Content}
				b,err := proto.Marshal(data)
				if err != nil {
					log.Fatal(err)
				}
				err = ioutil.WriteFile(name+"-crumb/"+name+strconv.Itoa(i)+".cbc", b, 0644)
				prev_crum = *loc
  }
/*
	fmt.Println("Encoding to base64....")
	encoded := base64.StdEncoding.EncodeToString([]byte(content.String()))

	count := len(encoded)/chuck +1
	fmt.Println("Reading The File")
	bar := pb.StartNew(count)
	for i := 0; ; i++ {

		if i < len(encoded)/chuck {
			loc := new(Crumb)
			if i==0 {
				loc.PrevHash = []byte("0000000")
			}else{
				loc.PrevHash = cmb[i-1].Hash
			}
			h := sha256.New()
			h.Write([]byte(encoded[i*chuck : (i+1)*chuck]))
			h.Write([]byte(loc.PrevHash))
			loc.Index = int64(i)
			loc.Name = name
			loc.Hash = h.Sum(nil)
			loc.Content = encoded[i*chuck : (i+1)*chuck]
			cmb = append(cmb, loc)
			bar.Increment()
		} else {
			loc := new(Crumb)
			if i==0 {
				loc.PrevHash = []byte("0000000")
			}else{
				loc.PrevHash = cmb[i-1].Hash
			}
			h := sha256.New()
			h.Write([]byte(encoded[i*chuck:]))
			h.Write([]byte(loc.PrevHash))
			loc.Index = int64(i)
			loc.Name = name
			loc.Hash = h.Sum(nil)
			loc.Content = encoded[i*chuck:]
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
		data := &Crumb{Index:cont.Index, Name:cont.Name, Hash:cont.Hash, PrevHash:cont.PrevHash, Content:cont.Content}
		b,err := proto.Marshal(data)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile(name+"-crumb/"+name+strconv.Itoa(i)+".cbc", b, 0644)
		bar.Increment()
	}
	bar.Finish()
*/
}
