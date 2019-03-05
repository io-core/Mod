MODULE DepSets;

  IMPORT SYSTEM, Files, Fonts, DepEntrys;

  CONST 
    replace* = 0; insert* = 1; delete* = 2; unmark* = 3;  (*op-codes*)

  TYPE 
    Method* = POINTER TO MethodDesc;
    DepItem* = POINTER TO DepItemDesc;
    DepItemDesc* = RECORD
      x*, y*, w*, h*, d*, p*: INTEGER;
      sx*, sy*: REAL;
      sw*, sh*: INTEGER;
      selected*, marked*: BOOLEAN;
      do*: Method;
      next: DepItem
    END;

    DepSet* = POINTER TO DepSetDesc;
    Notifier* = PROCEDURE (S: DepSet; op: INTEGER; beg, end: LONGINT);
    DepSetDesc* = RECORD
      len*: LONGINT;
      changed*: BOOLEAN;
      flattened*: BOOLEAN;
      sel*, first: DepItem;
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
        copy*: PROCEDURE (from, to: DepItem);
        draw*, change*: PROCEDURE (I: DepItem; VAR msg: Msg);
        selectable*: PROCEDURE (I: DepItem; x, y: INTEGER): BOOLEAN;
        read*: PROCEDURE (I: DepItem; VAR r: Files.Rider); (*; VAR c: Context); *)
        write*: PROCEDURE (I: DepItem; cno: INTEGER; VAR r: Files.Rider); (*; VAR c: Context);*)
        print*: PROCEDURE (I: DepItem; x, y: INTEGER)
      END ;

    
VAR 
 
  res*: INTEGER;
  width: INTEGER;
  new*: DepItem; 
  LineMethod*, CapMethod*, MacMethod* : Method;

  PROCEDURE New*(I: DepItem);
  BEGIN   new := I 
  END New;

  PROCEDURE Add*(S: DepSet; I: DepItem);
  BEGIN I.marked := FALSE; I.selected := TRUE; I.next := S.first;
    S.first := I; S.sel := I; S.changed := TRUE
  END Add;

  PROCEDURE ThisDepItem*(S: DepSet; x, y: INTEGER): DepItem;
    VAR I: DepItem;
  BEGIN I := S.first;
    WHILE (I # NIL) & ~I.do.selectable(I, x ,y) DO I := I.next END ;
    RETURN I
  END ThisDepItem;

  PROCEDURE SelectDepItem*(S: DepSet; I: DepItem);
  BEGIN
    IF I # NIL THEN I.selected := TRUE; S.sel := I END
  END SelectDepItem;

  PROCEDURE SelectArea*(S: DepSet; x0, y0, x1, y1: INTEGER);
    VAR I: DepItem; t: INTEGER;
  BEGIN I := S.first;
    IF x1 < x0 THEN t := x0; x0 := x1; x1 := t END ;
    IF y1 < y0 THEN t := y0; y0 := y1; y1 := t END ;
    WHILE I # NIL DO
      IF (x0 <= I.x) & (I.x + I.w <= x1) & (y0 <= I.y) & (I.y + I.h <= y1) THEN
        I.selected := TRUE; S.sel := I
      END ;
      I := I.next
    END ;
    IF S.sel # NIL THEN  END
  END SelectArea;

  PROCEDURE Draw*(S: DepSet; VAR m: Msg);
    VAR I: DepItem;
  BEGIN I := S.first;
    WHILE I # NIL DO I.do.draw(I, m); I := I.next END
  END Draw;

(*
  PROCEDURE List*(S: DepSet);
    VAR I: DepItem; tag: INTEGER;
  BEGIN I := S.first;
    WHILE I # NIL DO
      Texts.Write(XW, 9X); Texts.WriteHex(XW, ORD(I)); Texts.Write(XW, 9X);
      Texts.WriteInt(XW, I.x, 5); Texts.WriteInt(XW, I.y, 5); Texts.WriteInt(XW, I.w, 5); Texts.WriteInt(XW, I.h, 5);
      Texts.Write(XW, "/"); SYSTEM.GET(ORD(I)-8, tag); Texts.WriteHex(XW, tag);
      SYSTEM.GET(ORD(I)-4, tag); Texts.WriteHex(XW, tag); Texts.WriteLn(XW); I := I.next
    END ;
    Texts.Append(Oberon.Log, XW.buf)
  END List;
*)

  (*----------------procedures operating on selection -------------------*)

  PROCEDURE Deselect*(S: DepSet);
    VAR I: DepItem;
  BEGIN I := S.first; S.sel := NIL; 
    WHILE I # NIL DO I.selected := FALSE; I := I.next END
  END Deselect;

  PROCEDURE DrawSel*(S: DepSet; VAR m: Msg);
    VAR I: DepItem;
  BEGIN I := S.first;
    WHILE I # NIL DO
      IF I.selected THEN I.do.draw(I, m) END ;
      I := I.next
    END
  END DrawSel;

  PROCEDURE Change*(S: DepSet; VAR m: Msg);
    VAR I: DepItem;
  BEGIN I := S.first; S.changed := TRUE;
    WHILE I # NIL DO
      IF I.selected THEN I.do.change(I, m) END ;
      I := I.next
    END
  END Change;



  PROCEDURE WriteFile*(S: DepSet; name: ARRAY OF CHAR);
    VAR f: Files.File; w: Files.Rider; 
  BEGIN f := Files.New(name); Files.Set(w, f, 0); 
  (*  Files.Write(w, GraphFileId); InitContext(c); StoreElems(w, c, S.first);  *)
    Files.Register(f)
  END WriteFile;


  PROCEDURE SetWidth*(w: INTEGER);
  BEGIN width := w
  END SetWidth;


  PROCEDURE Open*(S: DepSet; name: ARRAY OF CHAR);
    VAR tag: CHAR;
          b: BYTE;
          i: INTEGER;
      f: Files.File; r: Files.Rider; 
  BEGIN S.first := NIL; S.sel := NIL; S.changed := FALSE; f := Files.Old(name);
    IF f # NIL THEN
      Files.Set(r, f, 0);
      Files.ReadByte(r, b); i:=b;
      Files.ReadByte(r, b); i:=i+(b*100H);
     
    ELSE res := 2
    END
  END Open;



  PROCEDURE Delete*(S: DepSet);
    VAR I, pred: DepItem;
  BEGIN S.sel := NIL; S.changed := TRUE; I := S.first;
    WHILE (I # NIL) & I.selected DO I := I.next END ;
    S.first := I;
    IF I # NIL THEN
      pred := I; I := I.next;
      WHILE I # NIL DO
        IF I.selected THEN pred.next := I.next ELSE pred := I END ;
        I := I.next
      END
    END
  END Delete;






  PROCEDURE InstallDrawMethods*(drawLine, drawCaption, drawMacro: PROCEDURE (I: DepItem; VAR msg: Msg));
  BEGIN LineMethod.draw := drawLine; CapMethod.draw := drawCaption; MacMethod.draw := drawMacro
  END InstallDrawMethods;

BEGIN 
END DepSets.

