--
-- Copyright 2019 Marco Helmich
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.
--

CREATE OR REPLACE PACKAGE BODY main AS

    PROCEDURE main IS
      li INT := 99;
      lstr varchar := 'this is a string';
      lstr2 varchar := 'hello';
    BEGIN
      dbms.print(li);
      dbms.print(lstr);
      IF li > 50 THEN
        dbms.print(50);
      END IF;

      dbms.print(47);

      WHILE li > 50 LOOP
        dbms.print(li);
        li := li - 1;
      END LOOP;

      IF lstr2 = 'hello' THEN
        dbms.print('string is hello :)');
      END IF;

      IF li = 50 THEN
        dbms.print('i is 50 :)');
      END IF;
    END;

END main;
/
