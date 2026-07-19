module.exports = {
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      [
        'feat',
        'fix',
        'docs',
        'test',
        'refactor',
        'perf',
        'chore',
        'security'
      ]
    ],
    'scope-enum': [
      2,
      'always',
      [
        'gateway',
        'scanner',
        'policy-engine',
        'rag-monitor',
        'redteam',
        'ui',
        'infra',
        'helm',
        'terraform',
        'docs'
      ]
    ]
  }
};
