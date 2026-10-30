BEGIN;

ALTER TABLE cefr_levels ADD COLUMN description text NOT NULL DEFAULT '';
ALTER TABLE units ADD COLUMN description text NOT NULL DEFAULT '';
ALTER TABLE units ADD COLUMN published boolean NOT NULL DEFAULT true;
ALTER TABLE admin_audit ADD COLUMN target_key varchar(100) NOT NULL DEFAULT '';

CREATE TABLE educational_resources (
 key varchar(40) PRIMARY KEY,
 title varchar(150) NOT NULL,
 kind varchar(20) NOT NULL CHECK(kind IN ('image','audio')),
 url varchar(500) NOT NULL,
 attribution text NOT NULL DEFAULT '',
 active boolean NOT NULL DEFAULT true,
 created_at timestamptz NOT NULL DEFAULT now()
);
INSERT INTO educational_resources(key,title,kind,url,attribution) VALUES
 ('apple','Apple illustration','image','/images/apple.svg','Original EnglishQuest illustration'),
 ('banana','Banana illustration','image','/images/banana.svg','Original EnglishQuest illustration'),
 ('orange','Orange illustration','image','/images/orange.svg','Original EnglishQuest illustration'),
 ('book','Book illustration','image','/images/book.svg','Original EnglishQuest illustration'),
 ('pen','Pen illustration','image','/images/pen.svg','Original EnglishQuest illustration'),
 ('school','School illustration','image','/images/school.svg','Original EnglishQuest illustration'),
 ('breakfast','Breakfast illustration','image','/images/breakfast.svg','Original EnglishQuest illustration'),
 ('night','Night illustration','image','/images/night.svg','Original EnglishQuest illustration');
-- Keep the legacy image key for old snapshots/clients; add an enforced FK.
ALTER TABLE answer_options ADD COLUMN image_resource_key varchar(40) REFERENCES educational_resources(key);
UPDATE answer_options SET image_resource_key=NULLIF(image,'');

COMMIT;
