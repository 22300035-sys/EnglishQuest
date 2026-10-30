BEGIN;
INSERT INTO educational_resources(key,title,kind,url,attribution) VALUES
 ('listen_a1_1','A1 listening clip 1','audio','/audio/eq-listen-a1-1.wav','Original EnglishQuest text; synthetic English narration generated with eSpeak NG'),
 ('listen_a1_2','A1 listening clip 2','audio','/audio/eq-listen-a1-2.wav','Original EnglishQuest text; synthetic English narration generated with eSpeak NG'),
 ('listen_a1_3','A1 listening clip 3','audio','/audio/eq-listen-a1-3.wav','Original EnglishQuest text; synthetic English narration generated with eSpeak NG'),
 ('listen_a2_1','A2 listening clip 1','audio','/audio/eq-listen-a2-1.wav','Original EnglishQuest text; synthetic English narration generated with eSpeak NG'),
 ('listen_a2_2','A2 listening clip 2','audio','/audio/eq-listen-a2-2.wav','Original EnglishQuest text; synthetic English narration generated with eSpeak NG'),
 ('listen_a2_3','A2 listening clip 3','audio','/audio/eq-listen-a2-3.wav','Original EnglishQuest text; synthetic English narration generated with eSpeak NG');

DO $$
DECLARE pack jsonb; game jsonb; question jsonb; uid uuid; lid uuid; aid uuid; qid uuid;
 idx integer; j integer; position integer; mode text; labels jsonb; solution jsonb;
 packs jsonb := '[
 {"level":"A1","unit":"Everyday adventures","games":[
 {"kind":"sound_detective","title":"Sound Detective","objective":"Identify colors, numbers and locations in short spoken sentences.","content":"Listen for a key word, not every word. Colors describe objects, numbers tell us how many, and on describes a position. Replay the recording when you need to.","instructions":"Play each recording and choose the answer. These are real prerecorded synthetic English clips; transcripts appear after you submit.","skill":"listening","vocabulary_count":3,"questions":[
 {"prompt":"What color is the bag?","audio":"listen_a1_1","options":["Blue","Green","Red"],"answer":0,"explanation":"He has a blue bag. The color is blue."},
 {"prompt":"How many apples does the speaker have?","audio":"listen_a1_2","options":["One","Two","Three"],"answer":1,"explanation":"I have two apples. The number is two."},
 {"prompt":"Where is the book?","audio":"listen_a1_3","options":["In the bag","Under the chair","On the desk"],"answer":2,"explanation":"The book is on the desk. On describes its position."}]},
 {"kind":"grammar_rescue","title":"Grammar Rescue","objective":"Correct subject–verb agreement in simple statements.","content":"Use have with I, you, we and they; use has with he, she and it. In the present simple, add -s or -es with he/she/it: She goes to school. Use are with they.","instructions":"Rewrite each sentence correctly. Preserve its meaning. Capitalization, repeated spaces and final punctuation do not affect your score.","skill":"grammar","questions":[
 {"prompt":"Correct: She go to school.","answer":"She goes to school.","explanation":"With she, go becomes goes in the present simple."},
 {"prompt":"Correct: I has a book.","answer":"I have a book.","explanation":"Use have with I: I have a book."},
 {"prompt":"Correct: They is happy.","answer":"They are happy.","explanation":"They takes are, not is."}]},
 {"kind":"spell_master","title":"Spell Master","objective":"Spell familiar English words using every letter once.","content":"Look at the meaning first. Say the English word, then build it with the letters. Some words have repeated letters: apple has two letters p.","instructions":"Build the English word from the letter tiles. Each tile can be used once, even when two tiles have the same letter.","skill":"writing","vocabulary_count":3,"questions":[
 {"prompt":"Spell the English word for: manzana","answer":"apple","explanation":"Apple is spelled a-p-p-l-e."},
 {"prompt":"Spell the English word for: escuela","answer":"school","explanation":"School is spelled s-c-h-o-o-l."},
 {"prompt":"Spell the English word for: plátano","answer":"banana","explanation":"Banana is spelled b-a-n-a-n-a."}]},
 {"kind":"conversation_quest","title":"Conversation Quest","objective":"Choose polite responses in short everyday conversations.","content":"Use a greeting to begin a conversation. How are you? asks about someone, not their name. Thank you expresses appreciation; You are welcome is a polite response.","instructions":"Read the conversation and choose a response that fits the situation.","skill":"reading","questions":[
 {"prompt":"Choose your response.","passage":"Alex: Hello! How are you?\nYou: ___","options":["My name is a book.","I am fine, thank you.","It is a banana."],"answer":1,"explanation":"I am fine, thank you answers How are you? politely."},
 {"prompt":"Choose your response.","passage":"Teacher: What is your name?\nYou: ___","options":["My name is Sam.","It is under the desk.","I am ten books."],"answer":0,"explanation":"My name is Sam answers a question about your name."},
 {"prompt":"Choose your response.","passage":"Friend: Thank you for your help!\nYou: ___","options":["I have three.","Good night, pencil.","You are welcome."],"answer":2,"explanation":"You are welcome is a polite reply to Thank you."}]},
 {"kind":"story_explorer","title":"Story Explorer","objective":"Find explicit information in a short everyday story.","content":"Read the whole story once to understand the situation. Then read each question and find the sentence that gives the information. You do not need to translate every word.","instructions":"Read the original short story and answer the questions. You may reread it while answering.","skill":"reading","passage":"Milo is a student. Every morning, he puts a blue book and a pen in his bag. He walks to school with his friend Ana. At school, they read a short story and write three sentences. At lunchtime, Milo eats an apple and drinks water. Ana has a banana. After school, they go home and practice English for ten minutes. Milo likes learning new words with his friend.","questions":[
 {"prompt":"What color is Milo’s book?","options":["Red","Blue","Yellow"],"answer":1,"explanation":"The story says that Milo puts a blue book in his bag."},
 {"prompt":"What does Milo eat at lunchtime?","options":["An apple","A banana","A sandwich"],"answer":0,"explanation":"Milo eats an apple; Ana has a banana."},
 {"prompt":"How long do they practice English after school?","options":["Three minutes","Thirty minutes","Ten minutes"],"answer":2,"explanation":"They practice English for ten minutes."}]},
 {"kind":"word_rain","title":"Word Rain","objective":"Recognize fruit vocabulary among unrelated words.","content":"Group words by meaning. Apple, banana and orange are fruits. Book and pen are classroom objects. Read the word before selecting it.","instructions":"Catch only the fruit words before the server-controlled time runs out. Words fall repeatedly. You can remove a selection before submitting. Reduced-motion mode shows static choices.","skill":"vocabulary","vocabulary_count":3,"seconds":90,"questions":[
 {"prompt":"Catch only the fruit words.","options":["apple","banana","orange","book","pen","school"],"answer":[0,1,2],"explanation":"Apple, banana and orange are fruits. The other words are not fruits."}]},
 {"kind":"english_rush","title":"English Rush","objective":"Apply vocabulary, word order and verb agreement in a mixed checkpoint.","content":"Review a few essentials: the plural of child is children; statements start with a subject; she takes goes in the present simple. Accuracy is more important than speed.","instructions":"Complete the mixed checkpoint within 120 seconds. The server checks the deadline. You may submit partial answers; unanswered questions are incorrect.","skill":"mixed","seconds":120,"assessment":true,"questions":[
 {"response":"choice","prompt":"Choose the plural of child.","options":["childs","children","childes"],"answer":1,"explanation":"Child has the irregular plural children."},
 {"response":"order","prompt":"Order the words: Estoy leyendo.","options":["I","am","reading."],"explanation":"The sentence is I am reading."},
 {"response":"text","prompt":"Complete with one word: She ___ (go) to school every day.","answer":"goes","explanation":"Use goes with she in the present simple."}]}
 ]},
 {"level":"A2","unit":"Plans and experiences","games":[
 {"kind":"sound_detective","title":"Sound Detective","objective":"Extract travel, time and order details from short recordings.","content":"Listen for contrasts such as usually versus today. Half past nine means 9:30. In a request, listen for the object after I would like.","instructions":"Listen to each original synthetic recording and answer. Replay as needed; the transcript is shown after submission.","skill":"listening","vocabulary_count":3,"questions":[
 {"prompt":"How is Mia traveling today?","audio":"listen_a2_1","options":["On foot","By bus","By bicycle"],"answer":1,"explanation":"Mia usually walks, but today she is taking the bus."},
 {"prompt":"When does the library open?","audio":"listen_a2_2","options":["At 9:00 a.m.","At 10:00 a.m.","At 9:30 a.m."],"answer":2,"explanation":"Half past nine in the morning means 9:30 a.m."},
 {"prompt":"Which drink does the speaker order?","audio":"listen_a2_3","options":["Water","Tea","Juice"],"answer":0,"explanation":"The speaker asks for a sandwich and a glass of water."}]},
 {"kind":"grammar_rescue","title":"Grammar Rescue","objective":"Correct basic past-simple and comparative forms.","content":"After did not, use the base form: did not go, not did not went. The past of buy is bought. Use than after a comparative adjective: bigger than.","instructions":"Correct each sentence while preserving its words and meaning where possible. This task checks the target correction, not every possible paraphrase.","skill":"grammar","questions":[
 {"prompt":"Correct: She did not went to school.","answer":"She did not go to school.","explanation":"Did not is followed by the base form go."},
 {"prompt":"Correct: I buyed a book yesterday.","answer":"I bought a book yesterday.","explanation":"Buy has the irregular past form bought."},
 {"prompt":"Correct: This bag is bigger then that bag.","answer":"This bag is bigger than that bag.","explanation":"Use than, not then, for comparisons."}]},
 {"kind":"spell_master","title":"Spell Master","objective":"Spell everyday routine and place words accurately.","content":"Pay attention to repeated letters and consonant groups. Read each meaning, recall the English word and check every letter before submitting.","instructions":"Build each English word with all its letter tiles. Identical letters can be chosen in either tile order.","skill":"writing","vocabulary_count":3,"questions":[
 {"prompt":"Spell the English word for: mañana (parte del día)","answer":"morning","explanation":"Morning is spelled m-o-r-n-i-n-g."},
 {"prompt":"Spell the English word for: biblioteca","answer":"library","explanation":"Library is spelled l-i-b-r-a-r-y."},
 {"prompt":"Spell the English word for: desayuno","answer":"breakfast","explanation":"Breakfast is spelled b-r-e-a-k-f-a-s-t."}]},
 {"kind":"conversation_quest","title":"Conversation Quest","objective":"Respond appropriately when buying tickets and asking for help.","content":"Polite requests often use could or would. Ask for clarification when you do not understand. Read the situation before choosing a response.","instructions":"Complete each conversation with the response that fits the context.","skill":"reading","questions":[
 {"prompt":"Choose your response.","passage":"Clerk: A single ticket or a return ticket?\nYou want to travel there and back.\nYou: ___","options":["A return ticket, please.","A single apple, please.","I returned yesterday."],"answer":0,"explanation":"A return ticket covers travel to a destination and back."},
 {"prompt":"Choose your response.","passage":"Visitor: Could you tell me where the library is?\nYou: ___","options":["It tasted good.","It is next to the school.","I have two tickets."],"answer":1,"explanation":"Next to the school gives the library’s location."},
 {"prompt":"Choose your response.","passage":"Someone gives instructions too quickly.\nYou: ___","options":["It is blue.","I will eat it.","Could you say that again, please?"],"answer":2,"explanation":"Ask politely for repetition when you need clarification."}]},
 {"kind":"story_explorer","title":"Story Explorer","objective":"Understand sequence and reasons in a short past-tense account.","content":"Words such as first, then and finally show order. Because introduces a reason. Use these clues to connect events in a short text.","instructions":"Read the story, then answer using information in the text.","skill":"reading","passage":"Last Saturday, Lea and Omar planned a picnic in the park. First, they bought bread, fruit and water at a small shop. They wanted to walk to the park, but it started to rain, so they took a bus. When they arrived, they found a covered table near the lake. They ate lunch and talked about a book they were reading. After an hour, the rain stopped. They walked around the lake and took some photographs. Finally, they went home by bus because they were tired. They enjoyed their day even though the weather changed their plans.","questions":[
 {"prompt":"What did Lea and Omar do first?","options":["They walked around the lake.","They bought food and water.","They went home."],"answer":1,"explanation":"First, they bought bread, fruit and water at a shop."},
 {"prompt":"Why did they take a bus to the park?","options":["It started to rain.","The shop was closed.","They did not like walking."],"answer":0,"explanation":"The rain changed their plan to walk to the park."},
 {"prompt":"What happened after the rain stopped?","options":["They bought a table.","They read at school.","They walked around the lake."],"answer":2,"explanation":"After the rain stopped, they walked around the lake and took photographs."}]},
 {"kind":"word_rain","title":"Word Rain","objective":"Recognize irregular past-simple verb forms.","content":"Some past forms are irregular: go becomes went, see becomes saw, and eat becomes ate. Do not add -ed to these verbs.","instructions":"Catch only the past-simple forms before the server deadline. Words reappear; avoid the base forms. Static choices are provided for reduced-motion users.","skill":"grammar","seconds":90,"questions":[
 {"prompt":"Catch only past-simple verb forms.","options":["went","saw","ate","go","see","eat"],"answer":[0,1,2],"explanation":"Went, saw and ate are past forms; go, see and eat are base forms."}]},
 {"kind":"english_rush","title":"English Rush","objective":"Apply everyday A2 grammar and word order in a mixed checkpoint.","content":"Use at for clock times. The past of go is went. In questions with did, use a subject and a base-form verb: Did you visit the park?", "instructions":"Solve the mixed checkpoint within 120 seconds. The API controls timing and scores; late attempts receive no XP.","skill":"mixed","seconds":120,"assessment":true,"questions":[
 {"response":"choice","prompt":"Complete: The lesson starts ___ nine o’clock.","options":["in","on","at"],"answer":2,"explanation":"Use at with a specific clock time."},
 {"response":"order","prompt":"Order the question: ¿Visitaste el parque?","options":["Did","you","visit","the","park?"],"explanation":"Use Did + subject + base-form verb: Did you visit the park?"},
 {"response":"text","prompt":"Complete with one word: We ___ (go) to the park yesterday.","answer":"went","explanation":"Go has the irregular past form went."}]}
 ]}
 ]'::jsonb;
BEGIN
 FOR pack IN SELECT value FROM jsonb_array_elements(packs) LOOP
  INSERT INTO units(level_code,title,description,position) SELECT pack->>'level',pack->>'unit','Integrated listening, grammar, spelling, reading and contextual practice.',COALESCE(max(u.position),0)+1 FROM units u WHERE level_code=pack->>'level' RETURNING id INTO uid;
  position:=0;
  FOR game IN SELECT value FROM jsonb_array_elements(pack->'games') LOOP
   position:=position+1;
   INSERT INTO lessons(unit_id,title,objective,content,position) VALUES(uid,game->>'title',game->>'objective',game->>'content',position) RETURNING id INTO lid;
   INSERT INTO activities(lesson_id,title,kind,instructions,time_limit_seconds,assessment,skill,vocabulary_count)
    VALUES(lid,game->>'title',game->>'kind',game->>'instructions',COALESCE((game->>'seconds')::int,0),COALESCE((game->>'assessment')::boolean,false),game->>'skill',COALESCE((game->>'vocabulary_count')::int,0)) RETURNING id INTO aid;
   idx:=0;
   FOR question IN SELECT value FROM jsonb_array_elements(game->'questions') LOOP
    idx:=idx+1;
    mode:=COALESCE(question->>'response',CASE game->>'kind' WHEN 'grammar_rescue' THEN 'text' WHEN 'spell_master' THEN 'letters' WHEN 'word_rain' THEN 'multi' ELSE 'choice' END);
    labels:=COALESCE(question->'options','[]'::jsonb);
    IF mode='letters' THEN SELECT jsonb_agg(letter ORDER BY n) INTO labels FROM regexp_split_to_table(question->>'answer','') WITH ORDINALITY AS chars(letter,n);END IF;
    IF mode='choice' THEN solution:=to_jsonb('o'||(question->>'answer'));
    ELSIF mode='order' THEN SELECT jsonb_agg('o'||n::text ORDER BY n) INTO solution FROM generate_series(0,jsonb_array_length(labels)-1) n;
    ELSIF mode='multi' THEN SELECT jsonb_agg('o'||value) INTO solution FROM jsonb_array_elements_text(question->'answer');
    ELSE solution:=question->'answer';END IF;
    INSERT INTO questions(activity_id,prompt,explanation,solution,position,response,passage,audio_key)
     VALUES(aid,question->>'prompt',question->>'explanation',solution,idx,mode,COALESCE(question->>'passage',game->>'passage',''),question->>'audio') RETURNING id INTO qid;
    IF jsonb_array_length(labels)>0 THEN FOR j IN 0..jsonb_array_length(labels)-1 LOOP INSERT INTO answer_options(id,question_id,label) VALUES('o'||j::text,qid,labels->>j);END LOOP;END IF;
   END LOOP;
  END LOOP;
 END LOOP;
END $$;
COMMIT;
