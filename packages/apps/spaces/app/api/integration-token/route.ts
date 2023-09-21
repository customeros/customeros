import { NextRequest, NextResponse } from 'next/server';
import jwt from 'jsonwebtoken';

type ResponseData = {
  token: string;
};

const tokenData = {
  id: '650aec2b46c5eb01406d56b1',
  name: 'Alex Calinica',
};

const WORKSPACE_KEY = '39aefeb6-fdc3-4790-93c3-568e10b5e694';
const WORKSPACE_SECRET =
  '48184ecc1a1ff4aa6bd1b01b468d644ed58a84b2333ed99111f50e85d60a141a';
const options = {
  issuer: WORKSPACE_KEY,
  expiresIn: 7200,
};

export async function GET(req: NextRequest, res: NextResponse<ResponseData>) {
  const token = jwt.sign(tokenData, WORKSPACE_SECRET, options);
  return new Response(JSON.stringify(token));
}
