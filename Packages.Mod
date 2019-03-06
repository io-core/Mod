MODULE Packages;

  IMPORT SYSTEM, Files, Fonts, Resources;

  CONST 
    replace* = 0; insert* = 1; delete* = 2; unmark* = 3;  (*op-codes*)

  TYPE 
    Method* = POINTER TO MethodDesc;
    Component* = POINTER TO ComponentDesc;
    ComponentDesc* = RECORD
      x*, y*, w*, h*, d*, p*: INTEGER;
      sx*, sy*: REAL;
      sw*, sh*: INTEGER;
      selected*, marked*: BOOLEAN;
      do*: Method;
      next: Component
    END;

    Package* = POINTER TO PackageDesc;
    Notifier* = PROCEDURE (P: Package; op: INTEGER; beg, end: LONGINT);
    PackageDesc* = RECORD
      len*: LONGINT;
      changed*: BOOLEAN;
      flattened*: BOOLEAN;
      sel*, first: Component;
      notify*: Notifier
    END;

    Msg* = RECORD END ;
    WidMsg* = RECORD (Msg) w*: INTEGER END ;
    ColorMsg* = RECORD (Msg) col*: INTEGER END ;
    FontMsg* = RECORD (Msg) fnt*: Fonts.Font END ;
    

    bmpHeaderEntry = RECORD
        a: INTEGER;
        b: INTEGER;
        c: INTEGER;
        d: INTEGER
    END;

    jpgHeaderEntry = RECORD
        a: INTEGER;
        b: INTEGER;
        c: INTEGER;
        d: INTEGER
    END;

    pngHeaderEntry = RECORD
        a: INTEGER;
        b: INTEGER;
        c: INTEGER;
        d: INTEGER
    END;

    MethodDesc* = RECORD
(*        module*, allocator*: Name;
        new*: Modules.Command;             *)
        copy*: PROCEDURE (from, to: Component);
        draw*, change*: PROCEDURE (C: Component; VAR msg: Msg);
        selectable*: PROCEDURE (C: Component; x, y: INTEGER): BOOLEAN;
        read*: PROCEDURE (C: Component; VAR r: Files.Rider); (*; VAR c: Context); *)
        write*: PROCEDURE (C: Component; cno: INTEGER; VAR r: Files.Rider); (*; VAR c: Context);*)
        print*: PROCEDURE (C: Component; x, y: INTEGER)
      END ;

    
VAR 
 
  res*: INTEGER;
  width: INTEGER;
  new*: Component; 
  LineMethod*, CapMethod*, MacMethod* : Method;

  PROCEDURE New*(C: Component);
  BEGIN   new := C 
  END New;

  PROCEDURE Add*(P: Package; C: Component);
  BEGIN C.marked := FALSE; C.selected := TRUE; C.next := P.first;
    P.first := C; P.sel := C; P.changed := TRUE
  END Add;

  PROCEDURE ThisComponent*(P: Package; x, y: INTEGER): Component;
    VAR C: Component;
  BEGIN C := P.first;
    WHILE (C # NIL) & ~C.do.selectable(C, x ,y) DO C := C.next END ;
    RETURN C
  END ThisComponent;

  PROCEDURE SelectComponent*(P: Package; C: Component);
  BEGIN
    IF C # NIL THEN C.selected := TRUE; P.sel := C END
  END SelectComponent;

  PROCEDURE SelectArea*(P: Package; x0, y0, x1, y1: INTEGER);
    VAR C: Component; t: INTEGER;
  BEGIN C := P.first;
    IF x1 < x0 THEN t := x0; x0 := x1; x1 := t END ;
    IF y1 < y0 THEN t := y0; y0 := y1; y1 := t END ;
    WHILE C # NIL DO
      IF (x0 <= C.x) & (C.x + C.w <= x1) & (y0 <= C.y) & (C.y + C.h <= y1) THEN
        C.selected := TRUE; P.sel := C
      END ;
      C := C.next
    END ;
    IF P.sel # NIL THEN  END
  END SelectArea;

  PROCEDURE Draw*(P: Package; VAR m: Msg);
    VAR C: Component;
  BEGIN C := P.first;
    WHILE C # NIL DO C.do.draw(C, m); C := C.next END
  END Draw;

(*
  PROCEDURE List*(P: Package);
    VAR C: Component; tag: INTEGER;
  BEGIN C := P.first;
    WHILE C # NIL DO
      Texts.Write(XW, 9X); Texts.WriteHex(XW, ORD(C)); Texts.Write(XW, 9X);
      Texts.WriteInt(XW, C.x, 5); Texts.WriteInt(XW, C.y, 5); Texts.WriteInt(XW, C.w, 5); Texts.WriteInt(XW, C.h, 5);
      Texts.Write(XW, "/"); SYSTEM.GET(ORD(C)-8, tag); Texts.WriteHex(XW, tag);
      SYSTEM.GET(ORD(C)-4, tag); Texts.WriteHex(XW, tag); Texts.WriteLn(XW); C := C.next
    END ;
    Texts.Append(Oberon.Log, XW.buf)
  END List;
*)

  (*----------------procedures operating on selection -------------------*)

  PROCEDURE Deselect*(P: Package);
    VAR C: Component;
  BEGIN C := P.first; P.sel := NIL; 
    WHILE C # NIL DO C.selected := FALSE; C := C.next END
  END Deselect;

  PROCEDURE DrawSel*(P: Package; VAR m: Msg);
    VAR C: Component;
  BEGIN C := P.first;
    WHILE C # NIL DO
      IF C.selected THEN C.do.draw(C, m) END ;
      C := C.next
    END
  END DrawSel;

  PROCEDURE Change*(P: Package; VAR m: Msg);
    VAR C: Component;
  BEGIN C := P.first; P.changed := TRUE;
    WHILE C # NIL DO
      IF C.selected THEN C.do.change(C, m) END ;
      C := C.next
    END
  END Change;



  PROCEDURE WriteFile*(P: Package; name: ARRAY OF CHAR);
    VAR f: Files.File; w: Files.Rider; 
  BEGIN f := Files.New(name); Files.Set(w, f, 0); 
  (*  Files.Write(w, GraphFileId); InitContext(c); StoreElems(w, c, P.first);  *)
    Files.Register(f)
  END WriteFile;


  PROCEDURE SetWidth*(w: INTEGER);
  BEGIN width := w
  END SetWidth;


  PROCEDURE Open*(P: Package; name: ARRAY OF CHAR);
    VAR tag: CHAR;
          b: BYTE;
          i: INTEGER;
      f: Files.File; r: Files.Rider; 
  BEGIN P.first := NIL; P.sel := NIL; P.changed := FALSE; f := Files.Old(name);
    IF f # NIL THEN
      Files.Set(r, f, 0);
      Files.ReadByte(r, b); i:=b;
      Files.ReadByte(r, b); i:=i+(b*100H);
     
    ELSE res := 2
    END
  END Open;



  PROCEDURE Delete*(P: Package);
    VAR C, pred: Component;
  BEGIN P.sel := NIL; P.changed := TRUE; C := P.first;
    WHILE (C # NIL) & C.selected DO C := C.next END ;
    P.first := C;
    IF C # NIL THEN
      pred := C; C := C.next;
      WHILE C # NIL DO
        IF C.selected THEN pred.next := C.next ELSE pred := C END ;
        C := C.next
      END
    END
  END Delete;






  PROCEDURE InstallDrawMethods*(drawLine, drawCaption, drawMacro: PROCEDURE (C: Component; VAR msg: Msg));
  BEGIN LineMethod.draw := drawLine; CapMethod.draw := drawCaption; MacMethod.draw := drawMacro
  END InstallDrawMethods;

BEGIN 
END Packages.

