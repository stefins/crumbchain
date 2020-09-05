package main

import (
  "flag"

  "github.com/iamstefin/crumbchain/crumb"
  "github.com/iamstefin/crumbchain/crumbjoiner"
)

func main()  {
  filename := flag.String("c", "", "Do the crumbing")
  join := flag.String("j","","Join the crumbs")
  intr := flag.String("i","","Check the integrity")
  size := flag.Int("s",10,"Set the size of the crumbs")
  flag.Parse()
  if (*filename == "") {
  }else{
    crumb.Crumber(*filename,*size)
    return
  }
  if (*join == ""){
  }else{
    crumbjoiner.Joiner(*join)
    return
  }
  if (*intr == "") {
  }else{
    crumbjoiner.Integrity(*intr)
    return
  }
  flag.Usage()
}
