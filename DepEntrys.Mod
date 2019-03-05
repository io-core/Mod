MODULE DepEntrys;
IMPORT SYSTEM, Kernel, Oberon;

TYPE
  DepEntry* = POINTER TO DepEntryDesc;
  DepEntryDesc* = RECORD
    value*: INTEGER;
    next: DepEntry
  END;

PROCEDURE New*(E: DepEntry; v: INTEGER) : DepEntry;
BEGIN
  IF E = NIL THEN NEW(E) END;
  IF E # NIL THEN
    E.value := v;
    E.next := NIL
  END
  RETURN E 
END New;

PROCEDURE Get*(E: DepEntry; x: INTEGER): INTEGER;
  VAR r,i: INTEGER; E1:DepEntry;
BEGIN i:=0;r:=0;
  E1:=E; WHILE E1 # NIL DO 
    IF i = x THEN r := E1.value END;
    E1:=E1.next;
    INC(i)
  END
  RETURN r 
END Get;

PROCEDURE Clear*(E: DepEntry);
BEGIN
  E.value:=0;
  E.next:=NIL
END Clear;

BEGIN
END DepEntrys.
