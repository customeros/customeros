import { toJS, computed, observable } from 'mobx';

import type { Store } from './_store';

export class Record<T extends object = object> {
  @observable accessor value: T;

  constructor(public store: Store<Record>, data: T) {
    this.value = data;
  }

  @computed
  get id() {
    return this.store.options.getId(this.value);
  }

  public commit(
    opts: {
      syncOnly?: boolean;
      onFailled?: () => void;
      onCompleted?: () => void;
    } = { syncOnly: false },
  ) {
    this.store.commit(this.id, opts);
  }

  public toRaw() {
    return toJS(this.value);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  static default(): any {
    return {};
  }
}

export type RecordFactoryClass<R extends Record> = {
  toRaw?(): R;
  default?(): R;
  new (store: Store<R>, data: object): R;
};
