/* eslint no-console: off */
const { copyFileSync, chmod } = require('fs');
const isCI = require('is-ci');

if (!isCI) {
  console.log('Setting up pre-commit hook');

  copyFileSync(
    process.cwd() + '/scripts/pre-commit',
    process.cwd() + '/../../../.git/hooks/pre-commit',
  );

  console.log('pre-commit installed succesfully');

  chmod(process.cwd() + '/../../../.git/hooks/pre-commit', 0o755, (err) => {
    if (err) throw err;
    console.log('permissions for pre-commit have changed successfully!');
  });
} else {
  console.log('CI env detected. Skipping setup of pre-commit-hook');
}
