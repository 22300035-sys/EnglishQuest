#!/usr/bin/env python3
"""Reproducible authoring of original, prerecorded synthetic learning clips."""
import pathlib
import json
import subprocess
import wave

CLIPS = {
    "eq-listen-a1-1": "He has a blue bag.",
    "eq-listen-a1-2": "I have two apples.",
    "eq-listen-a1-3": "The book is on the desk.",
    "eq-listen-a2-1": "Mia usually walks to school, but today she is taking the bus.",
    "eq-listen-a2-2": "The library opens at half past nine in the morning.",
    "eq-listen-a2-3": "I'd like a sandwich and a glass of water, please.",
}

if __name__ == "__main__":
    target = pathlib.Path(__file__).resolve().parents[1] / "frontend/public/audio"
    target.mkdir(parents=True, exist_ok=True)
    for name, text in CLIPS.items():
        path = target / (name + ".wav")
        subprocess.run(["espeak-ng", "-v", "en-us", "-s", "145", "-w", str(path), text], check=True)
        with wave.open(str(path)) as audio:
            duration = audio.getnframes() / audio.getframerate()
            assert duration > 0.5
            print(f"{path.name}: {duration:.2f}s PCM audio")
    source = target.parents[2] / 'backend/curriculum/course.json'
    rates = {'A1': 135, 'A2': 145, 'B1': 155, 'B2': 165, 'C1': 175, 'C2': 180}
    for unit in json.loads(source.read_text()):
        path = target / f"course-{unit['key']}.wav"
        subprocess.run(['espeak-ng', '-v', 'en-us', '-s', str(rates[unit['level']]), '-w', str(path), unit['listening']], check=True)
        with wave.open(str(path)) as audio:
            duration = audio.getnframes() / audio.getframerate()
            assert duration > 0.5
            print(f'{path.name}: {duration:.2f}s original synthetic narration')
