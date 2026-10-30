BEGIN;

CREATE TABLE cefr_levels (
 code varchar(2) PRIMARY KEY,
 title varchar(80) NOT NULL,
 position smallint NOT NULL UNIQUE,
 available boolean NOT NULL DEFAULT false
);
INSERT INTO cefr_levels VALUES
 ('A1','Beginner',1,true),('A2','Elementary',2,true),('B1','Intermediate',3,false),
 ('B2','Upper Intermediate',4,false),('C1','Advanced',5,false),('C2','Proficiency',6,false);

CREATE TABLE units (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 level_code varchar(2) NOT NULL REFERENCES cefr_levels(code),
 title varchar(150) NOT NULL,
 position smallint NOT NULL CHECK (position > 0),
 UNIQUE(level_code,position)
);
CREATE TABLE lessons (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 unit_id uuid NOT NULL REFERENCES units(id),
 title varchar(150) NOT NULL,
 objective text NOT NULL,
 content text NOT NULL,
 position smallint NOT NULL CHECK (position > 0),
 UNIQUE(unit_id,position)
);
CREATE TABLE activities (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 lesson_id uuid NOT NULL UNIQUE REFERENCES lessons(id),
 title varchar(150) NOT NULL,
 kind varchar(30) NOT NULL CHECK (kind IN ('word_hunter','sentence_builder','memory_match')),
 instructions text NOT NULL,
 published boolean NOT NULL DEFAULT true,
 assessment boolean NOT NULL DEFAULT false,
 version integer NOT NULL DEFAULT 1
);
CREATE INDEX activities_lesson_idx ON activities(lesson_id);
CREATE TABLE questions (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 activity_id uuid NOT NULL REFERENCES activities(id),
 prompt text NOT NULL,
 explanation text NOT NULL,
 solution jsonb NOT NULL,
 position smallint NOT NULL,
 UNIQUE(activity_id,position)
);
CREATE TABLE answer_options (
 id varchar(40) NOT NULL,
 question_id uuid NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
 label varchar(120) NOT NULL,
 image varchar(40) NOT NULL DEFAULT '',
 PRIMARY KEY(question_id,id)
);
CREATE TABLE activity_attempts (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id uuid NOT NULL REFERENCES users(id),
 activity_id uuid NOT NULL REFERENCES activities(id),
 snapshot jsonb NOT NULL,
 answers jsonb,
 result jsonb,
 score integer CHECK (score BETWEEN 0 AND 100),
 started_at timestamptz NOT NULL DEFAULT now(),
 completed_at timestamptz
);
CREATE UNIQUE INDEX one_open_attempt ON activity_attempts(user_id,activity_id) WHERE completed_at IS NULL;
CREATE INDEX attempts_user_idx ON activity_attempts(user_id,completed_at);
CREATE TABLE user_progress (
 user_id uuid NOT NULL REFERENCES users(id),
 activity_id uuid NOT NULL REFERENCES activities(id),
 best_score integer NOT NULL CHECK (best_score BETWEEN 0 AND 100),
 completed_at timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY(user_id,activity_id)
);
CREATE TABLE xp_transactions (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id uuid NOT NULL REFERENCES users(id),
 source varchar(100) NOT NULL,
 amount integer NOT NULL CHECK(amount > 0 AND amount <= 200),
 created_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE(user_id,source)
);
CREATE TABLE streaks (
 user_id uuid PRIMARY KEY REFERENCES users(id),
 current_days integer NOT NULL DEFAULT 0 CHECK (current_days >= 0),
 best_days integer NOT NULL DEFAULT 0 CHECK (best_days >= current_days),
 last_day date
);
CREATE TABLE practice_days (
 user_id uuid NOT NULL REFERENCES users(id),
 day date NOT NULL,
 PRIMARY KEY(user_id,day)
);
CREATE TABLE achievements (
 code varchar(40) PRIMARY KEY,
 title varchar(100) NOT NULL,
 description text NOT NULL,
 threshold integer NOT NULL
);
INSERT INTO achievements VALUES
 ('first_activity','Primer paso','Aprueba tu primera actividad.',1),
 ('ten_activities','Explorador constante','Aprueba diez actividades diferentes.',10),
 ('first_assessment','Desafío superado','Aprueba una evaluación de unidad.',1),
 ('week_streak','Una semana aprendiendo','Practica siete días consecutivos.',7);
CREATE TABLE user_achievements (
 user_id uuid NOT NULL REFERENCES users(id),
 achievement_code varchar(40) NOT NULL REFERENCES achievements(code),
 earned_at timestamptz NOT NULL DEFAULT now(),
 PRIMARY KEY(user_id,achievement_code)
);
CREATE TABLE daily_missions (
 code varchar(40) PRIMARY KEY,
 title varchar(150) NOT NULL,
 target integer NOT NULL CHECK(target > 0),
 reward integer NOT NULL CHECK(reward > 0 AND reward <= 100)
);
INSERT INTO daily_missions VALUES ('three_activities','Aprueba tres actividades nuevas hoy',3,30);
CREATE TABLE user_daily_missions (
 user_id uuid NOT NULL REFERENCES users(id),
 mission_code varchar(40) NOT NULL REFERENCES daily_missions(code),
 day date NOT NULL,
 progress integer NOT NULL DEFAULT 0 CHECK(progress >= 0),
 rewarded boolean NOT NULL DEFAULT false,
 PRIMARY KEY(user_id,mission_code,day)
);
CREATE TABLE admin_audit (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 actor_id uuid NOT NULL REFERENCES users(id),
 action varchar(100) NOT NULL,
 resource_id uuid NOT NULL,
 created_at timestamptz NOT NULL DEFAULT now()
);

COMMIT;
