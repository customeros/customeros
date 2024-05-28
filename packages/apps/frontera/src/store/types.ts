import { rdiffResult } from 'recursive-diff';

export type Operation = { id: number; diff: rdiffResult[] };
export type GroupOperation = {
  ids: string[];
  action: 'APPEND' | 'DELETE' | 'SYNC';
};

export type SyncPacket = {
  version: number;
  entity_id: string;
  operation: Operation;
};

export type GroupSyncPacket = {
  ids: string[];
  action: 'APPEND' | 'DELETE' | 'SYNC';
};

export type LatestDiff = {
  version: number;
  entity_id: string;
  operations: Operation[];
};
