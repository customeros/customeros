'use client';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Divider } from '@ui/presentation/Divider';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { useDisclosure } from '@ui/utils';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { RenewalLikelihoodModal } from './RenewalLikelihoodModal';
import {
  Maybe,
  RenewalLikelihood as RenewalLikelihoodT,
  RenewalLikelihoodProbability,
} from '@graphql/types';
import { getUserDisplayData } from '@spaces/utils/getUserEmail';
import { DateTimeUtils } from '@spaces/utils/date';
import { getFeatureIconColor } from '@organization/components/Tabs/panels/AccountPanel/utils';

export type RenewalLikelihoodType = RenewalLikelihoodT;

interface RenewalLikelihoodProps {
  data: RenewalLikelihoodType;
  name: string;
}

export const RenewalLikelihood = ({ data, name }: RenewalLikelihoodProps) => {
  const updateModal = useDisclosure();
  const infoModal = useDisclosure();

  return (
    <>
      <Card
        p='4'
        w='full'
        size='lg'
        boxShadow='xs'
        variant='outline'
        cursor='pointer'
        onClick={updateModal.onOpen}
      >
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={getFeatureIconColor(data?.probability)}
          >
            <Icons.HeartActivity />
          </FeaturedIcon>
          <Flex ml='5' align='center' justify='space-between' w='full'>
            <Flex flexDir='column'>
              <Flex align='center'>
                <Heading
                  size='sm'
                  fontWeight='semibold'
                  color='gray.700'
                  mr={2}
                >
                  Renewal likelihood
                </Heading>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Help'
                  onClick={(e) => {
                    e.stopPropagation();
                    infoModal.onOpen();
                  }}
                  icon={<Icons.HelpCircle color='gray.400' />}
                />
              </Flex>
              <Text fontSize='xs' color='gray.500'>
                {!data?.probability
                  ? 'Not set yet'
                  : `Set by 
                ${getUserDisplayData(data?.updatedBy)}
                 ${DateTimeUtils.timeAgo(data?.updatedAt, {
                   addSuffix: true,
                 })}`}
              </Text>
            </Flex>

            <Heading fontSize='2xl' color={getRenewalColor(data?.probability)}>
              {parseRenewalLabel(data?.probability)}
            </Heading>
          </Flex>
        </CardBody>

        {data?.probability && data?.updatedBy && (
          <CardFooter p='0' as={Flex} flexDir='column'>
            <Divider my='4' />
            <Flex align='flex-start'>
              {data?.comment ? (
                <Icons.File2 color='gray.400' />
              ) : (
                <Icons.FileCross viewBox='0 0 16 16' color='gray.400' />
              )}

              <Text color='gray.500' fontSize='xs' ml='1' noOfLines={2}>
                {data?.comment || 'No reason provided'}
              </Text>
            </Flex>
          </CardFooter>
        )}
      </Card>

      <RenewalLikelihoodModal
        name={name}
        renewalLikelihood={data}
        isOpen={updateModal.isOpen}
        onClose={updateModal.onClose}
      />

      <InfoDialog
        isOpen={infoModal.isOpen}
        onClose={infoModal.onClose}
        onConfirm={infoModal.onClose}
        confirmButtonLabel='Got it'
        label='Renewal likelihood'
      >
        <Text fontSize='sm' fontWeight='normal'>
          Renewal likelihood is a rough forecast of how likely {name} is to
          renew their account. This value can be manually set by you or
          automatically based on certain criteria.
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          It is used to prioritise actions and calculate Renewal forecasts.
        </Text>
      </InfoDialog>
    </>
  );
};

function parseRenewalLabel(
  data?: Maybe<RenewalLikelihoodProbability> | undefined,
) {
  switch (data) {
    case 'HIGH':
      return 'High';
    case 'MEDIUM':
      return 'Medium';
    case 'LOW':
      return 'Low';
    case 'ZERO':
      return 'Zero';
    default:
      return 'Not set';
  }
}

function getRenewalColor(
  data?: Maybe<RenewalLikelihoodProbability> | undefined,
) {
  switch (data) {
    case 'HIGH':
      return 'success.500';
    case 'MEDIUM':
      return 'warning.500';
    case 'LOW':
      return 'error.500';
    case 'ZERO':
      return 'gray.700';
    default:
      return 'gray.400';
  }
}
