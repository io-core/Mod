# Package
### Oberon source code package manager
A tool for managing packages of Oberon (and perhaps Component Pascal, and Modula2) source modules and resource files.

# NOT YET FUNCTIONAL --- WORK IN PROGRESS
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

```
