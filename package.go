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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"bytes"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func NormalizeNewlines(d []byte) []byte {
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}

func getWorkspaceSettings(wk string) (map[string]string,map[string]string,map[string]string){
	WSV := make(map[string]string)
        REPOS := make(map[string]string)
        METAS := make(map[string]string)
//        fmt.Println("Loading workspace settings", wk)
	if _, err := os.Stat(path.Clean(wk)+"/Packaging.csv"); err == nil {
		b, err := ioutil.ReadFile(path.Clean(wk)+"/Packaging.csv")
		if err != nil {
			fmt.Print("Couldn't read Packaging.csv")
			os.Exit(1)
		}
		a:=strings.Split(string(NormalizeNewlines(b)),"\n")
                for _,b:= range a[1:]{
			c:=strings.Split(b,",")
			if len(c)>1{
			  if c[0]=="meta-repo" {
			    t:=strings.Split(c[1],":")
			    if len(t)>1{
                              METAS[t[0]]=t[1]
			    }
			  }else if c[0]=="repo"{
                            t:=strings.Split(c[1],":")
                            if len(t)>1{
                              REPOS[t[0]]=t[1]
                            }
			  }else{
			    WSV[c[0]]=c[1]
			  }
			}
		}
	}else{
                fmt.Println("Workspace",wk,"is not initialized, exiting.")
                os.Exit(1)	
	}
	if _, ok := WSV["workspace-module-line-ending"]; ! ok {
	        fmt.Println("workspace-module-line-ending missing in workspace settings. exiting.")
		os.Exit(1)
	}
        if _, ok := WSV["workspace-packages-dirstyle"]; ! ok {
                fmt.Println("workspace-packages-dirstyle missing in workspace settings. exiting.")
                os.Exit(1)
        }
	return WSV, METAS, REPOS
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

	flag.Parse()

	
	tail:= flag.Args()
	var contents []byte

        if len(tail)==1 {
          if tail[0]=="init"{
                initWorkspace(*wkPtr,*lePtr,*dsPtr)
          }else if tail[0]=="repolist"{
          	_,_,REPOS := getWorkspaceSettings(*wkPtr)
                for r,v:=range REPOS { fmt.Println(r,v)}
          }else if tail[0]=="metalist"{
                _,METAS,_ := getWorkspaceSettings(*wkPtr)
                for m,v:=range METAS { fmt.Println(m,v)}
          }else{
                fmt.Println("Incomplete command. exiting.")
	  }
        }else if len(tail)>1 {

          WSV,METAS,REPOS := getWorkspaceSettings(*wkPtr)
	  if 1==2 { fmt.Println(WSV,METAS,REPOS) }

          if tail[0]=="addrepo"{
		t:=strings.Split(tail[1],":")
		if len(t)>1{
		  if _, ok := REPOS[t[0]]; ok {
	                fmt.Println(t[0],"already in workspace.")
        	  }else{
			fmt.Println("adding repo "+tail[1])
		  }
		}else{
			fmt.Println("need repo:path")
		}
          }else if tail[0]=="addmeta"{
                t:=strings.Split(tail[1],":")
                if len(t)>1{
                  if _, ok := METAS[t[0]]; ok {
                        fmt.Println(t[0],"already in workspace.")
                  }else{
                        fmt.Println("adding metarepo "+tail[1])
                  }
                }else{
                        fmt.Println("need metarepo:repo")
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

}
