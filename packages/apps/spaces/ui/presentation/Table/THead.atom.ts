import { atomFamily, useRecoilState } from 'recoil';

const THeadAtom = atomFamily<boolean, string>({
  key: 'THeadAtom',
  default: false,
});

export const useTHeadState = (id: string) => {
  return useRecoilState(THeadAtom(id));
};
