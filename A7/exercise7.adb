with Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;
use  Ada.Text_IO, Ada.Integer_Text_IO, Ada.Numerics.Float_Random;

-- (Ada tabs = 3 spaces)

procedure exercise7 is

   Count_Failed   : exception;   -- Exception to be raised when counting fails
   Gen            : Generator;   -- Random number generator

   protected type Transaction_Manager (N : Positive) is
      entry       Finished;
      function    Commit return Boolean;
      procedure   Signal_Abort;
   private
      Finished_Gate_Open   : Boolean := False;
      Aborted              : Boolean := False;
      Will_Commit          : Boolean := True;
   end Transaction_Manager;
   protected body Transaction_Manager is
      entry Finished when Finished_Gate_Open or Finished'Count = N is
      begin
         if Finished'Count = N-1 then   -- when first enters
            Finished_Gate_Open := True;
            Will_Commit := not Aborted;

         end if;
         if Finished'Count = 0 then
            Finished_Gate_Open := False;
            Aborted := False;
         end if;
         ------------------------------------------
         -- PART 3: Complete the exit protocol here
         ------------------------------------------
      end Finished;

      procedure Signal_Abort is
      begin
         Aborted := True;
      end Signal_Abort;

      function Commit return Boolean is
      begin
         return Will_Commit;
      end Commit;
      
   end Transaction_Manager;



   
   procedure Unreliable_Slow_Add (x :in out Integer) is
   Error_Rate : Constant := 0.15;  -- (between 0 and 1)
   begin
      if Random(Gen)>Error_Rate then
         x := x + 10;
         delay Duration(4.0);
      else
         Put_Line("NEI");
         raise Count_Failed;
      end if;
   end Unreliable_Slow_Add;




   task type Transaction_Worker (Initial : Integer; Manager : access Transaction_Manager);
   task body Transaction_Worker is
      Num         : Integer   := Initial;
      Prev        : Integer   := Num;
      Round_Num   : Integer   := 0;
   begin
      Put_Line ("Worker" & Integer'Image(Initial) & " started");

      loop
         Put_Line ("Worker" & Integer'Image(Initial) & " started round" & Integer'Image(Round_Num));
         Round_Num := Round_Num + 1;
         begin
            Unreliable_Slow_Add(Num);
         exception
            when Count_Failed =>
              Manager.Signal_Abort;
         end;

         Manager.Finished;
        
         ---------------------------------------
         -- PART 2: Do the transaction work here          
         ---------------------------------------
         
         if Manager.Commit = True then
            Put_Line ("  Worker" & Integer'Image(Initial) & " comitting" & Integer'Image(Num));
         else
            Put_Line ("  Worker" & Integer'Image(Initial) &
                      " reverting from" & Integer'Image(Num) &
                      " to" & Integer'Image(Prev));
            -------------------------------------------
            -- PART 2: Roll back to previous value here
            -------------------------------------------
            Num := Prev;
         end if;

         Prev := Num;
         delay 0.5;

      end loop;
   end Transaction_Worker;

   Manager : aliased Transaction_Manager (3);

   Worker_1 : Transaction_Worker (0, Manager'Access);
   Worker_2 : Transaction_Worker (1, Manager'Access);
   Worker_3 : Transaction_Worker (2, Manager'Access);

begin
   Reset(Gen); -- Seed the random number generator
end exercise7;

