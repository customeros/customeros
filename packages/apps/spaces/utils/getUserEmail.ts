import { User } from '@graphql/types';

export const getUserDisplayData = (user?: User | null): string => {
  if (!user) {
    // if user object does not exist the change was made by the system
    return 'default';
  }

  if (user?.emails?.[0]?.email) {
    return user.emails?.[0]?.email;
  }
  if (user?.firstName || user?.lastName) {
    return `${user?.firstName} ${user?.lastName}`.trim();
  }
  return 'Unknown';
};
