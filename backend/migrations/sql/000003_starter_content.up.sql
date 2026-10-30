BEGIN;
-- Original introductory practice, not a complete CEFR course or certification.
DO $$
DECLARE topic jsonb; word jsonb; sentence jsonb; uid uuid; lid uuid; aid uuid; qid uuid;
 topics jsonb := '[
 {"level":"A1","position":1,"title":"Fruit garden","words":[{"en":"apple","es":"Manzana","image":"apple"},{"en":"banana","es":"Plátano","image":"banana"},{"en":"orange","es":"Naranja","image":"orange"}],"sentences":[{"es":"Me gustan las manzanas.","words":["I","like","apples."]},{"es":"Esto es un plátano.","words":["This","is","a","banana."]},{"es":"La naranja es dulce.","words":["The","orange","is","sweet."]}]},
 {"level":"A1","position":2,"title":"My classroom","words":[{"en":"book","es":"Libro","image":"book"},{"en":"pen","es":"Bolígrafo","image":"pen"},{"en":"school","es":"Escuela","image":"school"}],"sentences":[{"es":"Esto es un libro.","words":["This","is","a","book."]},{"es":"Tengo un bolígrafo.","words":["I","have","a","pen."]},{"es":"Voy a la escuela.","words":["I","go","to","school."]}]},
 {"level":"A2","position":1,"title":"Everyday routines","words":[{"en":"breakfast","es":"Desayuno","image":"breakfast"},{"en":"school","es":"Escuela","image":"school"},{"en":"night","es":"Noche","image":"night"}],"sentences":[{"es":"Desayuno a las siete.","words":["I","have","breakfast","at","seven."]},{"es":"Ella va a la escuela cada día.","words":["She","goes","to","school","every","day."]},{"es":"Nosotros estudiamos por la noche.","words":["We","study","at","night."]}]}
 ]'::jsonb;
 i integer; j integer; tokens jsonb; key jsonb; pairs jsonb; prompt text;
BEGIN
 FOR topic IN SELECT value FROM jsonb_array_elements(topics) LOOP
  INSERT INTO units(level_code,title,position) VALUES(topic->>'level',topic->>'title',(topic->>'position')::int) RETURNING id INTO uid;
  INSERT INTO lessons(unit_id,title,objective,content,position) VALUES(uid,'Discover words','Recognize three everyday English words.',
   'Look at each picture and say the English word. Read the word aloud before choosing. A picture helps connect meaning directly to English.',1) RETURNING id INTO lid;
  INSERT INTO activities(lesson_id,title,kind,instructions) VALUES(lid,'Word Hunter','word_hunter','Read the English word. Choose the matching picture. Select one answer for every question, then check your work.') RETURNING id INTO aid;
  FOR i IN 0..2 LOOP
   INSERT INTO questions(activity_id,prompt,explanation,solution,position) VALUES(aid,'Find: '||(topic->'words'->i->>'en'),
    (topic->'words'->i->>'en')||' means '||(topic->'words'->i->>'es')||'.',to_jsonb('c'||i::text),i+1) RETURNING id INTO qid;
   FOR j IN 0..2 LOOP
    word:=topic->'words'->j;
    INSERT INTO answer_options VALUES('c'||j::text,qid,word->>'es',word->>'image');
   END LOOP;
  END LOOP;
  INSERT INTO lessons(unit_id,title,objective,content,position) VALUES(uid,'Build a sentence','Use English word order in simple statements.',
   CASE WHEN topic->>'level'='A2' THEN 'In the present simple, add -s or -es with he/she/it: She goes to school. Use at with clock times and at night. Start statements with a subject, then a verb.'
    ELSE 'English statements usually follow subject + verb + object. For example: I like apples. Use This is a/an to identify one object. Start with a capital letter and finish with a full stop.' END,2) RETURNING id INTO lid;
  INSERT INTO activities(lesson_id,title,kind,instructions) VALUES(lid,'Sentence Builder','sentence_builder','Read the Spanish meaning. Select the English words in order. You can undo your choices before checking.') RETURNING id INTO aid;
  FOR i IN 0..2 LOOP
   sentence:=topic->'sentences'->i;tokens:=sentence->'words';key:='[]'::jsonb;
   FOR j IN 0..jsonb_array_length(tokens)-1 LOOP key:=key||jsonb_build_array('t'||j::text);END LOOP;
   INSERT INTO questions(activity_id,prompt,explanation,solution,position) VALUES(aid,sentence->>'es',
    (SELECT string_agg(value,' ' ORDER BY ordinality) FROM jsonb_array_elements_text(tokens) WITH ORDINALITY),key,i+1) RETURNING id INTO qid;
   FOR j IN 0..jsonb_array_length(tokens)-1 LOOP INSERT INTO answer_options VALUES('t'||j::text,qid,tokens->>j,'');END LOOP;
  END LOOP;
  INSERT INTO lessons(unit_id,title,objective,content,position) VALUES(uid,'Unit checkpoint','Recall the meaning of three words without a picture.',
   'Review the words from this unit. Try to recall the meaning before revealing each card. Match each English word with its Spanish meaning. This short checkpoint requires all pairs to be correct.',3) RETURNING id INTO lid;
  INSERT INTO activities(lesson_id,title,kind,instructions,assessment) VALUES(lid,'Memory Match','memory_match','Reveal cards and form English–Spanish pairs. Pairs are checked together when you finish. Use Reset if you want to change them.',true) RETURNING id INTO aid;
  pairs:='[]'::jsonb;
  FOR j IN 0..2 LOOP pairs:=pairs||jsonb_build_array(jsonb_build_array('e'||j::text,'s'||j::text));END LOOP;
  INSERT INTO questions(activity_id,prompt,explanation,solution,position) VALUES(aid,'Match the English words with their Spanish meanings.',
   'Read each English word aloud, then recall its meaning. Review the lesson if you need more practice.',pairs,1) RETURNING id INTO qid;
  FOR j IN 0..2 LOOP
   word:=topic->'words'->j;
   INSERT INTO answer_options VALUES('e'||j::text,qid,word->>'en',''),('s'||j::text,qid,word->>'es','');
  END LOOP;
 END LOOP;
END $$;
COMMIT;
