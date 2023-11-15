'use client';
import { Command, CommandMenu, useCommands, CommandWrapper } from 'kmenu';

import { useLastTouchpointFilter } from '@organizations/components/Columns/Filters/LastTouchpoint/LastTouchpointFilter.atom';

import 'kmenu/dist/index.css';

export const KMenu = () => {
  const [_, setTouchpointFilter] = useLastTouchpointFilter();

  const main: Command[] = [
    {
      category: 'Filter',
      commands: [
        {
          text: 'Show organizations that need attention',
          perform: () => alert('Show organizations that need attention'),
        },
        {
          text: 'Show least active organizations',
          perform: () => alert('Show least active organizations'),
        },
        {
          text: 'Show most active organizations',
          perform: () => alert('Show most active organizations'),
        },
        {
          text: 'Show recently emailed organizations',
          perform: () => {
            setTouchpointFilter((prev) => {
              return {
                ...prev,
                isActive: true,
                value: ['EMAIL'],
              };
            });
          },
        },
        {
          text: 'Show recently created organizations',
          perform: () => alert('Show recently created organizations'),
        },
        {
          text: 'Clear all filters',
          perform: () =>
            setTouchpointFilter((prev) => ({
              ...prev,
              value: [],
              isActive: false,
            })),
        },
      ],
    },
  ];

  const [mainCommands] = useCommands(main);

  return (
    <CommandWrapper>
      <CommandMenu
        commands={mainCommands}
        crumbs={['Organizations']}
        index={1}
        placeholder='What would you like to do?'
      />
    </CommandWrapper>
  );
};
