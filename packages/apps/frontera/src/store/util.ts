import set from 'lodash/set';

import { Operation } from './types';

export function makePayload<T extends object>(operation: Operation) {
  return operation.diff.reduce((acc, curr) => {
    const path = curr.path;
    const value = curr.val;

    set(acc, path, value);

    return acc;
  }, {} as T);
}
