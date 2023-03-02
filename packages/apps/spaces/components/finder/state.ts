import { atom } from 'recoil';

export const tableMode = atom<'PREVIEW' | 'MERGE_ORG' | 'MERGE_CONTACT'>({
  key: 'tableMode',
  default: 'PREVIEW',
});

export const selectedItemsIds = atom<Array<string>>({
  key: 'selectedItemsIds',
  default: [],
});
