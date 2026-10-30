BEGIN;
ALTER TABLE activities DROP CONSTRAINT activities_kind_check;
ALTER TABLE activities ADD CONSTRAINT activities_kind_check CHECK(kind IN (
 'word_hunter','sentence_builder','memory_match','sound_detective','grammar_rescue',
 'spell_master','conversation_quest','story_explorer','word_rain','english_rush'));
ALTER TABLE activities ADD COLUMN time_limit_seconds integer NOT NULL DEFAULT 0 CHECK(time_limit_seconds=0 OR time_limit_seconds BETWEEN 30 AND 300);
ALTER TABLE activities ADD COLUMN vocabulary_count integer NOT NULL DEFAULT 0 CHECK(vocabulary_count BETWEEN 0 AND 100);
ALTER TABLE activities ADD COLUMN skill varchar(20) NOT NULL DEFAULT 'vocabulary' CHECK(skill IN ('vocabulary','grammar','listening','reading','writing','mixed'));
UPDATE activities SET skill='grammar' WHERE kind='sentence_builder';
UPDATE activities a SET vocabulary_count=(SELECT count(*) FROM questions q WHERE q.activity_id=a.id) WHERE kind='word_hunter';
UPDATE activities a SET vocabulary_count=(SELECT count(*)/2 FROM answer_options o JOIN questions q ON q.id=o.question_id WHERE q.activity_id=a.id) WHERE kind='memory_match';
ALTER TABLE questions ADD COLUMN response varchar(12) NOT NULL DEFAULT '' CHECK(response IN ('','choice','order','pairs','text','letters','multi'));
ALTER TABLE questions ADD COLUMN passage text NOT NULL DEFAULT '';
ALTER TABLE questions ADD COLUMN audio_key varchar(40) REFERENCES educational_resources(key);
COMMIT;
