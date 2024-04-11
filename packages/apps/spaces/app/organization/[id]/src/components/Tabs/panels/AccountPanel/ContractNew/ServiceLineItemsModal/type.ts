import { ServiceLineItem } from '@graphql/types';

export type ServiceItem = Partial<ServiceLineItem> & {
  type: string;
  isDeleted: boolean;
};
