require('dotenv').config();
const fs = require('fs');
const path = require('path');

const content = `// This file is auto-generated. Do not edit manually.
export const BASE_URL = '${process.env.DEFAULT_BASE_URL || ''}';
export const AUTHORIZE_PATH = '${process.env.DEFAULT_AUTHORIZE_PATH || ''}';
export const LOGOUT_PATH = '${process.env.DEFAULT_LOGOUT_PATH || ''}';
export const AUTH_PAGE = '${process.env.DEFAUTL_AUTH_PAGE || ''}';
`;

fs.writeFileSync(
  path.join(__dirname, '../src/constants.generated.ts'),
  content
);