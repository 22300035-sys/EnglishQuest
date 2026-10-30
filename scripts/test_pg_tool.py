import importlib.util
import os
from pathlib import Path
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location("pg_tool", Path(__file__).with_name("pg-tool.py"))
pg_tool = importlib.util.module_from_spec(spec)
spec.loader.exec_module(pg_tool)


class ConnectionEnvironmentTest(unittest.TestCase):
    def test_credentials_are_not_passed_as_arguments(self):
        env = {"DATABASE_URL": "postgres://learner:p%40ss@localhost:55439/quest?sslmode=require", "PGHOST": "untrusted-host"}
        with patch.dict(os.environ, env, clear=True), patch("sys.argv", ["pg-tool.py", "DATABASE_URL", "pg_dump", "--format=custom"]), patch.object(pg_tool.subprocess, "run") as run:
            run.return_value.returncode = 0
            self.assertEqual(pg_tool.main(), 0)
            args, kwargs = run.call_args
            self.assertEqual(args[0], ["pg_dump", "--format=custom"])
            self.assertEqual(kwargs["env"]["PGPASSWORD"], "p@ss")
            self.assertEqual(kwargs["env"]["PGHOST"], "localhost")
            self.assertEqual(kwargs["env"]["PGPORT"], "55439")
            self.assertEqual(kwargs["env"]["PGSSLMODE"], "require")

    def test_rejects_unsupported_or_duplicate_parameters(self):
        for query in ["options=unsafe", "sslmode=require&sslmode=disable"]:
            with self.subTest(query=query), patch.dict(os.environ, {"DATABASE_URL": "postgres://localhost/quest?"+query}), patch("sys.argv", ["pg-tool.py", "DATABASE_URL", "psql"]):
                with self.assertRaises(ValueError):
                    pg_tool.main()

    def test_rejects_other_executables(self):
        with patch("sys.argv", ["pg-tool.py", "DATABASE_URL", "sh"]):
            with self.assertRaises(ValueError):
                pg_tool.main()

    def test_rejects_invalid_database_url(self):
        with patch.dict(os.environ, {"DATABASE_URL": "mysql://localhost/quest"}), patch("sys.argv", ["pg-tool.py", "DATABASE_URL", "psql"]):
            with self.assertRaises(ValueError):
                pg_tool.main()


if __name__ == "__main__":
    unittest.main()
