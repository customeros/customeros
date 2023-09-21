import { NextRequest, NextResponse } from 'next/server';
import jwt from 'jsonwebtoken';

type ResponseData = {
  token: string;
};

export async function GET(req: NextRequest, res: NextResponse<ResponseData>) {
  const WORKSPACE_KEY = process.env.INTEGRATION_APP_WORKSPACE_KEY;
  const WORKSPACE_SECRET = process.env.INTEGRATION_APP_WORKSPACE_SECRET;
  const options = {
    issuer: WORKSPACE_KEY,
    expiresIn: 7200,
  };

  const tokenData = {
    id: '650aec2b46c5eb01406d56b1',
    name: 'Alex Calinica',
  };

  if (!WORKSPACE_KEY || !WORKSPACE_SECRET) {
    return new Response('Missing integration app credentials', {
      status: 500,
    });
  }

  const token = jwt.sign(tokenData, WORKSPACE_SECRET, options);
  return new Response(JSON.stringify(token));
}
