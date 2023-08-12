import { User } from '@graphql/types';

export const getUserDisplayData = (user?: User | null): string => {
  if (user?.emails?.[0]?.email) {
    return user.emails?.[0]?.email;
  }
  if (user?.firstName || user?.lastName) {
    return `${user?.firstName} ${user?.lastName}`.trim();
  }
  return 'Unknown';
};
