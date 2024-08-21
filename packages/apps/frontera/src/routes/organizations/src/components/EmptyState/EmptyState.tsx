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
      case TableIdType.MyPortfolio:
        return {
          title: 'No organizations assigned to you yet',
          description:
            'Currently, you have not been assigned to any organizations.\n' +
            '\n' +
            'Head to your list of organizations and assign yourself as an owner to one of them.',
          buttonLabel: 'Go to Organizations',
          dataTest: 'portfolio-go-to-all-orgs',
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

      case TableIdType.Nurture:
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
      case TableIdType.Churn:
        return {
          title: 'Smooth sailing',
          description:
            'Seems like your customers are loyal! No one has churned yet. Keep up the strong relationships.',
          buttonLabel: 'Go to Organizations',
          dataTest: 'churn-go-to-all-orgs',
          onClick: () => {
            navigate(`/finder?preset=${allOrgsView}`);
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

          {currentPreset !== TableIdType.Leads && (
            <Button
              variant='outline'
              onClick={options.onClick}
              data-test={options.dataTest}
              className='mt-4 min-w-min text-sm bg-white'
            >
              {options.buttonLabel}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
});
