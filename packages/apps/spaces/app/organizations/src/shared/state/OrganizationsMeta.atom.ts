import { atom, useRecoilState } from 'recoil';
import { GetOrganizationsQueryVariables } from '../../graphql/getOrganizations.generated';

interface OrganizationsMeta {
  getUsers: {
    hasFetched: boolean;
  };
  getOrganization: GetOrganizationsQueryVariables;
}

export const OrganizationsMetaAtom = atom<OrganizationsMeta>({
  key: 'OrganizationsMeta',
  default: {
    getUsers: {
      hasFetched: false,
    },
    getOrganization: {
      pagination: {
        page: 1,
        limit: 40,
      },
    },
  },
});

export const useOrganizationsMeta = () => {
  return useRecoilState(OrganizationsMetaAtom);
};
