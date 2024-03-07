import { useQuery } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';

import { ReminderPostit } from '../../shared/ReminderPostit';

type Reminder = {
  id: string;
  date: string;
  content: string;
};

const mockData: Reminder[] = [
  {
    id: '1',
    date: '2021-10-01',
    content: 'Reminder 1',
  },
  {
    id: '2',
    date: '2021-10-02',
    content: 'Reminder 2',
  },
  {
    id: '3',
    date: '2021-10-03',
    content: 'Reminder 3',
  },
];

export const Reminders = () => {
  const { data, isPending } = useQuery<Reminder[]>({
    queryKey: ['reminders'],
    queryFn: async () => {
      return new Promise((resolve) => resolve(mockData));
    },
  });

  if (isPending) return <p>Loading...</p>;

  return (
    <VStack align='flex-start'>
      {data?.map((r) => (
        <ReminderPostit key={r.id}>
          <Flex px='4' mb='2'>
            <Text>{r.content}</Text>
          </Flex>
          <Flex align='center' px='4' w='full' justify='space-between' mb='2'>
            <Text>24 Mar • 09:09</Text>
            <Button variant='ghost' colorScheme='yellow' size='sm'>
              Dismiss
            </Button>
          </Flex>
        </ReminderPostit>
      ))}
    </VStack>
  );
};
