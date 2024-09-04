import type { Transport } from '@store/transport';

import { gql } from 'graphql-request';

import type {
  BankAccountCreateInput,
  BankAccountUpdateInput,
} from '@graphql/types';

class BankAccountService {
  private static instance: BankAccountService | null = null;
  private transport: Transport;

  constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport): BankAccountService {
    if (!BankAccountService.instance) {
      BankAccountService.instance = new BankAccountService(transport);
    }

    return BankAccountService.instance;
  }

  async updateBankAccount(
    payload: BANK_ACCOUNT_UPDATE_INPUT,
  ): Promise<BANK_ACCOUNT_UPDATE_RESPONSE> {
    return this.transport.graphql.request<
      BANK_ACCOUNT_UPDATE_RESPONSE,
      BANK_ACCOUNT_UPDATE_INPUT
    >(UPDATE_BANK_ACCOUNT_MUTATION, payload);
  }

  async createBankAccount(
    payload: BANK_ACCOUNT_CREATE_INPUT,
  ): Promise<BANK_ACCOUNT_CREATE_RESPONSE> {
    return this.transport.graphql.request<
      BANK_ACCOUNT_CREATE_RESPONSE,
      BANK_ACCOUNT_CREATE_INPUT
    >(CREATE_BANK_ACCOUNT_MUTATION, payload);
  }

  async deleteBankAccount(
    id: string,
  ): Promise<{ accepted: boolean; completed: boolean }> {
    return this.transport.graphql.request<{
      accepted: boolean;
      completed: boolean;
    }>(DELETE_BANK_ACCOUNT_MUTATION, { id });
  }
}

type BANK_ACCOUNT_UPDATE_RESPONSE = {
  bankAccount_Update: {
    bic: string;
    iban: string;
    currency: string;
    bankName: string;
    sortCode: string;
    accountNumber: string;
    routingNumber: string;
    bankTransferEnabled: boolean;
  };
};
type BANK_ACCOUNT_UPDATE_INPUT = {
  input: BankAccountUpdateInput;
};
export const UPDATE_BANK_ACCOUNT_MUTATION = gql`
  mutation updateBankAccount($input: BankAccountUpdateInput!) {
    bankAccount_Update(input: $input) {
      currency
      bankName
      bankTransferEnabled
      iban
      bic
      sortCode
      accountNumber
      routingNumber
    }
  }
`;
type BANK_ACCOUNT_CREATE_RESPONSE = {
  bankAccount_Create: {
    bic: string;
    iban: string;
    currency: string;
    bankName: string;
    sortCode: string;
    accountNumber: string;
    routingNumber: string;
    bankTransferEnabled: boolean;
    metadata: {
      id: string;
    };
  };
};
type BANK_ACCOUNT_CREATE_INPUT = {
  input: BankAccountCreateInput;
};

export const CREATE_BANK_ACCOUNT_MUTATION = gql`
  mutation createBankAccount($input: BankAccountCreateInput!) {
    bankAccount_Create(input: $input) {
      currency
      bankName
      bankTransferEnabled
      iban
      bic
      sortCode
      accountNumber
      routingNumber
    }
  }
`;

export const DELETE_BANK_ACCOUNT_MUTATION = gql`
  mutation deleteBankAccount($id: ID!) {
    bankAccount_Delete(id: $id) {
      accepted
      completed
    }
  }
`;

export { BankAccountService };
