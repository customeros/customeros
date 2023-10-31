import type { NextApiRequest, NextApiResponse } from 'next';
import jwt from 'jsonwebtoken';

const WORKSPACE_KEY = process.env.INTEGRATION_APP_WORKSPACE_KEY;
const WORKSPACE_SECRET = process.env.INTEGRATION_APP_WORKSPACE_SECRET;

type ResponseData =
  | {
      token: string;
    }
  | { message: string };

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<ResponseData>,
) {
  const tenant = req.query.tenant;

  if (!tenant) {
    return res.status(500).json({ message: 'Missing tenant query param' });
  }

  const tokenData = {
    id: tenant,
    name: tenant,
  };

  if (!WORKSPACE_KEY || !WORKSPACE_SECRET) {
    return res
      .status(500)
      .json({ message: 'Missing integration app credentials' });
  }

  const token = jwt.sign(tokenData, WORKSPACE_SECRET, {
    issuer: WORKSPACE_KEY,
    expiresIn: 3600,
    algorithm: 'HS256'
  });

  res.status(200).json({ token });
}
