/* eslint-disable @typescript-eslint/no-explicit-any */
import { Filter, ComparisonOperator } from '@graphql/types';

import { type Index, type IndexedFieldRecord } from './indexer';

type PrimaryKey = string | number;

export class QueryClient {
  private index: Index;

  constructor(index: Index) {
    this.index = index;
  }

  query(filter: Filter): PrimaryKey[] {
    // Base case: If this is a simple filter (FilterItem)
    if (filter.filter) {
      const { property, operation, value, caseSensitive } = filter.filter;

      return this.queryIndex(property, operation!, value, caseSensitive);
    }

    // Handle AND conditions (intersection of all sub-filters)
    if (filter.AND && filter.AND.length > 0) {
      return filter.AND.map((subFilter) => this.query(subFilter)).reduce(
        (acc, result) => intersect(acc, result),
      );
    }

    // Handle OR conditions (union of all sub-filters)
    if (filter.OR && filter.OR.length > 0) {
      return filter.OR.map((subFilter) => this.query(subFilter)).reduce(
        (acc, result) => union(acc, result),
      );
    }

    // Handle NOT condition (negation of the sub-filter)
    if (filter.NOT) {
      const positiveResults = this.query(filter.NOT);
      const allIds = Object.keys(this.index)
        .flatMap((field) => Object.values(this.index.record[field]))
        .flat(); // Get all possible IDs

      return allIds.filter((pk) => !positiveResults.includes(pk));
    }

    // Default case: return an empty result set
    return [];
  }

  private queryIndex(
    field: string,
    operator: ComparisonOperator,
    value: any,
    caseSensitive: boolean | null = true,
  ): PrimaryKey[] {
    const fieldIndex = this.index.record[field];

    if (!fieldIndex) return [];

    let pks: PrimaryKey[] = [];

    switch (operator) {
      case ComparisonOperator.Eq:
        pks = handleEq(fieldIndex, value);
        break;

      case ComparisonOperator.Gt:
        pks = handleGt(fieldIndex, value);
        break;

      case ComparisonOperator.Gte:
        pks = handleGte(fieldIndex, value);
        break;

      case ComparisonOperator.Lt:
        pks = handleLt(fieldIndex, value);
        break;

      case ComparisonOperator.Lte:
        pks = handleLte(fieldIndex, value);
        break;

      case ComparisonOperator.Contains:
        pks = handleContains(fieldIndex, value, caseSensitive);
        break;

      case ComparisonOperator.NotContains:
        pks = handleNotContains(fieldIndex, value, caseSensitive);
        break;

      case ComparisonOperator.StartsWith:
        pks = handleStartsWith(fieldIndex, value, caseSensitive);
        break;

      case ComparisonOperator.In:
        pks = handleIn(fieldIndex, value);
        break;

      case ComparisonOperator.Between:
        pks = handleBetween(fieldIndex, value);
        break;

      case ComparisonOperator.IsNull:
        pks = handleIsNull(fieldIndex);
        break;

      case ComparisonOperator.IsNotEmpty:
        pks = handleIsNotEmpty(fieldIndex);
        break;

      case ComparisonOperator.IsEmpty:
        pks = handleIsEmpty(fieldIndex);
        break;

      case ComparisonOperator.IsNoneOf:
        pks = handleIsNoneOf(fieldIndex, value);
        break;

      default:
        throw new Error(`Operator ${operator} is not supported.`);
    }

    return pks;
  }
}

function handleEq(fieldIndex: IndexedFieldRecord, value: any): PrimaryKey[] {
  return fieldIndex[value] || [];
}

function handleGt(fieldIndex: IndexedFieldRecord, value: any): PrimaryKey[] {
  return Object.keys(fieldIndex)
    .filter((key) => +key > value)
    .flatMap((key) => fieldIndex[key]);
}

function handleGte(fieldIndex: IndexedFieldRecord, value: any): PrimaryKey[] {
  return Object.keys(fieldIndex)
    .filter((key) => +key >= value)
    .flatMap((key) => fieldIndex[key]);
}

function handleLt(fieldIndex: IndexedFieldRecord, value: any): PrimaryKey[] {
  return Object.keys(fieldIndex)
    .filter((key) => +key < value)
    .flatMap((key) => fieldIndex[key]);
}

function handleLte(fieldIndex: IndexedFieldRecord, value: any): PrimaryKey[] {
  return Object.keys(fieldIndex)
    .filter((key) => +key <= value)
    .flatMap((key) => fieldIndex[key]);
}

function handleContains(
  fieldIndex: IndexedFieldRecord,
  value: string,
  caseSensitive: boolean | null = true,
): PrimaryKey[] {
  return Object.keys(fieldIndex)
    .filter((key) =>
      caseSensitive
        ? key.includes(value)
        : key.toLowerCase().includes(value.toLowerCase()),
    )
    .flatMap((key) => fieldIndex[key]);
}

function handleNotContains(
  fieldIndex: IndexedFieldRecord,
  value: string,
  caseSensitive: boolean | null = true,
): PrimaryKey[] {
  return Object.keys(fieldIndex)
    .filter((key) =>
      caseSensitive
        ? !key.includes(value)
        : !key.toLowerCase().includes(value.toLowerCase()),
    )
    .flatMap((key) => fieldIndex[key]);
}

function handleStartsWith(
  fieldIndex: IndexedFieldRecord,
  value: string,
  caseSensitive: boolean | null = true,
): PrimaryKey[] {
  return Object.keys(fieldIndex)
    .filter((key) =>
      caseSensitive
        ? key.startsWith(value)
        : key.toLowerCase().startsWith(value.toLowerCase()),
    )
    .flatMap((key) => fieldIndex[key]);
}

function handleIn(fieldIndex: IndexedFieldRecord, values: any[]): PrimaryKey[] {
  return values.flatMap((value) => fieldIndex[value] || []);
}

function handleBetween(
  fieldIndex: IndexedFieldRecord,
  range: [any, any],
): PrimaryKey[] {
  const [min, max] = range;

  return Object.keys(fieldIndex)
    .filter((key) => +key >= min && +key <= max)
    .flatMap((key) => fieldIndex[key]);
}

function handleIsNull(fieldIndex: IndexedFieldRecord): PrimaryKey[] {
  // We use a special key '__NULL__' to store null values
  return fieldIndex['__NULL__'] || [];
}

function handleIsNotEmpty(fieldIndex: IndexedFieldRecord): PrimaryKey[] {
  // Get all keys except null, undefined, and empty string
  return Object.keys(fieldIndex)
    .filter(
      (key) =>
        key !== '__NULL__' && key !== '__UNDEFINED__' && key !== '__EMPTY__',
    )
    .flatMap((key) => fieldIndex[key]);
}

function handleIsEmpty(fieldIndex: IndexedFieldRecord): PrimaryKey[] {
  // Return all records that are either null, undefined, or an empty string
  return [
    ...(fieldIndex['__NULL__'] || []),
    ...(fieldIndex['__UNDEFINED__'] || []),
    ...(fieldIndex['__EMPTY__'] || []),
  ];
}

function handleIsNoneOf(
  fieldIndex: IndexedFieldRecord,
  values: any[],
): PrimaryKey[] {
  return Object.keys(fieldIndex)
    .filter((key) => !values.includes(key))
    .flatMap((key) => fieldIndex[key]);
}

///

function intersect(arr1: PrimaryKey[], arr2: PrimaryKey[]): PrimaryKey[] {
  return arr1.filter((id) => arr2.includes(id));
}

function union(arr1: PrimaryKey[], arr2: PrimaryKey[]): PrimaryKey[] {
  return Array.from(new Set([...arr1, ...arr2])); // Remove duplicates
}
