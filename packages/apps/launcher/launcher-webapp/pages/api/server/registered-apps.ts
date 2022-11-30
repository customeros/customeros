// This is an example of to protect an API route
import { unstable_getServerSession } from "next-auth/next"
import { authOptions } from "../auth/[...nextauth]"

import type { NextApiRequest, NextApiResponse } from "next"
import {NextResponse} from "next/server";

export const config = {
  api: {
    externalResolver: true,
  },
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const session = await unstable_getServerSession(req, res, authOptions)

  if (session) {
    return NextResponse.rewrite(new URL('http://localhost:8070/customer-os/registered-apps', req.url))
  }

  res.send({
    error: "You must be signed in to view the protected content on this page.",
  })
}
