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

func nnl(b []byte) []byte {
	b = bytes.Replace(b, []byte{13, 10}, []byte{10}, -1)
	b = bytes.Replace(b, []byte{13}, []byte{10}, -1)
	return b
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
		a:=strings.Split(string(nnl(b)),"\n")
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

func checkRepoPath( s string){
	fmt.Println(s,"looks legit")
}

func doCommand( wkPtr, lePtr, dsPtr *string, tail []string) {



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
			checkRepoPath(t[1])
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
                          if _, ok := REPOS[t[1]]; ! ok {
                                fmt.Println(t[1],"not in workspace.")
                          }else{
	                        fmt.Println("adding metarepo "+tail[1])
			  }
                  }
                }else{
                        fmt.Println("need metarepo:repo")
                }
          }else if tail[0]=="changerepo"{
                t:=strings.Split(tail[1],":")
                if len(t)>1{
                  if _, ok := REPOS[t[0]]; ! ok {
                        fmt.Println(t[0],"not in workspace.")
                  }else{
                        checkRepoPath(t[1])
                        fmt.Println("updated repo "+tail[1])
                  }
                }else{
                        fmt.Println("need repo:path")
                }
          }else if tail[0]=="changemeta"{
                t:=strings.Split(tail[1],":")
                if len(t)>1{
                  if _, ok := METAS[t[0]]; ! ok {
                        fmt.Println(t[0],"not in workspace.")
                  }else{
	                  if _, ok := REPOS[t[1]]; ! ok {
	                        fmt.Println(t[1],"not in workspace.")
	                  }else{
	                        fmt.Println("updated metarepo "+tail[1])
			  }
                  }
                }else{
                        fmt.Println("need metarepo:repo")
                }
          }else if tail[0]=="delrepo"{
                t:=tail[1]
                
                  if _, ok := REPOS[t]; ! ok {
                        fmt.Println(t,"not in workspace.")
                  }else{
                        fmt.Println("removing repo "+tail[1])
                  }
                
          }else if tail[0]=="delmeta"{
                t:=tail[1]
                
                  if _, ok := METAS[t]; ! ok {
                        fmt.Println(t,"not in workspace.")
                  }else{
                        fmt.Println("removing metarepo "+tail[1])
                  }
               
          }else if tail[0]=="checkrepo"{
                t:=tail[1]

                  if r, ok := REPOS[t]; ! ok {
                  	if r, ok := METAS[t]; ! ok {
                  	      fmt.Println(t,"not in workspace.")
                  	}else{
			      v, _:= REPOS[r]
                  	      fmt.Println("checking repo "+v)
                  	}
                  }else{
                        fmt.Println("checking repo "+r)
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
                 
              }else if tail[0]=="rehash"{
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Rehashing", p,":")
                        fmt.Println(string(contents))

              }else if tail[0]=="addto"{
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Adding to", p,":")
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

func main() {

        lePtr := flag.String("e", "noch", "Local Module Line Ending Style (cr|crlf|nl|noch)")
        dsPtr := flag.String("s", "combined", "Local Package Directory Style (combined|flat|paths)")
        wkPtr := flag.String("d", "./", "workspace location")

        flag.Usage = func() {
            fmt.Fprintf(os.Stderr, "\nUsage of %s: package %s\n\n", os.Args[0]," <flags> commmand ")
            fmt.Println("  Flags:\n")
            flag.PrintDefaults()
            fmt.Fprintf(os.Stderr, `
  Commands:

    init                         Initialize a workspace
    repolist                     List repos configured for workspace
    metalist                     List metarepos configured for workspace
    addrepo    <repo:path>       Add a repo (and path to that repo) to the workspace
    addmeta    <metarepo:repo>   Add a metarepo (and specific repo for the metarepo) to the workspace
    changerepo <repo:path>       Change the path for an existing repo for the workspace
    changemeta <metarepo:repo>   Change the specific repo for an existing metarepo for the workspace
    delrepo    <repo>            Remove a repo from the workspace
    delmeta    <metarepo>        Remove a metarepo from the workspace
    checkrepo  <repo>            Check the status of a repo for the workspace
    enroll     <package>         Enroll (create) a package in the workspace
    status     <package|all>     Check the status of a package or packages in the workspace and the repos
    latest     <package|all>     Retrieve the latest version of a package from the repos to the workspace
    rehash     <package|all>     Update the hashes of local files in the workspace for the package or packages
    addto      <package>         Add a local file to the package in the workspace
    updates    <package|all>     Report updates (from repos) to a package or packages in the workspace
    exact      <package|all>     Retrieve a specific version of a package from the repos to the workspace
    provider   <package|all>     Report which repo provided (if any) the package or packages in the workspace
`)


        }

        flag.Parse()


        tail:= flag.Args()
        doCommand( wkPtr, lePtr, dsPtr, tail)

}


