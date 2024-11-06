import { rdiffResult } from 'recursive-diff';

import {
  Filter as ServerFilter,
  FilterItem as ServerFilterItem,
} from '@graphql/types';

export type Operation = {
  id: number;
  ref?: string;
  entity: string;
  tenant?: string;
  entityId?: string;
  diff: rdiffResult[];
};
export type GroupOperation = {
  ref?: string;
  ids: string[];
  entity?: string;
  tenant?: string;
  entityId?: string;
  action: 'APPEND' | 'DELETE' | 'INVALIDATE';
};

export type SyncPacket = {
  ref?: string;
  version: number;
  entity_id: string;
  operation: Operation;
};

export type GroupSyncPacket = {
  ref?: string;
  ids: string[];
  action: 'APPEND' | 'DELETE' | 'INVALIDATE';
};

export type SystemSyncPacket = {
  action: 'NEW_VERSION_AVAILABLE';
};

export type LatestDiff = {
  version: number;
  entity_id: string;
  operations: Operation[];
};

export type FilterItem = ServerFilterItem & { active?: boolean };
export type Filter = Omit<ServerFilter, 'filter'> & { filter?: FilterItem };
