/* eslint no-console: off */
const { copyFileSync, chmod, constants } = require('fs');
const isCI = require('is-ci');

if (!isCI) {
  console.log('Setting up pre-commit hook');

  copyFileSync(
    process.cwd() + '/scripts/pre-commit',
    process.cwd() + '/../../../.git/hooks/pre-commit',
  );

  console.log('pre-commit installed succesfully');

  chmod(
    process.cwd() + '/../../../.git/hooks/pre-commit',
    constants.X_OK,
    (err) => {
      if (err) throw err;
      console.log('permissions for pre-commit have changed successfully!');
    },
  );
} else {
  console.log('CI env detected. Skipping setup of pre-commit-hook');
}
