import { useNavigate, useSearchParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';

import { TableIdType } from '@graphql/types';
import { Button } from '@ui/form/Button/Button';
import { useStore } from '@shared/hooks/useStore';
import { EmptyTable } from '@ui/media/logos/EmptyTable';

import HalfCirclePattern from '../../../../src/assets/HalfCirclePattern';

export const EmptyState = observer(() => {
  const store = useStore();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const preset = searchParams?.get('preset');

  const currentPreset = store.tableViewDefs
    ?.toArray()
    .find((e) => e.value.id === preset)?.value?.tableId;

  const allOrgsView = store.tableViewDefs.organizationsPreset;

  const options = (() => {
    switch (currentPreset) {
      case TableIdType.Organizations:
        return {
          title: "Let's get started",
          description:
            'Get started by manually adding an organization or connecting an app in Settings',
          buttonLabel: 'Add organization',
          dataTest: 'all-orgs-add-org',
          onClick: () => {
            store.ui.commandMenu.setType('AddNewOrganization');
            store.ui.commandMenu.setOpen(true);
          },
        };
      case TableIdType.Contacts:
        return {
          title: 'No contacts created yet',
          description: 'Currently, there are no contacts created yet.',
          buttonLabel: 'Go to Organizations',
          dataTest: 'contacts-go-to-orgs',
          onClick: () => {
            navigate(`/finder?preset=${allOrgsView}`);
          },
        };
      case TableIdType.Customers:
        return {
          title: 'Who will be first?',
          description:
            'No customers here yet. You can change prospects into customers by changing their relationship status in the About section.',
          buttonLabel: 'Go to Organizations',
          dataTest: 'customers-go-to-all-orgs',
          onClick: () => {
            navigate(`/finder?preset=${allOrgsView}`);
          },
        };

      case TableIdType.Targets:
        return {
          title: 'Bullseye pending',
          description:
            'We’re sorting through your Leads in the Organizations view using your Ideal Company Profile. Once qualified, they will automatically show up here as Targets.',
          buttonLabel: 'Go to Organizations',
          dataTest: 'targets-go-to-leads',
          onClick: () => {
            navigate(`/finder?preset=${allOrgsView}`);
          },
        };
      case TableIdType.Contracts:
        return {
          title: 'No signatures yet',
          description:
            'No contracts here yet. Once you create a contract for an organization, they will show up here.',
          buttonLabel: 'Go to Organizations',
          dataTest: 'contracts-go-to-all-orgs',
          onClick: () => {
            navigate(`/finder?preset=${allOrgsView}`);
          },
        };
      case TableIdType.PastInvoices:
      case TableIdType.UpcomingInvoices:
        return {
          title: 'No paper trails yet',
          description:
            'Once you generate an invoice from a customer’s contract, they will show up here.',
          buttonLabel: 'Go to Organizations',
          dataTest: 'invoices-go-to-all-orgs',
          onClick: () => {
            navigate(`/finder?preset=${allOrgsView}`);
          },
        };
      case TableIdType.FlowSequences:
        return {
          title: 'No sequences yet',
          description:
            'Your sequences are waiting to take their first steps. Go ahead and create your first sequence.',
          buttonLabel: 'New sequence',
          dataTest: 'sequence-create-new-sequence',
          onClick: () => {
            store.ui.commandMenu.setType('CreateNewSequence');
            store.ui.commandMenu.setOpen(true);
          },
        };
      default:
        return {
          title: "We couldn't find any organizations",
          description:
            'Manually add an organization or connect an app in Settings',
          buttonLabel: 'Add organization',
          dataTest: 'go-to-all-orgs',
          onClick: () => {
            store.ui.commandMenu.setType('AddNewOrganization');
            store.ui.commandMenu.setOpen(true);
          },
        };
    }
  })();

  return (
    <div className='flex items-center justify-center h-full bg-white'>
      <div className='flex flex-col h-[500px] w-[500px]'>
        <div className='flex relative'>
          <EmptyTable className='w-[152px] h-[120px] absolute top-[25%] right-[35%]' />
          <HalfCirclePattern width={500} height={500} />
        </div>
        <div className='flex flex-col text-center items-center top-[5vh] transform translate-y-[-230px]'>
          <p className='text-gray-900 text-md font-semibold'>{options.title}</p>
          <p className='max-w-[400px] text-sm text-gray-600 my-1'>
            {options.description}
          </p>

          <Button
            variant='outline'
            onClick={options.onClick}
            data-test={options.dataTest}
            className='mt-4 min-w-min text-sm bg-white'
          >
            {options.buttonLabel}
          </Button>
        </div>
      </div>
    </div>
  );
});
