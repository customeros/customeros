import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { FlowContact } from '@graphql/types';
import { Table } from '@ui/presentation/Table';
import { useStore } from '@shared/hooks/useStore';

export const SubjectsTable = observer(() => {
  const store = useStore();
  const params = useParams<{ id: string }>();
  const flow = store.flows.value.get(params.id as string);
  const subjects = flow?.value?.contacts || [];

  return <Table<FlowContact> columns={[]} data={subjects} />;
});
