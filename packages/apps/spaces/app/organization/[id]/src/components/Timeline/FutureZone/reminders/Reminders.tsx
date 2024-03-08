import { useForm } from 'react-inverted-form';

import { useQuery } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Button } from '@ui/form/Button';
import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useGlobalCacheQuery } from '@shared/graphql/global_Cache.generated';

import { ReminderPostit } from '../../shared/ReminderPostit';

type Reminder = {
  id: string;
  date: string;
  owner: string;
  content: string;
};

const mockData: Reminder[] = [
  {
    id: '1',
    date: '2021-10-01',
    content: 'Reminder 1',
    owner: 'customerostest',
  },
  {
    id: '2',
    date: '2021-10-02',
    content: 'Reminder 2',
    owner: 'Gigel',
  },
  {
    id: '3',
    date: '2021-10-03',
    content: 'Reminder 3',
    owner: 'Frone',
  },
];

export const Reminders = () => {
  const client = getGraphQLClient();
  const { data, isPending } = useQuery<Reminder[]>({
    queryKey: ['reminders'],
    queryFn: async () => {
      return new Promise((resolve) => resolve(mockData));
    },
  });

  const { state, setDefaultValues } = useForm<Reminder>({
    formId: 'reminder-edit-form',
    defaultValues: { id: '', date: '', owner: '', content: '' },
  });

  const { data: globalCacheData } = useGlobalCacheQuery(client);

  const user = globalCacheData?.global_Cache?.user;
  const currentOwner = [user?.firstName, user?.lastName]
    .filter(Boolean)
    .join(' ');

  if (isPending) return <p>Loading...</p>;

  return (
    <VStack align='flex-start'>
      {data?.map((r) => (
        <ReminderPostit key={r.id}>
          {state.values.id !== r.id ? (
            <Flex
              px='4'
              mb='2'
              whiteSpace='pre'
              fontFamily='sticky'
              onClick={() => setDefaultValues(r)}
            >
              <Text fontSize='sm'>
                {r.owner !== currentOwner && (
                  <Text as='span'>{`for ${currentOwner}: `}</Text>
                )}
                {r.content}
              </Text>
            </Flex>
          ) : (
            <FormAutoresizeTextarea
              autoFocus
              px='4'
              fontFamily='sticky'
              fontSize='sm'
              name='content'
              formId='reminder-edit-form'
              onBlur={() => {
                setDefaultValues({ id: '', date: '', owner: '', content: '' });
              }}
              placeholder='Type your reminder here'
              borderBottom='unset'
              _hover={{
                borderBottom: 'unset',
              }}
              _focus={{
                borderBottom: 'unset',
              }}
            />
          )}
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
