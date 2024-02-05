import { atom, useRecoilState } from 'recoil';

interface CustomerLogo {
  logoUrl?: string | null;
  dimensions: {
    width: number;
    height: number;
  };
}

export const CustomerLogo = atom<CustomerLogo>({
  key: 'CustomerLogo',
  default: {
    logoUrl: null,
    dimensions: {
      width: 0,
      height: 0,
    },
  },
});

export const useCustomerLogo = () => {
  return useRecoilState(CustomerLogo);
};
