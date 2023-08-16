'use client';
import { FC } from 'react';
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

export type RenewalLikelihoodType = RenewalLikelihoodT;

export const RenewalLikelihood: FC<{
  renewalLikelihood: RenewalLikelihoodType;
  name: string;
}> = ({ renewalLikelihood, name }) => {
  const updateModal = useDisclosure();
  const infoModal = useDisclosure();
  const { probability, comment } = renewalLikelihood;

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
            colorScheme={getFeatureIconColor(probability)}
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
                {!probability
                  ? 'Not set yet'
                  : `Set by 
                ${getUserDisplayData(renewalLikelihood?.updatedBy)}
                 ${DateTimeUtils.timeAgo(renewalLikelihood.updatedAt, {
                   addSuffix: true,
                 })}`}
              </Text>
            </Flex>

            <Heading fontSize='2xl' color={getRenewalColor(probability)}>
              {parseRenewalLabel(probability)}
            </Heading>
          </Flex>
        </CardBody>
        <CardFooter p='0' as={Flex} flexDir='column'>
          {comment && (
            <>
              <Divider mt='4' mb='2' />
              <Flex align='flex-start'>
                <Icons.File2 color='gray.400' />
                <Text color='gray.500' fontSize='xs' ml='1' noOfLines={2}>
                  {comment}
                </Text>
              </Flex>
            </>
          )}
        </CardFooter>
      </Card>

      <RenewalLikelihoodModal
        name={name}
        renewalLikelihood={renewalLikelihood}
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
          automatically based on certain criteria. It is used to prioritise
          actions and calculate renewal forecasts.
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          It is used to prioritise actions and calculate Renewal forecasts.
        </Text>
      </InfoDialog>
    </>
  );
};

function getFeatureIconColor(
  renewalLikelihood?: Maybe<RenewalLikelihoodProbability> | undefined,
) {
  switch (renewalLikelihood) {
    case 'HIGH':
      return 'success';
    case 'MEDIUM':
      return 'warning';
    case 'LOW':
      return 'error';
    case 'ZERO':
      return 'gray';
    default:
      return 'gray';
  }
}

function parseRenewalLabel(
  renewalLikelihood?: Maybe<RenewalLikelihoodProbability> | undefined,
) {
  switch (renewalLikelihood) {
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
  renewalLikelihood?: Maybe<RenewalLikelihoodProbability> | undefined,
) {
  switch (renewalLikelihood) {
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
