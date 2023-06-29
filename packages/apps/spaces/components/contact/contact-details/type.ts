import { GetContactQuery } from '@spaces/graphql';

export interface ContactDetailsProps {
  id: string;
  data: GetContactQuery['contact'];
  loading: boolean;
}
