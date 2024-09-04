import { Channel } from 'phoenix';
import { gql } from 'graphql-request';
import { Store } from '@store/store.ts';
import { RootStore } from '@store/root.ts';
import { Transport } from '@store/transport.ts';
import { GroupOperation } from '@store/types.ts';
import { runInAction, makeAutoObservable } from 'mobx';
import { GroupStore, makeAutoSyncableGroup } from '@store/group-store.ts';
import { BankAccountService } from '@store/BankAccounts/BankAccount.service.ts';

import { BankAccount, BankAccountCreateInput } from '@graphql/types';

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
  private service: BankAccountService;

  constructor(public root: RootStore, public transport: Transport) {
    makeAutoObservable(this);
    makeAutoSyncableGroup(this, {
      channelName: 'BankAccounts',
      getItemId: (item) => item?.metadata?.id,
      ItemStore: BankAccountStore,
    });
    this.service = BankAccountService.getInstance(transport);
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

  async invalidate() {
    this.isLoading = true;

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

  async create(
    payload: BankAccountCreateInput,
    options?: { onSuccess?: (serverId: string) => void },
  ) {
    const newBankAccount = new BankAccountStore(this.root, this.transport);
    const tempId = newBankAccount.value.metadata?.id;

    newBankAccount.value = {
      ...newBankAccount.value,
      ...newBankAccount.value,
      bankName: payload?.bankName ?? newBankAccount.value.bankName,
      currency: payload?.currency ?? newBankAccount.value.currency,
    };

    let serverId: string | undefined;

    this.value.set(tempId, newBankAccount);

    try {
      const { bankAccount_Create } = await this.service.createBankAccount({
        input: payload,
      });

      runInAction(() => {
        serverId = bankAccount_Create?.metadata.id;

        newBankAccount.setId(serverId);

        this.value.set(serverId, newBankAccount);
        this.value.delete(tempId);

        this.sync({ action: 'APPEND', ids: [serverId] });
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      serverId && options?.onSuccess?.(serverId);
      setTimeout(() => {
        if (serverId) {
          this.value.get(serverId)?.invalidate();
        }
      }, 1000);
    }
  }

  async remove(id: string) {
    try {
      await this.service.deleteBankAccount(id);
      runInAction(() => {
        this.value.delete(id);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error)?.message;
      });
    } finally {
      runInAction(() => {
        this.sync({ action: 'DELETE', ids: [id] });
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
