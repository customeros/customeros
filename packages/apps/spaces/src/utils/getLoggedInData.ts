import { Identity } from '@ory/client';

export function getUserName(identity: Identity): string {
  return identity.traits.email || identity.traits.username;
}
