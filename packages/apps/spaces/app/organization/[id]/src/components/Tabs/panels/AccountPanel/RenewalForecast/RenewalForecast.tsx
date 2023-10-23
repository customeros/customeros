'use client';
import { Flex } from '@ui/layout/Flex';
import { Heading } from '@ui/typography/Heading';
import { Text } from '@ui/typography/Text';
import { IconButton } from '@ui/form/IconButton';
import { Divider } from '@ui/presentation/Divider';
import { Icons, FeaturedIcon } from '@ui/media/Icon';
import { Card, CardBody, CardFooter } from '@ui/presentation/Card';
import { InfoDialog } from '@ui/overlay/AlertDialog/InfoDialog';
import { RenewalForecastModal } from './RenewalForecastModal';
import {
  RenewalForecast as RenewalForecastT,
  RenewalLikelihoodProbability,
} from '@graphql/types';
import { getUserDisplayData } from '@spaces/utils/getUserEmail';
import { DateTimeUtils } from '@spaces/utils/date';
import { formatCurrency } from '@spaces/utils/getFormattedCurrencyNumber';
import { getFeatureIconColor } from '@organization/src/components/Tabs/panels/AccountPanel/utils';
import { UseDisclosureReturn } from '@chakra-ui/hooks/dist/use-disclosure';

export type RenewalForecastType = RenewalForecastT & { amount?: string | null };

interface RenewalForecastProps {
  renewalForecast: RenewalForecastType;
  renewalProbability?: RenewalLikelihoodProbability | null;
  name: string;
  infoModal: UseDisclosureReturn;
  updateModal: UseDisclosureReturn;
}

export const RenewalForecast = ({
  renewalForecast,
  renewalProbability,
  name,
  infoModal,
  updateModal,
}: RenewalForecastProps) => {
  const getForecastMetaInfo = () => {
    if (!renewalForecast?.amount) {
      return 'Not calculated yet';
    }

    if (!renewalForecast?.updatedBy) {
      return 'Calculated from billing amount';
    }

    return `Set by ${getUserDisplayData(
      renewalForecast?.updatedBy,
    )} ${DateTimeUtils.timeAgo(renewalForecast?.updatedAt, {
      addSuffix: true,
    })}`;
  };

  const isAmountSet =
    renewalForecast?.amount !== null && renewalForecast?.amount !== undefined;

  return (
    <>
      <Card
        p='4'
        w='full'
        size='lg'
        variant='outline'
        cursor='pointer'
        boxShadow='xs'
        _hover={{
          boxShadow: 'md',
        }}
        transition='all 0.2s ease-out'
        onClick={updateModal.onOpen}
      >
        <CardBody as={Flex} p='0' align='center'>
          <FeaturedIcon
            size='md'
            minW='10'
            colorScheme={
              renewalForecast?.amount && !renewalForecast?.updatedBy
                ? getFeatureIconColor(renewalProbability)
                : 'gray'
            }
          >
            <Icons.Calculator />
          </FeaturedIcon>
          <Flex
            ml='5'
            w='full'
            align='center'
            columnGap={4}
            justify='space-between'
          >
            <Flex flexDir='column'>
              <Flex align='center'>
                <Heading
                  size='sm'
                  whiteSpace='nowrap'
                  fontWeight='semibold'
                  color='gray.700'
                  mr={2}
                >
                  ARR forecast
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
                {getForecastMetaInfo()}
              </Text>
            </Flex>

            <Heading
              fontSize='2xl'
              color={isAmountSet ? 'gray.700' : 'gray.400'}
            >
              {isAmountSet
                ? formatCurrency(renewalForecast?.amount ?? 0)
                : 'Unknown'}
            </Heading>
          </Flex>
        </CardBody>
        {!!renewalForecast?.amount && renewalForecast?.updatedBy && (
          <CardFooter p='0' as={Flex} flexDir='column'>
            <Divider mt='4' mb='2' />
            <Flex align='flex-start'>
              {renewalForecast?.comment ? (
                <Icons.File2 color='gray.400' />
              ) : (
                <Icons.FileCross viewBox='0 0 16 16' color='gray.400' />
              )}

              <Text color='gray.500' fontSize='xs' ml='1' noOfLines={2}>
                {renewalForecast?.comment || 'No reason provided'}
              </Text>
            </Flex>
          </CardFooter>
        )}
      </Card>

      <RenewalForecastModal
        renewalForecast={{
          amount: renewalForecast?.amount,
          comment: renewalForecast?.comment,
        }}
        renewalProbability={renewalProbability}
        name={name}
        isOpen={updateModal.isOpen}
        onClose={updateModal.onClose}
      />

      <InfoDialog
        isOpen={infoModal.isOpen}
        onClose={infoModal.onClose}
        onConfirm={infoModal.onClose}
        confirmButtonLabel='Got it'
        label='ARR forecast'
      >
        <Text fontSize='sm' fontWeight='normal'>
          The ARR forecast gives you a way to roughly project revenue per
          customer and across your entire portfolio.
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          {`It's calculated by discounting the renewal potential (billing amount * billings per cycle) based on the renewal likelihood—Medium, Low, or Zero.`}
        </Text>
        <br />
        <Text fontSize='sm' fontWeight='normal'>
          You can override this forecast at any time.
        </Text>
      </InfoDialog>
    </>
  );
};
