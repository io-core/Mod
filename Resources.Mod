MODULE Resources;
IMPORT SYSTEM, Kernel, Oberon;

TYPE
  Resource* = POINTER TO ResourceDesc;
  ResourceDesc* = RECORD
    value*: INTEGER;
    next: Resource
  END;

PROCEDURE New*(R: Resource; v: INTEGER) : Resource;
BEGIN
  IF R = NIL THEN NEW(R) END;
  IF R # NIL THEN
    R.value := v;
    R.next := NIL
  END
  RETURN R 
END New;

PROCEDURE Get*(R: Resource; x: INTEGER): INTEGER;
  VAR r,i: INTEGER; R1:Resource;
BEGIN i:=0;r:=0;
  R1:=R; WHILE R1 # NIL DO 
    IF i = x THEN r := R1.value END;
    R1:=R1.next;
    INC(i)
  END
  RETURN r 
END Get;

PROCEDURE Clear*(R: Resource);
BEGIN
  R.value:=0;
  R.next:=NIL
END Clear;

BEGIN
END Resources.
