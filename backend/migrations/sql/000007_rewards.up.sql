BEGIN;
ALTER TABLE daily_missions ADD COLUMN metric varchar(30) NOT NULL DEFAULT 'new_completion'
 CHECK(metric IN ('new_completion','vocabulary','grammar','listening','reading','writing','perfect'));
ALTER TABLE daily_missions ADD COLUMN description text NOT NULL DEFAULT '';
UPDATE daily_missions SET description='Aprueba tres actividades distintas por primera vez en tu día local.' WHERE code='three_activities';
INSERT INTO daily_missions(code,title,target,reward,metric,description) VALUES
 ('vocabulary_practice','Explora vocabulario',2,20,'vocabulary','Aprueba dos actividades nuevas de vocabulario.'),
 ('grammar_practice','Construye buenas bases',2,20,'grammar','Aprueba dos actividades nuevas de gramática.'),
 ('listening_practice','Entrena tu oído',1,20,'listening','Aprueba una actividad nueva de comprensión auditiva.'),
 ('reading_practice','Lee para descubrir',1,20,'reading','Aprueba una actividad nueva de lectura o conversación contextual.'),
 ('writing_practice','Escribe con precisión',1,20,'writing','Aprueba una actividad nueva de escritura controlada.'),
 ('perfect_practice','Un desafío sin errores',1,10,'perfect','Aprueba una actividad nueva con una puntuación de 100%.');

CREATE TABLE xp_boosters (
 code varchar(40) PRIMARY KEY,
 title varchar(150) NOT NULL,
 multiplier integer NOT NULL CHECK(multiplier=2),
 duration_seconds integer NOT NULL CHECK(duration_seconds BETWEEN 60 AND 3600)
);
INSERT INTO xp_boosters VALUES ('double_xp_15m','Doble XP base durante 15 minutos',2,900);
CREATE TABLE user_boosters (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 user_id uuid NOT NULL REFERENCES users(id),
 booster_code varchar(40) NOT NULL REFERENCES xp_boosters(code),
 source varchar(100) NOT NULL,
 earned_at timestamptz NOT NULL DEFAULT now(),
 activated_at timestamptz,
 expires_at timestamptz,
 active boolean NOT NULL DEFAULT false,
 UNIQUE(user_id,source),
 CHECK((activated_at IS NULL AND expires_at IS NULL AND NOT active) OR (activated_at IS NOT NULL AND expires_at>activated_at))
);
CREATE UNIQUE INDEX one_active_booster ON user_boosters(user_id) WHERE active;
CREATE INDEX booster_inventory_idx ON user_boosters(user_id,earned_at);
-- Existing earned daily rewards receive inventory, not retroactive extra XP.
INSERT INTO user_boosters(user_id,booster_code,source)
 SELECT user_id,'double_xp_15m','daily:'||day::text FROM user_daily_missions
 WHERE mission_code='three_activities' AND rewarded ON CONFLICT DO NOTHING;

CREATE TABLE user_skill_stats (
 user_id uuid PRIMARY KEY REFERENCES users(id),
 words_practiced integer NOT NULL DEFAULT 0 CHECK(words_practiced>=0)
);
INSERT INTO user_skill_stats(user_id,words_practiced)
 SELECT p.user_id,sum(a.vocabulary_count) FROM user_progress p JOIN activities a ON a.id=p.activity_id GROUP BY p.user_id;
INSERT INTO achievements VALUES
 ('hundred_words','Cien prácticas de palabras','Practica cien palabras objetivo en primeras aprobaciones; no necesariamente cien palabras diferentes.',100),
 ('level_path','Ruta de nivel completada','Aprueba todas las actividades publicadas de un nivel. No certifica dominio CEFR.',1);
COMMIT;
