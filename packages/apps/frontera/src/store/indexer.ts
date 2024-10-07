/* eslint-disable @typescript-eslint/no-explicit-any */

import { reaction } from 'mobx';

export type IndexedFieldRecord = Record<string, (number | string)[]>;
export type IndexRecord = Record<string, IndexedFieldRecord>;

type IndexOptions = {
  primaryKey: string;
};

export class Index {
  record: IndexRecord;
  private primaryKey: string;

  constructor({ primaryKey }: IndexOptions) {
    this.record = {};
    this.primaryKey = primaryKey;
  }

  set(
    obj: object,
    fields: string[],
    options?: { getPrimaryKey: (obj: object) => number | string },
  ) {
    fields.forEach((field) => {
      let value = Index.getValueAtPath(obj, field);
      const pk =
        options?.getPrimaryKey(obj) ||
        Index.getValueAtPath(obj, this.primaryKey);

      // Handle special values
      if (value === null) {
        value = '__NULL__';
      } else if (value === undefined) {
        value = '__UNDEFINED__';
      } else if (value === '') {
        value = '__EMPTY__';
      }

      if (value !== undefined) {
        if (!this.record[field]) this.record[field] = {};
        if (!this.record[field][value]) this.record[field][value] = [];
        this.record[field][value].push(pk);
      }
    });
  }

  update(pk: number | string, field: string, oldValue: any, newValue: any) {
    // Handle special cases for oldValue
    if (oldValue === null) oldValue = '__NULL__';
    if (oldValue === undefined) oldValue = '__UNDEFINED__';
    if (oldValue === '') oldValue = '__EMPTY__';

    // Handle special cases for newValue
    if (newValue === null) newValue = '__NULL__';
    if (newValue === undefined) newValue = '__UNDEFINED__';
    if (newValue === '') newValue = '__EMPTY__';

    // remove index for old value
    if (this.record[field][oldValue]) {
      this.record[field][oldValue] = this.record[field][oldValue].filter(
        (existingPk) => existingPk !== pk,
      );

      // Clean up the index if no pks remain for the old value
      if (this.record[field][oldValue].length === 0) {
        delete this.record[field][oldValue];
      }
    }

    // index the new value
    if (!this.record[field][newValue]) {
      this.record[field][newValue] = [];
    }
    this.record[field][newValue].push(pk);
  }

  remove(
    obj: object,
    options?: { getPrimaryKey: (obj: object) => number | string },
  ) {
    Object.keys(this.record).forEach((field) => {
      let fieldValue = Index.getValueAtPath(obj, field);
      const pk =
        options?.getPrimaryKey(obj) ||
        Index.getValueAtPath(obj, this.primaryKey);

      // Handle special cases
      if (fieldValue === null) fieldValue = '__NULL__';
      if (fieldValue === undefined) fieldValue = '__UNDEFINED__';
      if (fieldValue === '') fieldValue = '__EMPTY__';

      if (this.record[field] && this.record[field][fieldValue]) {
        this.record[field][fieldValue] = this.record[field][fieldValue].filter(
          (_pk) => _pk !== pk,
        );

        // Clean up the index if no pks remain for the old value
        if (this.record[field][fieldValue].length === 0) {
          delete this.record[field][fieldValue];
        }
      }
    });
  }

  observeIndexedFields(obj: object, fields: string[]) {
    fields.forEach((field) => {
      reaction(
        () => Index.getValueAtPath(obj, field),
        (newValue, oldValue) => {
          const pk = Index.getValueAtPath(obj, this.primaryKey);

          if (newValue !== oldValue) {
            this.update(pk, field, oldValue, newValue);
          }
        },
      );
    });
  }

  private static getValueAtPath(obj: object, path: string): any {
    return path.split('.').reduce((acc, part) => acc && acc[part], obj as any);
  }
}
