import type { NextApiRequest, NextApiResponse } from 'next';
import jwt from 'jsonwebtoken';

const WORKSPACE_KEY = process.env.INTEGRATION_APP_WORKSPACE_KEY;
const PRIVATE_KEY_VALUE = process.env.INTEGRATION_APP_PRIVATE_KEY_VALUE;

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

  if (!WORKSPACE_KEY || !PRIVATE_KEY_VALUE) {
    return res
      .status(500)
      .json({ message: 'Missing integration app credentials' });
  }

  const token = jwt.sign(tokenData, PRIVATE_KEY_VALUE, {
    issuer: WORKSPACE_KEY,
    expiresIn: 3600,
    algorithm: 'ES256'
  });

  res.status(200).json({ token });
}
