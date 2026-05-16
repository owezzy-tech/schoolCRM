const { defineConfig } = require('cz-git')

module.exports = defineConfig({
  extends: ['@commitlint/config-conventional'],
  rules: {
    'type-enum': [
      2,
      'always',
      ['feat', 'fix', 'docs', 'style', 'refactor', 'perf', 'test', 'build', 'ci', 'chore', 'revert']
    ],
    'scope-enum': [
      2,
      'always',
      ['frontend', 'rag', 'go', 'auth', 'sales', 'metrics', 'docs', 'config', 'deps', 'ci']
    ]
  },
  prompt: {
    useEmoji: false,
    alias: {
      fd: 'docs:frontend',
      fr: 'feat:rag',
      fg: 'feat:go',
      fx: 'fix'
    },
    messages: {
      type: 'Select the type of change:',
      scope: 'Select the scope of this change:',
      subject: 'Write a short, imperative summary:',
      body: 'Provide a longer description (optional):',
      breaking: 'Describe breaking changes (optional):',
      footer: 'List issue references or footer entries (optional):',
      confirmCommit: 'Confirm commit message?'
    },
    types: [
      { value: 'feat', name: 'feat:     a new feature' },
      { value: 'fix', name: 'fix:      a bug fix' },
      { value: 'docs', name: 'docs:     documentation only changes' },
      { value: 'style', name: 'style:    formatting, missing semi colons, etc' },
      { value: 'refactor', name: 'refactor: code change that neither fixes a bug nor adds a feature' },
      { value: 'perf', name: 'perf:     a code change that improves performance' },
      { value: 'test', name: 'test:     adding missing tests or correcting existing tests' },
      { value: 'build', name: 'build:    changes that affect the build system or dependencies' },
      { value: 'ci', name: 'ci:       changes to CI configuration files and scripts' },
      { value: 'chore', name: 'chore:    other changes that do not modify src or test files' },
      { value: 'revert', name: 'revert:   reverts a previous commit' }
    ],
    scopes: [
      { value: 'frontend', name: 'frontend: Angular app under api/frontends' },
      { value: 'rag', name: 'rag: Python RAG service under api/services/RAG' },
      { value: 'go', name: 'go: shared Go backend code' },
      { value: 'auth', name: 'auth: auth service and related packages' },
      { value: 'sales', name: 'sales: sales service and related packages' },
      { value: 'metrics', name: 'metrics: metrics service and related packages' },
      { value: 'docs', name: 'docs: documentation' },
      { value: 'config', name: 'config: repo or runtime configuration' },
      { value: 'deps', name: 'deps: dependency changes' },
      { value: 'ci', name: 'ci: git hooks or CI automation' }
    ],
    allowCustomScopes: true,
    allowEmptyScopes: false,
    upperCaseSubject: false,
    markBreakingChangeMode: true,
    breaklineNumber: 100,
    issuePrefixes: ['#']
  }
})
