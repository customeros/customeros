import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Store } from '@store/store.ts';
import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { GroupOperation } from '@store/types.ts';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store.ts';

import { BankAccount } from '@graphql/types';

import mock from './mock.json';
import { BankAccountStore } from './BankAccount.store.ts';

export class BankAccountsStore implements GroupStore<BankAccount> {
  version = 0;
  isLoading = false;
  history: GroupOperation[] = [];
  error: string | null = null;
  channel?: Channel | undefined;
  isBootstrapped: boolean = false;
  value: Map<string, Store<BankAccount>> = new Map();
  sync = makeAutoSyncableGroup.sync;
  subscribe = makeAutoSyncableGroup.subscribe;
  load = makeAutoSyncableGroup.load<BankAccount>();
  totalElements = 0;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'BankAccounts',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: BankAccountStore,
    });
  }

  toArray() {
    return Array.from(this.value.values());
  }

  async bootstrap() {
    if (this.root.demoMode) {
      this.load(mock.data.bankAccounts as BankAccount[]);
      this.isBootstrapped = true;

      return;
    }

    if (this.isBootstrapped || this.isLoading) return;

    try {
      const { bankAccounts } =
        await this.transport.graphql.request<BANK_ACCOUNTS_RESPONSE>(
          BANK_ACCOUNTS_QUERY,
        );

      runInAction(() => {
        this.load(bankAccounts);
      });
      runInAction(() => {
        this.isBootstrapped = true;
        this.totalElements = bankAccounts.length;
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }
}

type BANK_ACCOUNTS_RESPONSE = {
  bankAccounts: BankAccount[];
};
const BANK_ACCOUNTS_QUERY = gql`
  query getBankAccounts {
    bankAccounts {
      metadata {
        id
      }
      currency
      bic
      iban
      bankName
      sortCode
      allowInternational
      bankTransferEnabled
      otherDetails
      accountNumber
      routingNumber
    }
  }
`;
