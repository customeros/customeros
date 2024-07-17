module.exports = {
  root: true,
  env: { browser: true, es2020: true, node: true },
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react-hooks/recommended',
    'eslint-config-prettier',
  ],
  ignorePatterns: [
    'dist',
    '.eslintrc.cjs',
    '**/*.generated.ts',
    '**/*.generated.tsx',
  ],
  parser: '@typescript-eslint/parser',
  plugins: [
    '@stylistic/js',
    'perfectionist',
    '@typescript-eslint',
    'eslint-plugin-prettier',
  ],
  rules: {
    'no-fallthrough': 'off',
    '@typescript-eslint/no-var-requires': 'off',
    '@typescript-eslint/prefer-ts-expect-error': 'warn',
    'prettier/prettier': 'error',
    //Stop changing this rule .... it's a waste of time
    'no-console': ['error', { allow: ['warn', 'error', 'info'] }],
    'react/display-name': 'off',
    'react-hooks/exhaustive-deps': 'off',
    '@stylistic/js/padding-line-between-statements': [
      'error',
      { blankLine: 'always', prev: '*', next: 'return' },
    ],
    '@stylistic/js/no-multiple-empty-lines': ['error', { max: 1 }],
    '@typescript-eslint/no-unused-vars': [
      'error',
      { varsIgnorePattern: '^_', ignoreRestSiblings: true, args: 'none' },
    ],
    'perfectionist/sort-imports': [
      'error',
      {
        type: 'line-length',
        order: 'asc',
        groups: [
          'type',
          'react',
          ['builtin', 'external'],
          'internal-type',
          'internal',
          ['parent-type', 'sibling-type', 'index-type'],
          ['parent', 'sibling', 'index'],
          'side-effect',
          'style',
          'object',
          'unknown',
        ],
        'custom-groups': {
          value: {
            react: ['react', 'react-*', 'next', 'next-*', 'next/**/*'],
          },
          type: {
            react: ['react', 'next', 'next-*', 'next/*'],
          },
        },
        'newlines-between': 'always',
        'internal-pattern': [
          '@ui/**',
          '@utils/**',
          '@shared/**',
          '@graphql/types',
          '@organization/**',
          '@organizations/**',
          '@customerMap/**',
        ],
      },
    ],
    'perfectionist/sort-named-imports': [
      'error',
      {
        type: 'line-length',
        order: 'asc',
      },
    ],
    'perfectionist/sort-interfaces': [
      'error',
      {
        type: 'line-length',
        order: 'asc',
      },
    ],
    'perfectionist/sort-object-types': [
      'error',
      {
        type: 'line-length',
        order: 'asc',
      },
    ],
  },
};
