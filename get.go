// MIT License
//
// Copyright (c) 2018 the io-core authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
//	"crypto"
//	"crypto/rand"
//	"crypto/rsa"
//	"crypto/sha256"
//	"crypto/x509"
//	"encoding/base64"
//	"encoding/pem"
	"flag"
	"fmt"
//	"github.com/io-core/attest/s2r"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
//	"time"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func buildSourceList(s []string) []string {
	  var files []string

          if s[0]=="all" {
                fileInfo, _ := ioutil.ReadDir(".")
                for _, file := range fileInfo {
                  n:=file.Name()
                  if len(n)>4 {
                    if n[len(n)-4:]==".Pkg" {
                      files = append(files, n)
                    }
                  }
                }
          }else{
                for _, fn := range s {
                        if _, err := os.Stat("./"+fn+".Pkg"); err == nil {
                                files = append(files, fn+".Pkg")
                        }else{
                                fmt.Println("Package",fn,"Not Found, exiting.")
                                os.Exit(1)
                        }
                }
          }
    	  return files
}

func main() {
	
	inFilePtr := flag.String("i", "-", "input file")
	aMessagePtr := flag.String("a", "signed", "attest message")
	formatPtr := flag.String("f", "oberon", "attest comment style")
	pkeyPtr := flag.String("p", os.Getenv("HOME")+"/.ssh/id_rsa", "path to rsa private key file")
	bkeyPtr := flag.String("b", os.Getenv("HOME")+"/.ssh/id_rsa.pub", "path to rsa public key file")
        tkeysPtr := flag.String("t", os.Getenv("HOME")+"/.ssh/trusted_devs", "path to trusted_devs file")
	checkPtr := flag.Bool("c", false, "check instead of sign")
        rkeyPtr := flag.Bool("k", false, "retrieve public key from input file")

	flag.Parse()

	

	iam := filepath.Base(os.Args[0])
	if iam == "acheck" {
		f := true
		checkPtr = &f
	}

	tail:= flag.Args()
	var contents []byte

        if len(tail)>1 {

	  sPkgs := buildSourceList(strings.Split(tail[1],","))

	  fmt.Println(sPkgs)

          if tail[0]=="status"{
		if tail[1]=="all"{
		}else{
			contents, _ = ioutil.ReadFile(tail[1]+".Pkg")
                	fmt.Println("Status of", tail[1],":")
			fmt.Println(string(contents))
		}
          }else if tail[0]=="latest"{
                if tail[1]=="all"{
                }else{
                        contents, _ = ioutil.ReadFile(tail[1]+".Pkg")
                        fmt.Println("Status of", tail[1],":")
                        fmt.Println(string(contents))
                } 
          }else if tail[0]=="updates"{
                if tail[1]=="all"{
                }else{
                        contents, _ = ioutil.ReadFile(tail[1]+".Pkg")
                        fmt.Println("Status of", tail[1],":")
                        fmt.Println(string(contents))
                }
          }else if tail[0]=="exact"{
                if tail[1]=="all"{
                }else{
                        contents, _ = ioutil.ReadFile(tail[1]+".Pkg")
                        fmt.Println("Status of", tail[1],":")
                        fmt.Println(string(contents))
                }
          }else{
                fmt.Println("what?")
          }
	}else{
	  fmt.Println("Usage: get <command> <package> [options...]\n try: status latest dependencies")
	}

	if(1==2){
		fmt.Println(contents,*inFilePtr,*aMessagePtr,*formatPtr,*pkeyPtr,*bkeyPtr,*tkeysPtr,*checkPtr,*rkeyPtr)
	}
	
}
