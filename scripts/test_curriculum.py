import importlib.util
import json
import unittest
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
spec = importlib.util.spec_from_file_location('compiler', ROOT / 'scripts/build-curriculum.py')
compiler = importlib.util.module_from_spec(spec)
spec.loader.exec_module(compiler)


class CurriculumTests(unittest.TestCase):
    def setUp(self):
        self.source = json.loads(compiler.SOURCE.read_text())

    def test_all_routes_have_six_original_units(self):
        units = compiler.compile_course(self.source)
        self.assertEqual(36, len(units))
        self.assertEqual(180, sum(len(u['activities']) for u in units))
        self.assertEqual(36, len({u['reading'] for u in self.source}))
        self.assertEqual(36, len({u['listening'] for u in self.source}))

    def test_generated_migration_is_reproducible_and_additive(self):
        text = compiler.migration(self.source)
        self.assertEqual(text, compiler.TARGET.read_text())
        self.assertNotIn('DELETE FROM', text)
        self.assertNotIn('DROP TABLE', text)

    def test_questions_and_audio_are_real(self):
        for unit in compiler.compile_course(self.source):
            for activity in unit['activities']:
                self.assertGreater(len(activity['questions']), 0)
                for q in activity['questions']:
                    self.assertTrue(q['explanation'])
                    self.assertTrue(q['solution'])
                    for option in q['options']:
                        self.assertLessEqual(len(option['label'].encode()), 120, (unit['key'], option))
            path = ROOT / f"frontend/public/audio/course-{unit['key']}.wav"
            self.assertGreater(path.stat().st_size, 1000)

    def test_incomplete_or_duplicate_manifest_is_rejected(self):
        with self.assertRaises(AssertionError):
            compiler.compile_course(self.source[:-1])
        with self.assertRaises(AssertionError):
            compiler.compile_course(self.source + self.source[:1])


if __name__ == '__main__':
    unittest.main()
