import Fuse from 'fuse.js';

import type { Contact } from '../Contact.dto';

export function removeAccents(str: string) {
  return str
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '');
}

export function indexAndSearch(arr: Contact[], term: string) {
  return new Fuse(arr, {
    keys: [
      { name: 'name', getFn: (o) => o.name },
      {
        name: 'organization',
        getFn: (o) => o.value.primaryOrganizationName!,
      },
      {
        name: 'email',
        getFn: (o) => o.value?.emails?.[0]?.email || '',
      },
    ],
    threshold: 0.3,
    isCaseSensitive: false,
  })
    .search(removeAccents(term), { limit: 40 })
    .map((r) => r.item);
}
