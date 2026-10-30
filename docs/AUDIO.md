# Original listening assets

EnglishQuest ships prerecorded WAV audio rather than browser-dependent speech
synthesis or a paid narration service. The spoken sentences are original
educational text. Narration is synthetic, generated with eSpeak NG (`en-us`,
135–180 words/minute in the new routes; 145 in the original clips), not a recording of a human teacher. It should be reviewed
for clarity by an English educator before a public course launch.

The 42 clips include six original A1/A2 recordings and 36 additional recordings
for guided A1–C2 units. Their filenames do not encode answers. Original short
transcripts are in post-submission explanations; the new scripts are in the private
curriculum authoring source, with relevant feedback in each question's explanation.
No script or solution is supplied in the public listening attempt before submission.
Faster narration alone is not evidence of academically validated CEFR difficulty.

To reproduce the files locally with eSpeak NG installed:

```bash
python3 scripts/generate-audio.py
```

Generation is an authoring operation, not a Render runtime dependency. The WAV
files are included in the frontend Static Site build. No generated user audio
or recordings are stored on the Go Web Service. Resource attribution is also
stored in PostgreSQL. Authors adding external audio must have permission to use
and distribute it, and must test that it is accessible in target browsers.
