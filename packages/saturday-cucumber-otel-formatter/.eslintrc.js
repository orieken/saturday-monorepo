module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  parserOptions: {
    project: './tsconfig.json',
    ecmaVersion: 2020,
    sourceType: 'module',
  },
  plugins: [
    '@typescript-eslint'
  ],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended'
  ],
  rules: {

    // Complexity limits
    'complexity': ['error', 7], // Cyclomatic complexity limit
    'max-depth': ['error', 3], // Maximum block nesting
    'max-lines-per-function': ['error', {
      max: 10,
      skipBlankLines: true,
      skipComments: true
    }],
    'max-params': ['error', 3], // Limit function parameters

    // Clean code practices
    'max-len': ['error', {
      code: 10,
      ignoreComments: true,
      ignoreStrings: true,
      ignoreTemplateLiterals: true
    }],
    'no-magic-numbers': ['error', {
      ignore: [-1, 0, 1],
      ignoreArrayIndexes: true
    }],
    'prefer-const': 'error',
    'no-var': 'error',
    'no-multiple-empty-lines': ['error', { max: 1 }],

    // SOLID principles support
    '@typescript-eslint/interface-name-prefix': 'off',
    '@typescript-eslint/explicit-function-return-type': 'error',
    '@typescript-eslint/explicit-module-boundary-types': 'error',
    '@typescript-eslint/no-explicit-any': 'error',

    // Class and method structure
    'max-classes-per-file': ['error', 1],
    '@typescript-eslint/member-ordering': ['error', {
      default: [
        'static-field',
        'instance-field',
        'constructor',
        'static-method',
        'instance-method'
      ]
    }],

    // Code clarity
    '@typescript-eslint/no-unused-vars': ['error', {
      argsIgnorePattern: '^_',
      varsIgnorePattern: '^_'
    }],
    '@typescript-eslint/naming-convention': [
      'error',
      {
        selector: 'interface',
        format: ['PascalCase'],
        prefix: ['I']
      },
      {
        selector: 'class',
        format: ['PascalCase']
      },
      {
        selector: 'method',
        format: ['camelCase']
      },
      {
        selector: 'function',
        format: ['camelCase']
      },
      {
        selector: 'variable',
        format: ['camelCase', 'UPPER_CASE']
      }
    ],

    // Error prevention
    'no-console': ['warn', {
      allow: ['warn', 'error']
    }],
    '@typescript-eslint/no-floating-promises': 'error',
    '@typescript-eslint/await-thenable': 'error',
    '@typescript-eslint/no-misused-promises': 'error',

    // Documentation
    '@typescript-eslint/explicit-member-accessibility': ['error', {
      accessibility: 'explicit',
      overrides: {
        constructors: 'no-public'
      }
    }]
  },
  overrides: [
    {
      files: ['*.spec.ts', '*.test.ts'],
      rules: {
        'max-lines-per-function': 'off',
        'no-magic-numbers': 'off'
      }
    }
  ]
}