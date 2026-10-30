BEGIN;
CREATE TABLE user_level_unlocks (
 user_id uuid NOT NULL REFERENCES users(id),
 level_code varchar(2) NOT NULL REFERENCES cefr_levels(code),
 source varchar(12) NOT NULL CHECK(source IN ('progress','exam','legacy')),
 unlocked_at timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY(user_id,level_code)
);
-- Preserve access for learners who already started advanced content before gating.
INSERT INTO user_level_unlocks(user_id,level_code,source)
 SELECT DISTINCT t.user_id,u.level_code,'legacy' FROM activity_attempts t
 JOIN activities a ON a.id=t.activity_id JOIN lessons l ON l.id=a.lesson_id
 JOIN units u ON u.id=l.unit_id WHERE u.level_code IN ('B2','C1','C2')
 ON CONFLICT DO NOTHING;

CREATE TABLE access_exam_attempts (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id uuid NOT NULL REFERENCES users(id),
 target_level varchar(2) NOT NULL REFERENCES cefr_levels(code) CHECK(target_level IN ('B2','C1','C2')),
 snapshot jsonb NOT NULL,
 answers jsonb,
 result jsonb,
 started_at timestamptz NOT NULL DEFAULT now(),
 expires_at timestamptz NOT NULL,
 completed_at timestamptz,
 CHECK(expires_at>started_at)
);
CREATE UNIQUE INDEX one_open_access_exam ON access_exam_attempts(user_id,target_level) WHERE completed_at IS NULL;
CREATE INDEX access_exam_history ON access_exam_attempts(user_id,target_level,completed_at DESC);

CREATE TABLE course_diplomas (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id uuid NOT NULL UNIQUE REFERENCES users(id),
 recipient_name varchar(120) NOT NULL,
 issued_at timestamptz NOT NULL DEFAULT now(),
 course_title varchar(150) NOT NULL,
 activity_count integer NOT NULL CHECK(activity_count>0),
 activity_ids jsonb NOT NULL,
 level_counts jsonb NOT NULL
);
COMMIT;
