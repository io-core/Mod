(* begin-documentation-node
################################################
  Compiling, Style Guide, and Contrib. Process
################################################

Obtaining the Code
==================

Check out the main project on a Linux or Mac based system

This repo uses submodules, to get the whole thing:

 * git clone --recursive https://github.com/io-core/io

After cloning the submodules may not be checked out to the main branch so after the above command you should:

 * cd io; bash makemain.sh


Navigating the Code
===================

At the top of the repo are some shell scripts, folders for artifacts of the build process, and a root folder.

The directory structures under root mimic a potential future deployed self-hosting target system.

Within root the src/github.com/io-core path provides a home for the io-core project submodules needed to build the target system.

Within /root/src/github.com/io-core:

.. list-table::
   :widths: 40 60
   :header-rows: 0

   * - Boot
     - source code for booting an IO target system and also shell scripts for building target artifacts on the host system.
   * - Build
     - compiler source code used to compile the system, including the compiler itself.
   * - Kernel, Files, and Modules
     - the inner core of the IO system.
   * - Oberon, Edit, and System
     - the outer core of the IO system.
   * - Assist, Attest, Crypto, Doc, Draw, Paint, Mod, ArgParse, Test
     - standard packages in the IO system.
   * - Examples, BASIC, C, Go, Pascal
     - supplemental packages in the IO system.

The Doc package also contains shell scripts used in the target system to generate documentation for the system.

All documentation for the system resides in comments within the source modules of the system.

Compiling the Code
==================

Building the artifacts of an IO system on a Linux or Mac host requres several external packages:

 * cd root/src/github.com/io-core/Boot; bash compile.sh

The system will use a previously built image to perform compilation, a new image will be generated 

Building the documetation of the IO system on a Linux or Mac host requires several more external packages:

 * cd root/src/github.com/io-core/Doc; bash genhtml.sh; make clean; make html

The documentation will be found in Doc/gen/html/

Coding Guidelines
=================

On your own system, anything goes. Code intended for a Pull Request to Integrated Oberon however should aspire to simplicity, clarity, and generality.

Commit Message Conventions
==========================

Please be more clear and descriptive than charlesap

Pull Request Process
====================

fork, hack, pull, merge, post!

Bug Report Checklist
====================

.. list-table::  
   :widths: 20 80 
   :header-rows: 1 
 
   * - Subsystem 
     - Issue Tracker 
   * - IO Project Issue Tracker
     - https://github.com/io-core/io/issues)
   * - doc
     - https://github.com/io-core/doc/issues
   * - Kernel
     - https://github.com/io-core/Kernel/issues
   * - Files
     - https://github.com/io-core/Files/issues
   * - Modules
     - https://github.com/io-core/Modules/issues
   * - Edit
     - https://github.com/io-core/Edit/issues
   * - Oberon
     - https://github.com/io-core/Oberon/issues
   * - System
     - https://github.com/io-core/System/issues
   * - Build
     - https://github.com/io-core/Build/issues
   * - Draw
     - https://github.com/io-core/Draw/issues



https://github.com/io-core/Mod/blob/main/Vet.Mod

end-documentation-node *)


(* begin-module-use-description
Module Vet handles:
* Nothing yet. Placeholder.

end-module-use-description *)
MODULE Vet; 
  IMPORT SYSTEM;

  CONST 

  TYPE 

  VAR 

  (* begin-section-description
## ---------- Generate Documentation
  end-section-description *)

  (* begin-procedure-description
---
**Foo** is a placeholder
  end-procedure-description *)
  PROCEDURE Foo;
  BEGIN 
  END Foo;



END Vet.
