'use client';
// TODO uncomment when forecast
// import { useParams } from 'next/navigation';
import { useEffect } from 'react';

// import { Dot } from '@ui/media/Dot';
// import { Box } from '@ui/layout/Box';
// import { Flex } from '@ui/layout/Flex';
// import { Text } from '@ui/typography/Text';
// import { Select } from '@ui/form/SyncSelect';
import { FeaturedIcon } from '@ui/media/Icon';
// import { User02 } from '@ui/media/icons/User02';
import { Heading } from '@ui/typography/Heading';
// import { Button, ButtonGroup } from '@ui/form/Button';
// import { AutoresizeTextarea } from '@ui/form/Textarea';
// import { FormCurrencyInput } from '@ui/form/CurrencyInput';
// import { RenewalLikelihoodProbability } from '@graphql/types';
// import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
// import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
// import { useGetUsersQuery } from '@organizations/graphql/getUsers.generated';
import {
  Modal,
  // ModalBody,
  // ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';

interface RenewalDetailsProps {
  // data: any;
  isOpen: boolean;
  onClose: () => void;
}

export const RenewalDetailsModal = ({
  // data,
  isOpen,
  onClose,
}: RenewalDetailsProps) => {
  // const client = getGraphQLClient();
  //
  // const formId = 'renewalDetailsForm';
  // const timeoutRef = useRef<NodeJS.Timeout | null>(null);
  // // const id = useParams()?.id as string;
  // const [probability, setLikelihood] = useState<
  //   RenewalLikelihoodProbability | undefined | null
  // >(data?.probability);
  // const [reason, setReason] = useState<string>(data?.comment || '');
  // const { data: usersData } = useGetUsersQuery(client, {
  //   pagination: {
  //     limit: 100,
  //     page: 1,
  //   },
  // });
  //
  // const options = useMemo(() => {
  //   return usersData?.users?.content
  //     ?.filter((e) => Boolean(e.firstName) || Boolean(e.lastName))
  //     ?.map((o) => ({
  //       value: o.id,
  //       label: `${o.firstName} ${o.lastName}`.trim(),
  //     }));
  // }, [data]);

  // const { data: session } = useSession();
  //
  // const client = getGraphQLClient();
  // const queryClient = useQueryClient();
  //
  // const getOrganizationQueryKey = useOrganizationAccountDetailsQuery.getKey({
  //   id,
  // });

  // const updateData = useUpdateDataMutation(client, {
  //   onSuccess: () => {
  //     queryClient.setQueryData<OrganizationAccountDetailsQuery>(
  //       getOrganizationQueryKey,
  //       (oldData) => {
  //         if (!oldData || !oldData?.organization) return;
  //
  //         return {
  //           organization: {
  //             ...(oldData?.organization ?? {}),
  //             accountDetails: {
  //               ...(oldData?.organization?.accountDetails ?? {}),
  //               data: {
  //                 comment: reason,
  //                 previousProbability: data?.probability,
  //                 probability: probability,
  //                 updatedAt: new Date(),
  //                 updatedBy: [session?.user] as unknown as User,
  //               },
  //             },
  //           },
  //         };
  //       },
  //     );
  //
  //     queryClient.setQueryData(
  //       useInfiniteGetTimelineQuery.getKey({
  //         organizationId: id,
  //         from: NEW_DATE,
  //         size: 50,
  //       }),
  //       (oldData) => {
  //         const newEvent = {
  //           __typename: 'Action',
  //           id: `timeline-event-action-new-id-${new Date()}`,
  //           actionType: 'RENEWAL_LIKELIHOOD_UPDATED',
  //           appSource: 'customer-os-api',
  //           createdAt: new Date(),
  //           metadata: JSON.stringify({
  //             likelihood: probability,
  //             reason,
  //           }),
  //           actionCreatedBy: null,
  //           content: `Renewal likelihood set to ${probability} by ${session?.user?.name}`,
  //         };
  //
  //         // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
  //         if (!oldData || !oldData.pages?.length) {
  //           return {
  //             pages: [
  //               {
  //                 organization: {
  //                   id,
  //                   timelineEventsTotalCount: 1,
  //                   timelineEvents: [newEvent],
  //                 },
  //               },
  //             ],
  //           };
  //         }
  //
  //         // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
  //         const firstPage = oldData.pages[0] ?? {};
  //         // @ts-expect-error TODO: queryClient.setQueryClient should be typed in order to fix this line
  //         const pages = oldData.pages?.slice(1);
  //
  //         const firstPageWithEvent = {
  //           ...firstPage,
  //           organization: {
  //             ...firstPage?.organization,
  //             timelineEvents: [
  //               ...(firstPage?.organization?.timelineEvents ?? []),
  //               newEvent,
  //             ],
  //             timelineEventsTotalCount:
  //               (firstPage?.organization?.timelineEventsTotalCount ?? 0) + 1,
  //           },
  //         };
  //
  //         return {
  //           ...oldData,
  //           pages: [firstPageWithEvent, ...pages],
  //         };
  //       },
  //     );
  //   },
  //   onSettled: () => {
  //     if (timeoutRef.current) {
  //       clearTimeout(timeoutRef.current);
  //     }
  //     timeoutRef.current = setTimeout(() => {
  //       queryClient.invalidateQueries(getOrganizationQueryKey);
  //     }, 1000);
  //   },
  // });

  // const handleSet = () => {
  // updateData.mutate({
  //   input: { id, probability: probability, comment: reason },
  // });
  onClose();
  // };

  useEffect(() => {
    return () => {
      // if (timeoutRef.current) {
      //   clearTimeout(timeoutRef.current);
      // }
    };
  }, []);

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent
        borderRadius='2xl'
        backgroundImage='/backgrounds/organization/circular-bg-pattern.png'
        backgroundRepeat='no-repeat'
        sx={{
          backgroundPositionX: '1px',
          backgroundPositionY: '-7px',
        }}
      >
        <ModalCloseButton />
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <ClockFastForward />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Renewal details
          </Heading>
        </ModalHeader>
        {/*<ModalBody as={Flex} flexDir='column' pb='0'>*/}
        {/*  <Select*/}
        {/*    isClearable*/}
        {/*    // value={value}*/}
        {/*    isLoading={false}*/}
        {/*    placeholder='Owner'*/}
        {/*    backspaceRemovesValue*/}
        {/*    onChange={() => null}*/}
        {/*    options={options}*/}
        {/*    leftElement={<User02 color='gray.500' mr={3} />}*/}
        {/*  />*/}

        {/*  <ButtonGroup w='full' isAttached>*/}
        {/*    <Button*/}
        {/*      w='full'*/}
        {/*      variant='outline'*/}
        {/*      leftIcon={<Dot colorScheme='success' />}*/}
        {/*      onClick={() => setLikelihood(RenewalLikelihoodProbability.High)}*/}
        {/*      bg={probability === 'HIGH' ? 'gray.100' : 'white'}*/}
        {/*    >*/}
        {/*      High*/}
        {/*    </Button>*/}
        {/*    <Button*/}
        {/*      w='full'*/}
        {/*      variant='outline'*/}
        {/*      leftIcon={<Dot colorScheme='warning' />}*/}
        {/*      onClick={() => setLikelihood(RenewalLikelihoodProbability.Medium)}*/}
        {/*      bg={probability === 'MEDIUM' ? 'gray.100' : 'white'}*/}
        {/*    >*/}
        {/*      Medium*/}
        {/*    </Button>*/}
        {/*    <Button*/}
        {/*      w='full'*/}
        {/*      variant='outline'*/}
        {/*      leftIcon={<Dot colorScheme='error' />}*/}
        {/*      onClick={() => setLikelihood(RenewalLikelihoodProbability.Low)}*/}
        {/*      bg={probability === 'LOW' ? 'gray.100' : 'white'}*/}
        {/*    >*/}
        {/*      Low*/}
        {/*    </Button>*/}
        {/*    <Button*/}
        {/*      variant='outline'*/}
        {/*      w='full'*/}
        {/*      leftIcon={<Dot />}*/}
        {/*      onClick={() => setLikelihood(RenewalLikelihoodProbability.Zero)}*/}
        {/*      bg={probability === 'ZERO' ? 'gray.100' : 'white'}*/}
        {/*    >*/}
        {/*      Zero*/}
        {/*    </Button>*/}
        {/*  </ButtonGroup>*/}
        {/*  <Text>Last updated by </Text>*/}

        {/*  <FormCurrencyInput*/}
        {/*    formId={formId}*/}
        {/*    name='arrForecast'*/}
        {/*    w='full'*/}
        {/*    placeholder='Amount'*/}
        {/*    label='ARR forecast'*/}
        {/*    min={0}*/}
        {/*    leftElement={*/}
        {/*      <Box color='gray.500'>*/}
        {/*        <CurrencyDollar height='16px' />*/}
        {/*      </Box>*/}
        {/*    }*/}
        {/*  />*/}

        {/*  {!!probability && (*/}
        {/*    <>*/}
        {/*      <Text as='label' htmlFor='reason' mt='5' fontSize='sm'>*/}
        {/*        <b>Reason for change</b> (optional)*/}
        {/*      </Text>*/}
        {/*      <AutoresizeTextarea*/}
        {/*        pt='0'*/}
        {/*        id='reason'*/}
        {/*        value={reason}*/}
        {/*        spellCheck='false'*/}
        {/*        onChange={(e) => setReason(e.target.value)}*/}
        {/*        placeholder={`What is the reason for ${*/}
        {/*          !data?.probability ? 'setting' : 'updating'*/}
        {/*        } these details`}*/}
        {/*      />*/}
        {/*    </>*/}
        {/*  )}*/}
        {/*</ModalBody>*/}
        {/*<ModalFooter p='6'>*/}
        {/*  <Button variant='outline' w='full' onClick={onClose}>*/}
        {/*    Cancel*/}
        {/*  </Button>*/}
        {/*  <Button*/}
        {/*    ml='3'*/}
        {/*    w='full'*/}
        {/*    variant='outline'*/}
        {/*    colorScheme='primary'*/}
        {/*    onClick={handleSet}*/}
        {/*  >*/}
        {/*    Update*/}
        {/*  </Button>*/}
        {/*</ModalFooter>*/}
      </ModalContent>
    </Modal>
  );
};
