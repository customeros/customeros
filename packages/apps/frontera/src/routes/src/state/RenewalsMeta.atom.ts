import { atom, useRecoilState } from 'recoil';
import { GetRenewalsQueryVariables } from '@renewals/graphql/getRenewals.generated';

interface RenewalsMeta {
  getRenewals: GetRenewalsQueryVariables;
  getUsers: {
    hasFetched: boolean;
  };
}

export const RenewalsMetaAtom = atom<RenewalsMeta>({
  key: 'RenewalsMeta',
  default: {
    getUsers: {
      hasFetched: false,
    },
    getRenewals: {
      pagination: {
        page: 1,
        limit: 40,
      },
    },
  },
});

export const useRenewalsMeta = () => {
  return useRecoilState(RenewalsMetaAtom);
};
