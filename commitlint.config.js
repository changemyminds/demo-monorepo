// Conventional Commits config — release-please relies on this to compute bumps.
// Scopes map to deployable units so commits route to the right release PR.
module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      [
        'feat',     // minor bump
        'fix',      // patch bump
        'perf',
        'refactor',
        'docs',
        'test',
        'build',
        'ci',
        'chore',
        'revert',
      ],
    ],
    // Encourage (not force) scoping commits by deployable unit.
    'scope-enum': [
      1,
      'always',
      [
        'service-a',
        'service-b',
        'service-c',
        'service-d',
        'libs',
        'deploy',
        'ci',
        'release',
        'repo',
      ],
    ],
    'subject-case': [0],
  },
};
