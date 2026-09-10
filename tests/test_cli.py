import contextlib
import io
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest
from unittest.mock import Mock, patch

from codexcommits import cli as cc


class CommitWorkflow(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix='cc-test-')
        self.repo = Path(self.temp.name)
        self.git('init', '-q')
        for key, value in [('user.name', 'Commit Test'), ('user.email', 'test@example.invalid'),
                           ('commit.gpgsign', 'false'), ('core.hooksPath', str(self.repo / '.git/hooks'))]:
            self.git('config', key, value)
        (self.repo / 'main.py').write_text('value = 1\n')
        self.git('add', 'main.py')
        self.git('commit', '-qm', 'chore: baseline')
        (self.repo / 'main.py').write_text('value = 2\n')
        self.git('add', 'main.py')
        self.head = self.git('rev-parse', 'HEAD')
        self.old_cwd = Path.cwd()
        os.chdir(self.repo)

    def tearDown(self):
        os.chdir(self.old_cwd)
        self.temp.cleanup()

    def git(self, *args):
        return subprocess.check_output(['git', '-C', str(self.repo), *args], stderr=subprocess.STDOUT).decode().strip()

    def run_flow(self, choices, generated='fix: update the value', effect=None):
        with patch.object(sys, 'argv', ['codexcommits']), patch.object(sys, 'stdin', Mock(isatty=lambda: True)), \
             patch('builtins.input', side_effect=choices), patch.object(cc, 'generate', return_value=generated, side_effect=effect) as model, \
             patch.object(cc.shutil, 'which', return_value='/mock/tool'), \
             contextlib.redirect_stdout(io.StringIO()), contextlib.redirect_stderr(io.StringIO()):
            result = cc.main()
            return result, model

    def test_only_staged_content_is_sent_and_committed(self):
        (self.repo / 'main.py').write_text('value = 3  # unstaged marker\n')
        result, model = self.run_flow(['y'])
        self.assertEqual(result, 0)
        diff = model.call_args.args[0]
        self.assertIn(b'+value = 2', diff)
        self.assertNotIn(b'unstaged marker', diff)
        self.assertEqual(self.git('show', 'HEAD:main.py'), 'value = 2')
        self.assertIn('unstaged marker', (self.repo / 'main.py').read_text())

    def test_cancel_preserves_head_and_staging(self):
        tree = self.git('write-tree')
        self.run_flow([''])
        self.assertEqual(self.head, self.git('rev-parse', 'HEAD'))
        self.assertEqual(tree, self.git('write-tree'))

    def test_edit_requires_confirmation_and_uses_literal_message(self):
        message = 'fix: preserve $(touch SHOULD_NOT_EXIST) literally'
        self.run_flow(['e', message, 'y'])
        self.assertEqual(self.git('log', '-1', '--format=%s'), message)
        self.assertFalse((self.repo / 'SHOULD_NOT_EXIST').exists())

    def test_changed_index_aborts(self):
        def mutate(diff, args):
            (self.repo / 'extra.txt').write_text('new staged change')
            self.git('add', 'extra.txt')
            return 'fix: update the value'
        with self.assertRaisesRegex(cc.Failure, 'changed'):
            self.run_flow([], effect=mutate)
        self.assertEqual(self.head, self.git('rev-parse', 'HEAD'))

    def test_changed_head_aborts(self):
        def mutate(diff, args):
            self.git('commit', '--allow-empty', '-qm', 'chore: concurrent commit')
            return 'fix: update the value'
        with self.assertRaisesRegex(cc.Failure, 'changed'):
            self.run_flow([], effect=mutate)

    def test_failing_hook_is_respected(self):
        hook = self.repo / '.git/hooks/pre-commit'
        hook.write_text('#!/bin/sh\nexit 1\n')
        hook.chmod(0o755)
        with self.assertRaisesRegex(cc.Failure, 'git commit failed'):
            self.run_flow(['y'])
        self.assertEqual(self.head, self.git('rev-parse', 'HEAD'))

    def test_empty_index_never_calls_codex(self):
        self.git('reset', '-q', 'HEAD')
        with patch.object(cc, 'generate') as model:
            with self.assertRaisesRegex(cc.Failure, 'No staged'):
                cc.snapshot(self.repo)
            model.assert_not_called()

    def test_generation_failure_never_commits(self):
        with self.assertRaisesRegex(cc.Failure, 'unavailable'):
            self.run_flow([], effect=cc.Failure('unavailable'))
        self.assertEqual(self.head, self.git('rev-parse', 'HEAD'))


if __name__ == '__main__':
    unittest.main(verbosity=2)
