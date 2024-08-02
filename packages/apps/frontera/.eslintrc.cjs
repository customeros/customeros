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
    '@stylistic',
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
    '@stylistic/no-multiple-empty-lines': ['error', { max: 1 }],
    '@stylistic/lines-between-class-members': [
      'error',
      {
        enforce: [
          { blankLine: 'always', prev: 'field', next: 'method' },
          { blankLine: 'always', prev: 'method', next: '*' },
        ],
      },
    ],
    '@stylistic/padding-line-between-statements': [
      'error',
      { blankLine: 'always', prev: '*', next: 'block-like' },
      {
        blankLine: 'always',
        prev: ['const', 'let', 'var', 'function'],
        next: 'expression',
      },
      {
        blankLine: 'always',
        prev: 'expression',
        next: ['const', 'let', 'var', 'function'],
      },
      {
        blankLine: 'always',
        prev: ['expression', 'block-like', 'const', 'let', 'var'],
        next: 'if',
      },
      { blankLine: 'always', prev: '*', next: 'return' },
    ],
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
        ignoreCase: true,
        customGroups: {
          value: {
            react: ['react', 'react-*'],
          },
          type: {
            react: ['react', 'react-*'],
          },
        },
        newlinesBetween: 'always',
        internalPattern: [
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
    'perfectionist/sort-jsx-props': [
      'error',
      {
        type: 'line-length',
        order: 'asc',
      },
    ],
    'perfectionist/sort-variable-declarations': [
      'error',
      {
        type: 'line-length',
        order: 'asc',
        ignoreCase: true,
      },
    ],
  },
};
