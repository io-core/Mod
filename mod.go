// MIT License
//
// Copyright (c) 2019 the io-core authors
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
	"time"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

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

func putWorkspaceSettings(wk string, WSV, REPOS, METAS map[string]string){
	t := time.Now()
	err := os.Rename(path.Clean(wk)+"/Packaging.csv",path.Clean(wk)+"/Packaging.csv."+t.Format("20060102150405"))
	check(err)
	f, err := os.Create(path.Clean(wk)+"/Packaging.csv"); check(err)
	defer f.Close()
        _, err = f.WriteString("setting,value\n"); check(err)
	for k,v := range WSV {
		//  println(k+","+v)
		_, err := f.WriteString(k+","+v+"\n"); check(err)
	}
        for k,v := range REPOS {
        	//  println("repo,"+k+":"+v)
        	_, err := f.WriteString("repo,"+k+":"+v+"\n"); check(err)
        }
        for k,v := range METAS {
        	//  println("meta-repo,"+k+":"+v)
		_, err := f.WriteString("meta-repo,"+k+":"+v+"\n"); check(err)
        }

	f.Sync()
	err = os.Remove(path.Clean(wk)+"/Packaging.csv."+t.Format("20060102150405")); check(err)

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

func repoPathOK( s string) bool {
	r:=false
	t:=strings.Split(s,"/")
	if len(t)>1 {
	  fmt.Println(s,"looks legit")
	  r=true
	}else{
          fmt.Println(s,"expecting a path with at least one slash in it")
	}
	return r
}

func listPackages( wk string){
     //           _,_,REPOS := getWorkspaceSettings(wk)
     //           for r,v:=range REPOS { fmt.Println(r,v)}
}

func repoList( wk string){
                _,_,REPOS := getWorkspaceSettings(wk)
                for r,v:=range REPOS { fmt.Println(r,v)}
}

func metaList( wk string){
                _,METAS,_ := getWorkspaceSettings(wk)
                for m,v:=range METAS { fmt.Println(m,v)}
}

func addRepo( wk string, WSV, REPOS, METAS map[string]string, tail []string){
                t:=strings.Split(tail[1],":")
                if len(t)>1{
                  if _, ok := REPOS[t[0]]; ok {
                        fmt.Println(t[0],"already in workspace.")
                  }else{
                        if repoPathOK(t[1]){
                          fmt.Println("adding repo "+tail[1])
			  REPOS[t[0]]=t[1]
			  putWorkspaceSettings(wk,WSV,REPOS,METAS)
			}
                  }
                }else{
                        fmt.Println("need repo:path")
                }
}

func addMeta( wk string, WSV, REPOS, METAS map[string]string, tail []string){
                t:=strings.Split(tail[1],":")
                if len(t)>1{
                  if _, ok := METAS[t[0]]; ok {
                        fmt.Println(t[0],"already in workspace.")
                  }else{
                          if _, ok := REPOS[t[1]]; ! ok {
                                fmt.Println(t[1],"not in workspace.")
                          }else{
                                fmt.Println("adding metarepo "+tail[1])
                          	METAS[t[0]]=t[1]
                          	putWorkspaceSettings(wk,WSV,REPOS,METAS)
                          }     
                  }
                }else{
                        fmt.Println("need metarepo:repo")
                }       
}

func changeRepo( wk string, WSV, REPOS, METAS map[string]string, tail []string){
                t:=strings.Split(tail[1],":")
                if len(t)>1{
                  if _, ok := REPOS[t[0]]; ! ok {
                        fmt.Println(t[0],"not in workspace.")
                  }else{
                        if repoPathOK(t[1]){
                          fmt.Println("updating repo "+tail[1])
                          REPOS[t[0]]=t[1]
                          putWorkspaceSettings(wk,WSV,REPOS,METAS)
			}
                  }
                }else{
                        fmt.Println("need repo:path")
                }
}

func changeMeta( wk string, WSV, REPOS, METAS map[string]string, tail []string){
                t:=strings.Split(tail[1],":")
                if len(t)>1{
                  if _, ok := METAS[t[0]]; ! ok {
                        fmt.Println(t[0],"not in workspace.")
                  }else{
                          if _, ok := REPOS[t[1]]; ! ok {
                                fmt.Println(t[1],"not in workspace.")
                          }else{
                                fmt.Println("updating metarepo "+tail[1])
                                METAS[t[0]]=t[1]
                                putWorkspaceSettings(wk,WSV,REPOS,METAS)
                          }
                  }
                }else{
                        fmt.Println("need metarepo:repo")
                }
}

func delRepo( wk string, WSV, REPOS, METAS map[string]string, tail []string){
                t:=tail[1]                 
                  if _, ok := REPOS[t]; ! ok {
                        fmt.Println(t,"not in workspace.")
                  }else{
			metahas:=""
			for m,x:=range METAS{
				if x == t { metahas = m }
			}
			if metahas == ""{
                        	fmt.Println("removing repo "+tail[1])
                                delete(REPOS,t)
                                putWorkspaceSettings(wk,WSV,REPOS,METAS)

			}else{
				fmt.Println("metarepo "+metahas+" refers to "+t+". Please adjust the metarepo first.")
			}
                  }     
}

func delMeta( wk string, WSV, REPOS, METAS map[string]string, tail []string){
                t:=tail[1]
                  if _, ok := METAS[t]; ! ok {
                        fmt.Println(t,"not in workspace.")
                  }else{
                        fmt.Println("removing metarepo "+tail[1])
			delete(METAS,t)
			putWorkspaceSettings(wk,WSV,REPOS,METAS)
                  }
}

func checkRepo( REPOS, METAS map[string]string, tail []string){
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
}

func enrollPackage( wkPtr *string, tail []string){
                sPkgs := buildSourceList(*wkPtr,[]string{"all"})
                nPkgs := strings.Split(tail[1],",")
                if len(nPkgs)!=1{
                        fmt.Println("Only enroll one package at a time")
                }else{
                        fmt.Println("Enrolling",nPkgs[0],sPkgs)
                }
}

func packageStatus(p string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}

func latestPackage(p string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}

func rehashPackage(p string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Rehashing", p,":")
                        fmt.Println(string(contents))
}

func addtoPackage(p string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Adding to", p,":")
                        fmt.Println(string(contents))
}

func packageUpdates(p string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}

func exactPackage(p string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}

func packageProvider(p string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}



func doCommand( wkPtr, lePtr, dsPtr *string, tail []string) {

        if len(tail)==1 {
                if tail[0] == "init"     { initWorkspace(*wkPtr,*lePtr,*dsPtr)
          }else if tail[0] == "list"     { listPackages(*wkPtr)
          }else if tail[0] == "repolist" { repoList(*wkPtr)
          }else if tail[0] == "metalist" { metaList(*wkPtr)
          }else{
                fmt.Println("Incomplete command. exiting.")
	  }
        }else if len(tail)>1 {
          WSV,METAS,REPOS := getWorkspaceSettings(*wkPtr)
                if tail[0] == "addrepo"   { addRepo(*wkPtr,WSV,REPOS,METAS,tail)
          }else if tail[0] == "addmeta"   { addMeta(*wkPtr,WSV,REPOS,METAS,tail)
          }else if tail[0] == "changerepo"{ changeRepo(*wkPtr,WSV,REPOS,METAS,tail)
          }else if tail[0] == "changemeta"{ changeMeta(*wkPtr,WSV,REPOS,METAS,tail)
          }else if tail[0] == "delrepo"   { delRepo(*wkPtr,WSV,REPOS,METAS,tail)
          }else if tail[0] == "delmeta"   { delMeta(*wkPtr,WSV,REPOS,METAS,tail)
          }else if tail[0] == "checkrepo" { checkRepo(REPOS,METAS,tail)
          }else if tail[0] == "enroll"    { enrollPackage(wkPtr,tail)
	  }else{
	    sPkgs := buildSourceList(*wkPtr,strings.Split(tail[1],","))
	    for _, p := range sPkgs {
                    if tail[0]=="status"  { packageStatus(p)
              }else if tail[0]=="latest"  { latestPackage(p)
              }else if tail[0]=="rehash"  { rehashPackage(p)
              }else if tail[0]=="addto"   { addtoPackage(p)
              }else if tail[0]=="updates" { packageUpdates(p)
              }else if tail[0]=="exact"   { exactPackage(p)
              }else if tail[0]=="provider"{ packageProvider(p)
              }else{
                fmt.Println(tail[0]," means what?")
              }
	    }
	  }
	}else{
	  flag.Usage()
	}

}

func main() {

        lePtr := flag.String("e", "noch", "Local Module Line Ending Style (cr|crlf|nl|noch)")
        dsPtr := flag.String("s", "combined", "Local Package Directory Style (combined|flat|paths)")
        wkPtr := flag.String("d", "./", "workspace location")

        flag.Usage = func() {
            fmt.Fprintf(os.Stderr, "\nUsage of %s: mod %s\n\n", os.Args[0]," <flags> commmand ")
            fmt.Println("  Flags:\n")
            flag.PrintDefaults()
            fmt.Fprintf(os.Stderr, `
  Commands:

    init                         Initialize a workspace
    list                         List the packages in the workspace
    repolist                     List repos configured for workspace
    metalist                     List metarepos configured for workspace
    addrepo    <repo:path>       Add a repo (and path to that repo) to the workspace
    addmeta    <metarepo:repo>   Add a metarepo (and specific repo for the metarepo) to the workspace
    changerepo <repo:path>       Change the path for an existing repo for the workspace
    changemeta <metarepo:repo>   Change the specific repo for an existing metarepo for the workspace
    delrepo    <repo>            Remove a repo from the workspace
    delmeta    <metarepo>        Remove a metarepo from the workspace
    checkrepo  <repo>            Check the status of a repo for the workspace
    enroll     <package file(s)> Enroll (create) a package in the workspace with file(s)
    status     <package|all>     Check the status of a package or packages in the workspace and the repos
    latest     <package|all>     Retrieve the latest version of a package from the repos to the workspace
    rehash     <package|all>     Update the hashes of local files in the workspace for the package or packages
    addto      <package file(s)> Add local file(s) to the package in the workspace
    updates    <package|all>     Report updates (from repos) to a package or packages in the workspace
    exact      <package|all>     Retrieve a specific version of a package from the repos to the workspace
    provider   <package|all>     Report which repo provided (if any) the package or packages in the workspace
`)


        }

        flag.Parse()


        tail:= flag.Args()
        doCommand( wkPtr, lePtr, dsPtr, tail)

}


