import { atom } from 'recoil';

export const tableMode = atom<'PREVIEW' | 'MERGE' | 'ARCHIVE'>({
  key: 'tableMode',
  default: 'PREVIEW',
});

export const selectedItemsIds = atom<Array<string>>({
  key: 'selectedItemsIds',
  default: [],
});
