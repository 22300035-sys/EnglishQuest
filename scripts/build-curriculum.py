#!/usr/bin/env python3
"""Compile private original authoring data to an additive, reviewable migration."""
import argparse
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SOURCE = ROOT / 'backend/curriculum/course.json'
TARGET = ROOT / 'backend/migrations/sql/000008_guided_curriculum.up.sql'
PRODUCTION = {
    'A1': 'Write 40–80 words about yourself or a familiar person, using the unit situation as a model. Include clear basic facts and simple linked sentences.',
    'A2': 'Write 80–120 words about a comparable everyday situation. Describe events or plans, give a reason, and include a practical request or response.',
    'B1': 'Write 120–180 words explaining a comparable problem, its background and a proposed response. Connect events and reasons; consider an alternative.',
    'B2': 'Write 180–250 words presenting a reasoned response to the unit issue. Develop a position, engage with a counterargument, and qualify the conclusion.',
    'C1': 'Write 250–350 words for a named audience on the unit issue. Organise a sustained argument, distinguish evidence from inference, and calibrate register and certainty.',
    'C2': 'Write 300–450 words critically reformulating the unit argument for a different audience. Preserve its qualifications, analyse an alternative reading, and justify precise lexical and rhetorical choices.',
}


def choice(item, passage='', audio=''):
    prompt, labels, answer, explanation = item
    assert len(labels) == 3 and len(set(labels)) == 3 and 0 <= answer < 3
    return dict(prompt=prompt, options=[dict(id=f'o{i}', label=label) for i, label in enumerate(labels)],
                solution=f'o{answer}', explanation=explanation, passage=passage, audio_key=audio)


def compile_course(source):
    units = []
    seen = set()
    for unit in source:
        key, level = unit['key'], unit['level']
        assert key not in seen and level in ['A1', 'A2', 'B1', 'B2', 'C1', 'C2']
        seen.add(key)
        assert len(unit['reading']) >= 40 and len(unit['teaching']) >= 30
        reading = [choice(q, unit['reading']) for q in unit['comprehension']]
        grammar = [dict(prompt=q[0], solution=q[1], explanation=q[2], options=[], response='text') for q in unit['grammar']]
        conversation = [choice(q, unit['dialogue']) for q in unit['interaction']]
        listening = [choice(q, audio=f'course_{key}') for q in unit['listening_questions']]
        activities = []
        for title, kind, skill, questions, assessment, content in [
            ('Read and discover', 'story_explorer', 'reading', reading, False, unit['teaching']),
            ('Grammar in context', 'grammar_rescue', 'grammar', grammar, False, unit['teaching']),
            ('Respond with purpose', 'conversation_quest', 'reading', conversation, False, unit['teaching']),
            ('Listen for meaning', 'sound_detective', 'listening', listening, False,
             'Listen to the separate original recording. Listen once for the situation, then again for details and speaker intention. Replay as needed. The voice is synthetic; no speech is recorded from you.'),
            ('Unit review', 'conversation_quest', 'mixed', [], True, unit['teaching']),
        ]:
            # Untimed retrieval checkpoint; recognition only, not a free-writing grade.
            if assessment:
                questions = [choice(q, unit['reading']) for q in unit['comprehension']] + conversation
            activities.append(dict(title=title, kind=kind, skill=skill, questions=questions, assessment=assessment,
                                   content=content, instructions='Read the context and answer every question. Review the explanation after submitting.'
                                   if kind != 'grammar_rescue' else 'Follow the transformation exactly. Write the target sentence; case, repeated spaces and final punctuation are normalized. Other paraphrases are not automatically assessed.'))
        units.append(dict(key=key, level=level, title=unit['title'], objective=unit['objective'],
                          project_prompt=PRODUCTION[level] + '\nTopic: ' + unit['title'] + '\nFocus: ' + unit['objective'] + '\nSource context: ' + unit['reading'] + '\nOptional oral rehearsal: explain your response aloud, then note one change you would make. No microphone recording or automatic oral assessment is performed.',
                          project_rubric='Self-review (not an automatic grade): 1. Task and audience: did I answer the task and preserve relevant facts? 2. Organisation: are ideas linked and easy to follow? 3. Language: did I apply the unit focus with suitable vocabulary and register? 4. Accuracy and revision: what did I correct, and what needs teacher feedback?\nUnit language focus: ' + unit['teaching'],
                          description=unit['objective'] + ' Guided practice, not a CEFR certification.', activities=activities))
    for level in ['A1', 'A2', 'B1', 'B2', 'C1', 'C2']:
        assert sum(u['level'] == level for u in units) == 6, level
    return units


def migration(source):
    units = compile_course(source)
    payload = json.dumps(units, ensure_ascii=False, separators=(',', ':'))
    assert '$course$' not in payload
    return '''-- Generated from backend/curriculum/course.json; original guided practice.
-- Additive only: existing lessons, attempts, progress and administrator edits survive.
BEGIN;
CREATE TABLE unit_projects (
 unit_id uuid PRIMARY KEY REFERENCES units(id),
 prompt text NOT NULL,
 rubric text NOT NULL
);
CREATE TABLE user_unit_projects (
 user_id uuid NOT NULL REFERENCES users(id),
 unit_id uuid NOT NULL REFERENCES unit_projects(unit_id),
 body text NOT NULL DEFAULT '' CHECK(char_length(body)<=12000),
 reflection text NOT NULL DEFAULT '' CHECK(char_length(reflection)<=2000),
 updated_at timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY(user_id,unit_id)
);
DO $migration$
DECLARE course jsonb := $course$''' + payload + '''$course$::jsonb;
 pack jsonb; activity jsonb; question jsonb; option jsonb;
 uid uuid; lid uuid; aid uuid; qid uuid; n integer; p integer;
BEGIN
 FOR pack IN SELECT value FROM jsonb_array_elements(course) LOOP
  INSERT INTO units(level_code,title,description,position)
   SELECT pack->>'level',pack->>'title',pack->>'description',COALESCE(max(position),0)+1
   FROM units WHERE level_code=pack->>'level' RETURNING id INTO uid;
  INSERT INTO unit_projects(unit_id,prompt,rubric) VALUES(uid,pack->>'project_prompt',pack->>'project_rubric');
  INSERT INTO educational_resources(key,title,kind,url,attribution)
   VALUES('course_'||(pack->>'key'),(pack->>'level')||' '||(pack->>'title'),'audio',
    '/audio/course-'||(pack->>'key')||'.wav','Original EnglishQuest script; synthetic eSpeak NG narration');
  n:=0;
  FOR activity IN SELECT value FROM jsonb_array_elements(pack->'activities') LOOP
   n:=n+1;
   INSERT INTO lessons(unit_id,title,objective,content,position)
    VALUES(uid,activity->>'title',pack->>'objective',activity->>'content',n) RETURNING id INTO lid;
   INSERT INTO activities(lesson_id,title,kind,instructions,assessment,skill)
    VALUES(lid,activity->>'title',activity->>'kind',activity->>'instructions',
     (activity->>'assessment')::boolean,activity->>'skill') RETURNING id INTO aid;
   p:=0;
   FOR question IN SELECT value FROM jsonb_array_elements(activity->'questions') LOOP
    p:=p+1;
    INSERT INTO questions(activity_id,prompt,explanation,solution,position,response,passage,audio_key)
     VALUES(aid,question->>'prompt',question->>'explanation',question->'solution',p,
      COALESCE(question->>'response',''),COALESCE(question->>'passage',''),NULLIF(question->>'audio_key','')) RETURNING id INTO qid;
    FOR option IN SELECT value FROM jsonb_array_elements(question->'options') LOOP
     INSERT INTO answer_options(id,question_id,label) VALUES(option->>'id',qid,option->>'label');
    END LOOP;
   END LOOP;
  END LOOP;
 END LOOP;
 UPDATE cefr_levels SET available=true,description=CASE code
  WHEN 'A1' THEN 'Foundations: identity, routines, places, food, abilities and everyday plans.'
  WHEN 'A2' THEN 'Everyday interaction: past events, travel, comparisons, health, plans and experiences.'
  WHEN 'B1' THEN 'Independent practice: narratives, advice, work, community, media and decisions.'
  WHEN 'B2' THEN 'Extended practice: argument, evidence, negotiation, reporting, environment and culture.'
  WHEN 'C1' THEN 'Advanced practice: stance, implication, synthesis, register, evaluation and mediation.'
  WHEN 'C2' THEN 'Nuanced practice: rhetoric, ambiguity, irony, precision, conflicting sources and reformulation.'
 END || ' Guided route; completion is not CEFR accreditation.';
END $migration$;
COMMIT;
'''


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--check', action='store_true', help='Fail if committed migration differs; never write files.')
    args = parser.parse_args()
    result = migration(json.loads(SOURCE.read_text()))
    if args.check:
        assert TARGET.read_text() == result, 'Curriculum migration is out of sync.'
    else:
        TARGET.write_text(result)
    print('36 units; 180 guided lessons; six CEFR-labelled practice routes. No database was modified.')
