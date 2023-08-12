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

import {
  Likelihood,
  RenewalLikelihoodModal,
  Value as RenewalLikelihoodValue,
} from './RenewalLikelihoodModal';
import { useState } from 'react';

export const RenewalLikelihood = () => {
  const update = useDisclosure();
  const info = useDisclosure();
  const [renewalLikelihood, setRenewalLikelihood] =
    useState<RenewalLikelihoodValue>({ reason: '', likelihood: 'NOT_SET' });
  const { likelihood, reason } = renewalLikelihood;

  return (
    <>
      <Card
        p='4'
        w='full'
        size='lg'
        boxShadow='xs'
        variant='outline'
        cursor='pointer'
        onClick={update.onOpen}
      >
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon size='md' colorScheme={getFeatureIconColor(likelihood)}>
            <Icons.Building7 />
          </FeaturedIcon>
          <Flex ml='5' align='center' justify='space-between' w='full'>
            <Flex flexDir='column'>
              <Flex align='center'>
                <Heading size='sm' fontWeight='semibold' color='gray.700'>
                  Renewal likelihood
                </Heading>
                <IconButton
                  size='xs'
                  variant='ghost'
                  aria-label='Help'
                  onClick={(e) => {
                    e.stopPropagation();
                    info.onOpen();
                  }}
                  icon={<Icons.HelpCircle color='gray.400' />}
                />
              </Flex>
              <Text fontSize='xs' color='gray.500'>
                {!likelihood ? 'Not set yet' : 'Set by Unknown just now'}
              </Text>
            </Flex>

            <Heading fontSize='2xl' color={getRenewalColor(likelihood)}>
              {parseRenewalLabel(likelihood)}
            </Heading>
          </Flex>
        </CardBody>
        <CardFooter p='0' as={Flex} flexDir='column'>
          {reason && (
            <>
              <Divider mt='4' mb='2' />
              <Flex align='flex-start'>
                <Icons.File2 color='gray.400' />
                <Text color='gray.500' fontSize='xs' ml='1'>
                  {reason}
                </Text>
              </Flex>
            </>
          )}
        </CardFooter>
      </Card>

      <RenewalLikelihoodModal
        value={renewalLikelihood}
        onChange={setRenewalLikelihood}
        isOpen={update.isOpen}
        onClose={update.onClose}
      />

      <InfoDialog
        isOpen={info.isOpen}
        onClose={info.onClose}
        onConfirm={info.onClose}
        confirmButtonLabel='Got it'
        label='Renewal likelihood'
      >
        <Text fontSize='sm' fontWeight='normal'>
          Renewal likelihood is a rough forecast of how likely Acme Corp is to
          renew their account. This renewalLikelihood can be manually set by you
          or automatically based on certain criteria.
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          It is used to prioritise actions and calculate Renewal forecasts.
        </Text>
      </InfoDialog>
    </>
  );
};

function getFeatureIconColor(renewalLikelihood: Likelihood) {
  switch (renewalLikelihood) {
    case 'NOT_SET':
      return 'gray';
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

function parseRenewalLabel(renewalLikelihood: Likelihood) {
  switch (renewalLikelihood) {
    case 'NOT_SET':
      return 'Not set';
    case 'HIGH':
      return 'High';
    case 'MEDIUM':
      return 'Medium';
    case 'LOW':
      return 'Low';
    case 'ZERO':
      return 'Zero';
    default:
      'Not set';
  }
}

function getRenewalColor(renewalLikelihood: Likelihood) {
  switch (renewalLikelihood) {
    case 'NOT_SET':
      return 'gray.400';
    case 'HIGH':
      return 'success.500';
    case 'MEDIUM':
      return 'warning.500';
    case 'LOW':
      return 'error.500';
    case 'ZERO':
      return 'gray.400';
    default:
      return 'gray.400';
  }
}
