# Package
### Oberon source code package manager
A tool for managing packages of Oberon (and perhaps Component Pascal, and Modula2) source modules and resource files.

```
Usage of ./package: package  <flags> commmand 

  Flags:

  -d string
    	workspace location (default "./")
  -e string
    	Local Module Line Ending Style (cr|crlf|nl|noch) (default "noch")
  -s string
    	Local Package Directory Style (combined|flat|paths) (default "combined")

  Commands:

    init
    repolist
    metalist
    addrepo <repo:path>
    addmeta <metarepo:repo>
    changerepo <repo:path>
    changemeta <metarepo:repo>
    delrepo <repo>
    delmeta <metarepo>
    checkrepo <repo>
    enroll <package>
    status <package|all>
    latest <package|all>
    rehash <package|all>
    addto <package>
    updates <package|all>
    exact <package|all>
    provider <package|all>
```
