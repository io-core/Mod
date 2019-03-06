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
	"path"
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

func getWorkspaceSettings(wk string) map[string]string{
	var WSV map[string]string
        fmt.Println("Loading workspace settings", wk)
	if _, err := os.Stat(path.Clean(wk)+"/Packaging.csv"); err == nil {
		b, err := ioutil.ReadFile(path.Clean(wk)+"/Packaging.csv")
		if err != nil {
			fmt.Print("Couldn't read Packaging.csv")
			os.Exit(1)
		}
		fmt.Println(string(b))
	}else{
                fmt.Println("Workspace",wk,"is not initialized, exiting.")
                os.Exit(1)	
	}
	return WSV
}

func buildSourceList(wk string, s []string) []string {
	  var files []string

          if s[0]=="all" {
                fileInfo, _ := ioutil.ReadDir(path.Clean(wk))
                for _, file := range fileInfo {
                  n:=file.Name()
                  if len(n)>4 {
                    if n[len(n)-4:]==".Pkg" {
                      files = append(files, n[0:len(n)-4])
                    }
                  }
                }
          }else{
                for _, fn := range s {
                        if _, err := os.Stat(path.Clean(wk)+"/"+fn+".Pkg"); err == nil {
                                files = append(files, fn)
                        }else{
                                fmt.Println("Package",fn,"Not Found, exiting.")
                                os.Exit(1)
                        }
                }
          }
    	  return files
}

func initWorkspace(wk, le, ds string){
        fmt.Println("Initializing the workspace", wk)

	e:="\n"
        if le=="cr" {
                e="\r"
                fmt.Println("CR local package line ending style")
        }else if le=="crlf" {
		e="\r\n"
                fmt.Println("CRLF local package line ending style")
        }else if le=="nl"{
                fmt.Println("NL local package line ending style")
        }else{
                fmt.Println("Line ending style",le,"not understood.")
                os.Exit(1)
        }

	if ds=="combined" {
                fmt.Println("Combined Local Package Directory Style")
	}else if ds=="flat"{
                fmt.Println("Flat Local Package Directory Style")
	}else if ds=="paths"{
                fmt.Println("Paths Local Package Directory Style")
	}else{
		fmt.Println("Local Package Directory Style",ds,"not understood.")
		os.Exit(1)
	}

        if _, err := os.Stat(path.Clean(wk)+"/Packaging.csv"); err != nil {
	        c := []byte("setting,value"+e+"workspace-module-line-ending,"+le+e+"workspace-packages-dirstyle,"+ds+e)
	        err := ioutil.WriteFile(path.Clean(wk)+"/Packaging.csv", c, 0644)
	        if err != nil{
                        fmt.Println("Error Creating Packaging.csv file.")
			os.Exit(1)
		}else{
                	fmt.Println("Created Packaging.csv file.")
		}
        }else{
                fmt.Println("Packaging.csv already exists in",wk,"exiting.")
                os.Exit(1)
        }

}

func main() {
	
        lePtr := flag.String("e", "cr", "Local Module Line Ending Style (cr|crlf|nl)")
        dsPtr := flag.String("s", "combined", "Local Package Directory Style (combined|flat|paths)")
        wkPtr := flag.String("d", "./", "workspace location")


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

        if len(tail)==1 {
          if tail[0]=="init"{
                initWorkspace(*wkPtr,*lePtr,*dsPtr)
          }
        }else if len(tail)>1 {

          WSV := getWorkspaceSettings(*wkPtr)
	  fmt.Println(WSV)

          if tail[0]=="addrepo"{
                sPkgs := buildSourceList(*wkPtr,[]string{"all"})
                nPkgs := strings.Split(tail[1],",")
                if len(nPkgs)!=1{
                        fmt.Println("Only enroll one package at a time")
                }else{
                        fmt.Println("Enrolling",nPkgs[0],sPkgs)
                }
          }else if tail[0]=="enroll"{
          	sPkgs := buildSourceList(*wkPtr,[]string{"all"})
          	nPkgs := strings.Split(tail[1],",")
          	if len(nPkgs)!=1{
                	fmt.Println("Only enroll one package at a time")
          	}else{
          		fmt.Println("Enrolling",nPkgs[0],sPkgs)
          	}
	  }else{

	    sPkgs := buildSourceList(*wkPtr,strings.Split(tail[1],","))

	    fmt.Println(sPkgs)

	    for _, p := range sPkgs {

              if tail[0]=="status"{
			contents, _ = ioutil.ReadFile(p+".Pkg")
                	fmt.Println("Status of", p,":")
			fmt.Println(string(contents))
		
              }else if tail[0]=="latest"{
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
                 
              }else if tail[0]=="updates"{
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
                
              }else if tail[0]=="exact"{
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
                
              }else if tail[0]=="provider"{
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))

              }else{
                fmt.Println(tail[0]," means what?")
              }
	    }
	  }
	}else{
	  fmt.Println("Usage: get <command> <package> [options...]\n try: status latest dependencies")
	}

	if(1==2){
		fmt.Println(*lePtr,*dsPtr,contents,*inFilePtr,*aMessagePtr,*formatPtr,*pkeyPtr,*bkeyPtr,*tkeysPtr,*checkPtr,*rkeyPtr)
	}
	
}
