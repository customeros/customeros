import Fuse from 'fuse.js';

import { type Organization } from '../Organization.dto';

export function removeAccents(str: string) {
  return str
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '');
}

export function indexAndSearch(arr: Organization[], term: string) {
  return new Fuse(arr, {
    keys: [
      { name: 'name', getFn: (o) => o.value.name },
      {
        name: 'website',
        getFn: (o) => o.value?.website ?? '',
      },
      {
        name: 'socials',
        getFn: (o) => o.value?.socialMedia?.[0]?.url,
      },
    ],
    threshold: 0.3,
    isCaseSensitive: false,
  })
    .search(removeAccents(term), { limit: 40 })
    .map((r) => r.item);
}
