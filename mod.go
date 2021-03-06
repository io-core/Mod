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
	"net/http"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"bytes"
	"time"
	"log"
	"path/filepath"
	"io"
	"crypto/sha256"
	git "gopkg.in/src-d/go-git.v4"
        //"gopkg.in/src-d/go-git.v4/storage/memory"

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

func getPackageSettings(pn,pl string) (string,string,string,string,string,string,map[string]string,map[string]string){
        IMP := make(map[string]string)
        PRO := make(map[string]string)
	name:=""
	version:="v0.0.0"
	from:="local"
	retrieved:="?"
	license:="unknown"
	authors:="anonymous"
	var a []string
        if _, err := os.Stat(path.Clean(pl)+"/"+pn+".Pkg"); err == nil {
                b, err := ioutil.ReadFile(path.Clean(pl)+"/"+pn+".Pkg")
                if err != nil {
                        fmt.Print("Couldn't read Package", pn)
                        os.Exit(1)
                }
                a=strings.Split(string(nnl(b)),"\n")
	}else{
		fmt.Println("Couldn't stat ",path.Clean(pl)+"/"+pn+".Pkg")
	}
	c:=0
	for ;c<len(a);c++ {
	    l:=strings.Split(strings.TrimSpace(a[c]),",")	
	    if len(l)>2 {
                if l[0]=="package" {
                        name = l[1]
			version = l[2]
                }else if l[0]=="from" { 
                        from = l[1]  
			retrieved = l[2]
                }else if l[0]=="license" { 
                        license = l[1]  
			authors = l[2]
		}else if l[0]=="r" { 
			IMP[l[1]]=l[2]
                }else if l[0]=="p" { 
			PRO[l[1]]=l[2] 
		}
	    }
	}
	if 1==2 {fmt.Println(authors,retrieved)}
	return name,version,from,retrieved,license,authors,IMP,PRO
}

func putPackageSettings(pn,pl,name,version,from,retrieved,license,authors string, IMP,PRO map[string]string){
        t := time.Now()
        err := os.Rename(path.Clean(pl)+"/"+pn+".Pkg",path.Clean(pl)+"/"+pn+".Pkg."+t.Format("20060102150405"))
        check(err)
        f, err := os.Create(path.Clean(pl)+"/"+pn+".Pkg"); check(err)
        defer f.Close()
        _, err = f.WriteString("k,v1,v2\n"); check(err)
        _, err = f.WriteString("package,"+name+","+version+"\n"); check(err)
        _, err = f.WriteString("from,"+from+","+retrieved+"\n"); check(err)
        _, err = f.WriteString("license,"+license+","+authors+"\n"); check(err)
        for k,v := range IMP {
                _, err := f.WriteString("r,"+k+","+v+"\n"); check(err)
        }
        for k,v := range PRO {
                _, err := f.WriteString("p,"+k+","+v+"\n"); check(err)
        }       

        f.Sync()
        err = os.Remove(path.Clean(pl)+"/"+pn+".Pkg."+t.Format("20060102150405")); check(err)
}

func getWorkspaceSettings(wk string) (map[string]string,map[string]string,map[string]string){
	WSV := make(map[string]string)
        REPOS := make(map[string]string)
        METAS := make(map[string]string)
	if _, err := os.Stat(path.Clean(wk)+"/Packages.Wrk"); err == nil {
		b, err := ioutil.ReadFile(path.Clean(wk)+"/Packages.Wrk")
		if err != nil {
			fmt.Print("Couldn't read Packages.Wrk")
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
	err := os.Rename(path.Clean(wk)+"/Packages.Wrk",path.Clean(wk)+"/Packages.Wrk."+t.Format("20060102150405"))
	check(err)
	f, err := os.Create(path.Clean(wk)+"/Packages.Wrk"); check(err)
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
	err = os.Remove(path.Clean(wk)+"/Packages.Wrk."+t.Format("20060102150405")); check(err)

}

func buildSourceList(wk string, WSV map[string]string,s []string) map[string]string {
	  files := make(map[string]string)

	  dstyle, _ := WSV["workspace-packages-dirstyle"]

	  

          if s[0]=="all" {
                fileInfo, _ := ioutil.ReadDir(path.Clean(wk))
                for _, file := range fileInfo {
		  if dstyle=="flat"{
		    if file.IsDir() {
			n:=file.Name()
                        if _, err := os.Stat(path.Clean(wk)+"/"+n+"/"+n+".Pkg"); err == nil {
	                      files[n] = path.Clean(wk)+"/"+n
		        }
		    }
		  }else if dstyle=="paths"{
		  }else{ // dstyle=="combined"
                          n:=file.Name()
                          if len(n)>4 {
                            if n[len(n)-4:]==".Pkg" {
                              files[n[0:len(n)-4]] = path.Clean(wk)
                            }
                          }
		  }
                }
          }else{
                for _, fn := range s {
		        pe:=dsExtend(fn,path.Clean(wk),WSV)

                        if _, err := os.Stat(pe+"/"+fn+".Pkg"); err == nil {
                                files[fn]=pe
                        }else{
                                fmt.Println("Package",fn,"Not Found, exiting.")
                                os.Exit(1)
                        }
                }
          }
    	  return files
}

func leStr( le string) string {
	var e string
        if le=="cr" {
                e="\r"
        }else if le=="crlf" {
                e="\r\n"
        }else if le=="nl"{
		e="\n"
        }else{  
                fmt.Println("Line ending style",le,"not understood.")
                os.Exit(1)
        }       
	return e
}

func initWorkspace(wk, le, ds string){
        fmt.Println("Initializing the workspace", wk)

	e:=leStr(le)

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

        if _, err := os.Stat(path.Clean(wk)+"/Packages.Wrk"); err != nil {
	        c := []byte("setting,value"+e+"workspace-module-line-ending,"+le+e+"workspace-packages-dirstyle,"+ds+e)
	        err := ioutil.WriteFile(path.Clean(wk)+"/Packages.Wrk", c, 0644)
	        if err != nil{
                        fmt.Println("Error Creating Packages.Wrk file.")
			os.Exit(1)
		}else{
                	fmt.Println("Created Packages.Wrk file.")
		}
        }else{
                fmt.Println("Packages.Wrk already exists in",wk,"exiting.")
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

func listPackages(  wkPtr *string, WSV map[string]string, tail []string){
        sPkgs := buildSourceList(*wkPtr,WSV,[]string{"all"})
	for i,j := range sPkgs{
		_,v,f,_,l,_,IMP,PRO:=getPackageSettings(i,j)
		fmt.Println(i,v,"license",l,"from",f)
		for i,j:=range(IMP){
			fmt.Println("  imports:",i,j)
		}
                for i,j:=range(PRO){
                        fmt.Println("  provides:",i,j)
                }
	}
}

func repoList( wk string){
                _,_,REPOS := getWorkspaceSettings(wk)
                for r,v:=range REPOS { fmt.Println(r,v)}
}

func metaList( wk string){
                _,METAS,_ := getWorkspaceSettings(wk)
                for m,v:=range METAS { fmt.Println(m,v)}
}


func repubList( init bool, wkPtr *string, WSV map[string]string, tail []string){
    var f *os.File = nil
    var err error
    rmTmp := false

    t := time.Now()
    subDir:=""
    le,_:=WSV["workspace-module-line-ending"]
    e:=leStr(le)
    if WSV["workspace-packages-dirstyle"] == "flat" { subDir = "/Index" }
    if WSV["workspace-packages-dirstyle"] == "path" { subDir = "/Index" }
    n:=path.Clean(tail[2])+subDir+"/Packages.Ndx"
    if len(tail)>2 {
        if _, err = os.Stat(n); err == nil {
	    if ! init {
		rmTmp = true
                err = os.Rename(n,n+"."+t.Format("20060102150405")); check(err)
                f, err = os.Create(n); check(err)
	    }else{ fmt.Println(" list already exists. Try 'repub'")}
 	}else{
	    if init {
		f, err = os.Create(n); check(err)
	    }else{ fmt.Println(" list not found. Try 'prepub'")}
	}

	if f != nil {

        	loc:=tail[1]
        	sPkgs := buildSourceList(*wkPtr,WSV,[]string{"all"})
        	_, err = f.WriteString("package,license,version,location"+e); check(err)
        	for i,j := range sPkgs{
        	        _,v,_,_,l,_,_,_:=getPackageSettings(i,j)
			if WSV["workspace-packages-dirstyle"] == "combined" {
        	        	_, err = f.WriteString(i+","+l+","+v+","+loc+e); check(err)
				}else if WSV["workspace-packages-dirstyle"] == "flat" {
        	                _, err = f.WriteString(i+","+l+","+v+","+loc+"/"+i+e); check(err)
			}else{
				// TODO: decide on path format
			}
        	}

	}
	if rmTmp {
		err = os.Remove(n+"."+t.Format("20060102150405")); check(err)
	}
    }else{
	fmt.Println("need \"mod repub <options> publicLocation localDirectory\"")
    }
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


func getRepoDefaultBranch( repo string ) (db string) {

        flatapiurl:="https://api.github.com/repos/"+repo+"/Index"
	combinedapiurl:="https://api.github.com/repos/"+repo

	resp, err := http.Get(flatapiurl);check(err)
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body); check(err)
	ans:=string(body)
	if ans[:23]=="{\"message\":\"Not Found\"," {

	        resp2, err2 := http.Get(combinedapiurl);check(err)
	        defer resp2.Body.Close()
	        body2, err2 := ioutil.ReadAll(resp2.Body); check(err2)
		ans = string(body2)
	}
	
	x:=strings.Split(ans,",")
	db = "master"
	for _,e:=range x {
		i:=strings.Split(e,":")
		if i[0]=="\"default_branch\"" {
			c:=len(i[1])
			db = i[1][1:c-1]
		}
	}
	
	

	return db
}

func getRepoStats(t string, REPOS map[string]string ) (string, int) {
	p:=""
	n:=0
                        x:=strings.SplitN(REPOS[t],"/",2)
                        repolist:=""
                        if x[0] == "github.com" {

				dBranch := getRepoDefaultBranch( x[1] )

                                flaturl:="https://raw.githubusercontent.com/"+x[1]+"/Index/"+dBranch+"/Packages.Ndx"
                                combinedurl:="https://raw.githubusercontent.com/"+x[1]+"/"+dBranch+"/Packages.Ndx"
                                resp, err := http.Get(flaturl);check(err)
                                defer resp.Body.Close()
                                body, err := ioutil.ReadAll(resp.Body); check(err)
                                if string(body[:32]) == "package,license,version,location" {
                                        repolist=string(nnl(body))
                                }else{
                                        fmt.Println(combinedurl)
                                        resp2, err2 := http.Get(combinedurl);check(err2)
                                        defer resp2.Body.Close()
                                        body2, err2 := ioutil.ReadAll(resp2.Body); check(err2)

                                        repolist=string(nnl(body2))
                                }
                        }else{
                                dir, err := ioutil.TempDir("", "pkg-"); check(err)

                                defer os.RemoveAll(dir)

                                fmt.Println("https://"+string(REPOS[t])+"/Index")

                                _, err = git.PlainClone(dir, false, &git.CloneOptions{
                                        URL: "https://"+string(REPOS[t])+"/Index",
                                })
                                if err != nil {
                                        log.Fatal(err)
                                }
                                
                                packagelist, err := os.Open(filepath.Join(dir, "Packages.Ndx"))
                                if err != nil {
                                        log.Fatal(err)
                                }
                                if 1==2 {io.Copy(os.Stdout, packagelist)}
                                
                        }
                        packages:=strings.Split(repolist,"\n")
                        pc:=0
                        for i,v := range packages{
                                if i > 0 {
                                        e:=strings.Split(v,",")
                                        if len(e)>1{
                                                pc++
                                                p=p+fmt.Sprintf("%s ",e[0])
                                        }
                                }
                        }
			n=pc

	return p,n 
}

func checkRepo( REPOS, METAS map[string]string, tail []string){
                t:=tail[1]
		  p:=""
		  c:=0
                  if _, ok := REPOS[t]; ! ok {
                        if r, ok := METAS[t]; ! ok {
                              fmt.Println(t,"not in workspace.")
                        }else{
			      p,c = getRepoStats(r, REPOS)
                        }
                  }else{
			p,c = getRepoStats(t, REPOS)
		  }
		  if c > 0 {
			fmt.Println(p,"\n",c,"packages.")
		  }else{
			fmt.Println("no packages.")
		  }

                  
}


func enrollPackage( wkPtr *string, WSV map[string]string, tail []string){
                sPkgs := buildSourceList(*wkPtr,WSV,[]string{"all"})
                nPkgs := strings.Split(tail[1],",")
                if len(nPkgs)!=1{
                        fmt.Println("Only enroll one package at a time")
                }else{
                        if _, ok := sPkgs[nPkgs[0]]; ! ok {       
                                //fmt.Println("Enrolling",nPkgs[0],sPkgs)
				le,_:=WSV["workspace-module-line-ending"]
				e:=leStr(le)
                                ds,_:=WSV["workspace-packages-dirstyle"]

		                c :=        "package,"+nPkgs[0]+",v0.0.0"+e+
                                            "license,NOLICENSE,UNKNOWN"+e

                  		if ds=="flat"{
		                	err := ioutil.WriteFile(path.Clean(*wkPtr)+"/"+nPkgs[0]+
						"/"+nPkgs[0]+".Pkg", []byte(c), 0644)
		                	if err != nil{
		                	        fmt.Println("Error Enrolling Package.")
		                	        os.Exit(1)
		                	}else{  
		                	        fmt.Println("Enrolled.")
		                	}       
                                }else if ds=="paths"{
                                }else{ // ds=="combined"
                                        err := ioutil.WriteFile(path.Clean(*wkPtr)+
                                                "/"+nPkgs[0]+".Pkg", []byte(c), 0644)
                                        if err != nil{
                                                fmt.Println("Error Enrolling Package.")
                                                os.Exit(1)
                                        }else{
                                                fmt.Println("Enrolled.")
                                        }     
                                }

                        }else{  
                                fmt.Println(nPkgs[0],"already in workspace.")
                        }
                }
}

func withdrawPackage( wkPtr *string, WSV map[string]string, tail []string){
        nPkgs := strings.Split(tail[1],",")
        if len(tail)<2 {
                fmt.Println("usage: mod remove PACKAGE")
        }else{
                if len(nPkgs)!=1 || tail[1]=="all"{
                        fmt.Println("Only remove one package at a time")
                }else{
                        pe:=dsExtend(tail[1],*wkPtr,WSV)
                        if _, err := os.Stat(pe+"/"+tail[1]+".Pkg"); err == nil {
				err := os.Remove(pe+"/"+tail[1]+".Pkg"); check(err)
                                fmt.Println("Package",tail[1],"removed.")
                        }else{
                                fmt.Println("Package",tail[1],"not found.")
                        }
                }
        }                
}

func relicensePackage(wkPtr *string, WSV map[string]string, tail []string){
        nPkgs := strings.Split(tail[1],",")
	if len(tail)<3 {
		fmt.Println("usage: mod relicense PACKAGE LICENSE")
	}else{
                if len(nPkgs)!=1 || tail[1]=="all"{
                        fmt.Println("Only relicense one package at a time")
                }else{
		        pe:=dsExtend(tail[1],*wkPtr,WSV)
			if _, err := os.Stat(pe+"/"+tail[1]+".Pkg"); err == nil {
				n,v,f,r,l,a,IMP,PRO:=getPackageSettings(tail[1],pe)
                        	fmt.Println("Relicensing", tail[1],"from",l,"to",tail[2])
                        	l=tail[2]
			        putPackageSettings(tail[1],pe,n,v,f,r,l,a,IMP,PRO)

			}else{
				fmt.Println("Package",tail[1],"not found.")
			}
                }
	}
}

func reauthorPackage(wkPtr *string, WSV map[string]string, tail []string){
        nPkgs := strings.Split(tail[1],",")
        if len(tail)<3 {
                fmt.Println("usage: mod reauthor PACKAGE \"AUTHORS\"")
        }else{
                if len(nPkgs)!=1 || tail[1]=="all"{
                        fmt.Println("Only reauthor one package at a time")
                }else{
                        pe:=dsExtend(tail[1],*wkPtr,WSV)
                        if _, err := os.Stat(pe+"/"+tail[1]+".Pkg"); err == nil {
                                n,v,f,r,l,a,IMP,PRO:=getPackageSettings(tail[1],pe)
                                fmt.Println("Relicensing", tail[1],"from",a,"to",tail[2])
                                a=tail[2]
                                putPackageSettings(tail[1],pe,n,v,f,r,l,a,IMP,PRO)

                        }else{
                                fmt.Println("Package",tail[1],"not found.")
                        }
                }
	}
}

func resourcePackage(wkPtr *string, WSV map[string]string, tail []string){
        nPkgs := strings.Split(tail[1],",")
        if len(tail)<3 {
                fmt.Println("usage: mod resource PACKAGE PATH")
        }else{
                if len(nPkgs)!=1 || tail[1]=="all"{
                        fmt.Println("Only resource one package at a time")
                }else{
                        pe:=dsExtend(tail[1],*wkPtr,WSV)
                        if _, err := os.Stat(pe+"/"+tail[1]+".Pkg"); err == nil {
                                n,v,f,r,l,a,IMP,PRO:=getPackageSettings(tail[1],pe)
                                fmt.Println("Changing origin of", tail[1],"from",f,"to",tail[2])
                                f=tail[2]
                                putPackageSettings(tail[1],pe,n,v,f,r,l,a,IMP,PRO)

                        }else{
                                fmt.Println("Package",tail[1],"not found.")
                        }
                }
        }
}

func incrementPackage(wkPtr *string, WSV map[string]string, tail []string){
        nPkgs := strings.Split(tail[1],",")
        if len(tail)<3 {
                fmt.Println("usage: mod increment VERSION (e.g. v1.0.3-rc6)")
        }else{
                if len(nPkgs)!=1 || tail[1]=="all"{
                        fmt.Println("Only increment one package at a time")
                }else{
                        pe:=dsExtend(tail[1],*wkPtr,WSV)
                        if _, err := os.Stat(pe+"/"+tail[1]+".Pkg"); err == nil {
                                n,v,f,r,l,a,IMP,PRO:=getPackageSettings(tail[1],pe)
                                fmt.Println("Incrementing", tail[1],"from",v,"to",tail[2])
                                v=tail[2]
                                putPackageSettings(tail[1],pe,n,v,f,r,l,a,IMP,PRO)

                        }else{
                                fmt.Println("Package",tail[1],"not found.")
                        }
                }
        }
}

func packageStatus(i,p string, WSV map[string]string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+"/"+i+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}

func latestPackage(i,p string, WSV map[string]string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+"/"+i+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}

func versionPackage(i,p string, WSV map[string]string){
	_,v,_,_,_,_,_,_:=getPackageSettings(i,p)
	fmt.Println("  ",i,v)
}

func dsExtend(i,p string, WSV map[string]string) string {
        pe:=p
        ds,_:=WSV["workspace-packages-dirstyle"]
        if ds=="flat" {
                pe=pe+"/"+i
        }
	return pe
}

func rehashPackage(i,p string, WSV map[string]string){
        var contents []byte
	pe:=dsExtend(i,p,WSV)

        contents, _ = ioutil.ReadFile(pe+".Pkg")//; check(err)
        fmt.Println("Rehashing", i, pe+".Pkg",":")

        // fmt.Println(string(contents))

	sha_256 := sha256.New()
	sha_256.Write(contents)

	if 1==2 { fmt.Printf("sha256:\t\t%x\n", sha_256.Sum(nil)) }

        n,v,f,r,l,a,IMP,PRO:=getPackageSettings(i,p)
	for j,h := range(PRO){
          
	  item, e := ioutil.ReadFile(p+"/"+j)
	  have:="-"
          sha_item := sha256.New()
	  if e == nil {
                sha_item.Write(item)
		have = fmt.Sprintf("%x",sha_item.Sum(nil))
		if have != h {
		   fmt.Println("   ",j,"to",have)
		   PRO[j]=have
		}else{
		   fmt.Println("   ",j,"unchanged.")
		}
	  }else{
		fmt.Println("  ",j,"not found! please fix.")
		os.Exit(1)
	  }
          

	}
	putPackageSettings(i,p,n,v,f,r,l,a,IMP,PRO)

}

func addtoPackage(i,p string, WSV map[string]string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+"/"+i+".Pkg")
                        fmt.Println("Adding to", p,":")
                        fmt.Println(string(contents))
}

func packageUpdates(i,p string, WSV map[string]string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+"/"+i+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}

func exactPackage(i,p string, WSV map[string]string){
        var contents []byte
                        contents, _ = ioutil.ReadFile(p+"/"+i+".Pkg")
                        fmt.Println("Status of", p,":")
                        fmt.Println(string(contents))
}

func packageProvider(i,p string, WSV map[string]string){
        _,_,f,_,_,_,_,_:=getPackageSettings(i,p)
        fmt.Println("  ",i,"origin is",f)
}



func doCommand( wkPtr, lePtr, dsPtr *string, tail []string) {

        if len(tail)==1 {
                if tail[0] == "init"     { initWorkspace(*wkPtr,*lePtr,*dsPtr)
          }else if tail[0] == "list"     { WSV,_,_ := getWorkspaceSettings(*wkPtr); listPackages(wkPtr,WSV,tail)
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
          }else if tail[0] == "enroll"    { enrollPackage(wkPtr,WSV,tail)
          }else if tail[0] == "withdraw"  { withdrawPackage(wkPtr,WSV,tail)
          }else if tail[0] == "relicense" { relicensePackage(wkPtr,WSV,tail)
          }else if tail[0] == "reauthor"  { reauthorPackage(wkPtr,WSV,tail)
          }else if tail[0] == "resource"  { resourcePackage(wkPtr,WSV,tail)
          }else if tail[0] == "increment" { incrementPackage(wkPtr,WSV,tail)
          }else if tail[0] == "prepub"    { repubList(true,wkPtr,WSV,tail)
          }else if tail[0] == "repub"     { repubList(false,wkPtr,WSV,tail)
	  }else{
	    sPkgs := buildSourceList(*wkPtr,WSV,strings.Split(tail[1],","))
	    for i, p := range sPkgs {
                    if tail[0]=="status"    { packageStatus(i,p,WSV)
              }else if tail[0]=="latest"    { latestPackage(i,p,WSV)
              }else if tail[0]=="version"   { versionPackage(i,p,WSV)
              }else if tail[0]=="rehash"    { rehashPackage(i,p,WSV)
              }else if tail[0]=="addto"     { addtoPackage(i,p,WSV)
              }else if tail[0]=="updates"   { packageUpdates(i,p,WSV)
              }else if tail[0]=="exact"     { exactPackage(i,p,WSV)
              }else if tail[0]=="provider"  { packageProvider(i,p,WSV)
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

    init                             Initialize a workspace
    repub                            Regenerate the publish list of packages
    prepub                           Generate the publish list of packages
    list                             List the packages in the workspace
    repolist                         List repos configured for workspace
    metalist                         List metarepos configured for workspace
    addrepo    <repo:path>           Add a repo (and path to that repo) to the workspace
    addmeta    <metarepo:repo>       Add a metarepo (and specific repo for the metarepo) to the workspace
    changerepo <repo:path>           Change the path for an existing repo for the workspace
    changemeta <metarepo:repo>       Change the specific repo for an existing metarepo for the workspace
    delrepo    <repo>                Remove a repo from the workspace
    delmeta    <metarepo>            Remove a metarepo from the workspace
 x  checkrepo  <repo>                Check the status of a repo for the workspace
    enroll     <package>             Enroll (register) a local package into the workspace
    withdraw   <package>             Withdraw (de-register) a local package from the workspace
 x  status     <package|all>         Check the status of a package or packages in the workspace and the repos
    relicense  <package> <license>   Change the license of a package in the workspace
    reauthor   <package> <authors>   Change the authors of a package in the workspace
    resource   <package> <authors>   Change the origin of a package in the workspace
    increment  <package> <version>   Increment the version number of a package in the workspace
 x  latest     <package|all>         Retrieve the latest version of a package from the repos to the workspace
    version    <package|all>         Show the version(s) of a package or packages in the workspace
    rehash     <package|all>         Update the hashes of local files in the workspace for the package or packages
 x  addto      <package> <file(s)>   Add local file(s) to the package in the workspace
 x  updates    <package|all>         Report updates (from repos) to a package or packages in the workspace
 x  exact      <package|all>         Retrieve a specific version of a package from the repos to the workspace
    provider   <package|all>         Report which repo provided (if any) the package or packages in the workspace
`)


        }

        flag.Parse()


        tail:= flag.Args()
        doCommand( wkPtr, lePtr, dsPtr, tail)

}


