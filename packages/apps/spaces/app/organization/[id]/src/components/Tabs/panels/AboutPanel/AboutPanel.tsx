'use client';
import { useRef } from 'react';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';

import { useFeatureIsOn } from '@growthbook/growthbook-react';
import {
  useDebounce,
  useDidMount,
  useWillUnmount,
  useDeepCompareEffect,
} from 'rooks';

import { Box } from '@ui/layout/Box';
import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { Tag } from '@ui/presentation/Tag';
import { Text } from '@ui/typography/Text';
import { FormInput } from '@ui/form/Input';
import { Tooltip } from '@ui/overlay/Tooltip';
import { Organization } from '@graphql/types';
import { FormSelect } from '@ui/form/SyncSelect';
import { VStack, HStack } from '@ui/layout/Stack';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { FormNumberInputGroup } from '@ui/form/InputGroup';
import { CurrencyDollar } from '@ui/media/icons/CurrencyDollar';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { useOrganizationQuery } from '@organization/src/graphql/organization.generated';
import { Branches } from '@organization/src/components/Tabs/panels/AboutPanel/branches/Branches';
import { OwnerInput } from '@organization/src/components/Tabs/panels/AboutPanel/owner/OwnerInput';
import { ParentOrgInput } from '@organization/src/components/Tabs/panels/AboutPanel/branches/ParentOrgInput';

import { FormUrlInput } from './FormUrlInput';
import { FormSocialInput } from '../../shared/FormSocialInput';
import { useAboutPanelMethods } from './hooks/useAboutPanelMethods';
import {
  OrganizationAboutForm,
  OrganizationAboutFormDto,
} from './OrganizationAbout.dto';
import {
  industryOptions,
  employeesOptions,
  relationshipOptions,
  businessTypeOptions,
  lastFundingRoundOptions,
} from './util';

const placeholders = {
  valueProposition: `Value proposition (A company's value prop is its raison d'être, its sweet spot, its jam. It's the special sauce that makes customers come back for more. It's the secret behind "Shut up and take my money!")`,
};

export const AboutPanel = () => {
  const client = getGraphQLClient();
  const id = useParams()?.id as string;
  const [_, copyToClipboard] = useCopyToClipboard();
  const { data } = useOrganizationQuery(client, { id });
  const showParentRelationshipSelector = useFeatureIsOn(
    'show-parent-relationship-selector',
  );
  const nameRef = useRef<HTMLInputElement | null>(null);
  const parentRelationshipReadOnly = useFeatureIsOn(
    'parent-relationship-selector-read-only',
  );
  const { updateOrganization, addSocial, invalidateQuery } =
    useAboutPanelMethods({ id });

  const defaultValues: OrganizationAboutForm = new OrganizationAboutFormDto(
    data?.organization,
  );

  const mutateOrganization = (
    previousValues: OrganizationAboutForm,
    variables?: Partial<OrganizationAboutForm>,
  ) => {
    updateOrganization.mutate({
      input: OrganizationAboutFormDto.toPayload({
        ...previousValues,
        ...(variables ?? {}),
      }),
    });
  };

  const debouncedMutateOrganization = useDebounce(mutateOrganization, 500);

  const { setDefaultValues } = useForm<OrganizationAboutForm>({
    formId: 'organization-about',
    defaultValues,
    stateReducer: (state, action, next) => {
      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'isCustomer':
          case 'industry':
          case 'employees':
          case 'businessType':
          case 'lastFundingRound': {
            mutateOrganization(state.values, {
              [action.payload.name]: action.payload?.value,
            });
            break;
          }
          case 'name':
          case 'website':
          case 'valueProposition':
          case 'targetAudience':
          case 'lastFundingAmount': {
            const trimmedValue = (action.payload?.value || '')?.trim();
            if (
              //@ts-expect-error fixme
              trimmedValue === defaultValues?.[action.payload.name]
            ) {
              return next;
            }
            debouncedMutateOrganization(state.values, {
              [action.payload.name]: trimmedValue,
            });
            break;
          }
          default:
            return next;
        }
      }
      if (action.type === 'FIELD_BLUR') {
        if (action.payload.name === 'name') {
          const trimmedValue = (action.payload?.value || '')?.trim();
          if (!trimmedValue?.length) {
            mutateOrganization(state.values, {
              name: 'Unnamed',
            });

            return {
              ...next,
              values: {
                ...next.values,
                name: 'Unnamed',
              },
            };
          }
        }
      }

      return next;
    },
  });

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);

  useDidMount(() => {
    if (nameRef.current?.value === 'Unnamed') {
      nameRef.current?.focus();
      nameRef.current?.setSelectionRange(0, 7);
    }
  });

  useWillUnmount(() => {
    debouncedMutateOrganization.flush();
  });

  const handleAddSocial = ({
    newValue,
    onSuccess,
  }: {
    newValue: string;
    onSuccess: ({ id, url }: { id: string; url: string }) => void;
  }) => {
    addSocial.mutate(
      { organizationId: id, input: { url: newValue } },
      {
        onSuccess: ({ organization_AddSocial: { id, url } }) => {
          onSuccess({ id, url });
        },
      },
    );
  };

  return (
    <Flex
      pt='4'
      px='6'
      w='full'
      h='full'
      overflowY='auto'
      flex='1'
      background='gray.25'
      borderRadius='2xl'
    >
      <Flex
        h='full'
        flexDir='column'
        overflowY='auto'
        overflow='visible'
        w='full'
      >
        <Flex align='center'>
          <FormInput
            name='name'
            fontSize='lg'
            autoComplete='off'
            fontWeight='semibold'
            variant='unstyled'
            borderRadius='unset'
            placeholder='Company name'
            formId='organization-about'
            ref={nameRef}
          />
          {data?.organization?.referenceId && (
            <Box h='full' ml='4'>
              <Tooltip label={'Copy ID'}>
                <Tag
                  colorScheme='gray'
                  variant='outline'
                  color='gray.500'
                  borderRadius='full'
                  boxShadow='unset'
                  border='1px solid'
                  cursor='pointer'
                  borderColor='gray.300'
                  onClick={() => {
                    copyToClipboard(
                      data?.organization?.referenceId ?? '',
                      'Reference ID copied ',
                    );
                  }}
                >
                  <Text>{data?.organization?.referenceId}</Text>
                </Tag>
              </Tooltip>
            </Box>
          )}
        </Flex>
        <FormUrlInput
          name='website'
          autoComplete='off'
          placeholder='www.'
          variant='unstyled'
          borderRadius='unset'
          formId='organization-about'
        />

        <FormAutoresizeTextarea
          mb='6'
          spellCheck={false}
          name='valueProposition'
          formId='organization-about'
          placeholder={placeholders.valueProposition}
        />

        {!data?.organization?.subsidiaries?.length &&
          showParentRelationshipSelector && (
            <ParentOrgInput
              id={id}
              isReadOnly={parentRelationshipReadOnly}
              parentOrg={
                data?.organization?.subsidiaryOf?.[0]?.organization?.id
                  ? {
                      label:
                        data?.organization?.subsidiaryOf?.[0]?.organization
                          ?.name,
                      value:
                        data?.organization?.subsidiaryOf?.[0]?.organization?.id,
                    }
                  : null
              }
            />
          )}

        <VStack
          flex='1'
          align='flex-start'
          justify='flex-start'
          spacing='0'
          gap={0}
        >
          <HStack w='full'>
            <FormSelect
              isClearable
              name='isCustomer'
              formId='organization-about'
              placeholder='Relationship'
              options={relationshipOptions}
              leftElement={<Icons.HeartHand color='gray.500' mr='3' />}
            />
          </HStack>

          <FormSelect
            name='industry'
            isClearable
            placeholder='Industry'
            options={industryOptions}
            formId='organization-about'
            leftElement={<Icons.Building7 color='gray.500' mr='3' />}
          />

          <FormAutoresizeTextarea
            pl='30px'
            variant='flushed'
            name='targetAudience'
            formId='organization-about'
            placeholder='Target Audience'
            leftElement={<Icons.Target5 color='gray.500' />}
          />

          <FormSelect
            isClearable
            name='businessType'
            formId='organization-about'
            placeholder='Business Type'
            options={businessTypeOptions}
            leftElement={<Icons.DataFlow3 color='gray.500' mr='3' />}
          />

          <HStack w='full'>
            <FormSelect
              isClearable
              name='lastFundingRound'
              formId='organization-about'
              placeholder='Last funding round'
              options={lastFundingRoundOptions}
              leftElement={
                <Icons.HorizontalBarChart3 color='gray.500' mr='3' />
              }
            />
            <FormNumberInputGroup
              name='lastFundingAmount'
              formId='organization-about'
              placeholder='Last funding amount'
              min={0}
              leftElement={
                <Box color='gray.500' ml={1}>
                  <CurrencyDollar height={16} />
                </Box>
              }
            />
          </HStack>

          <FormSelect
            isClearable
            name='employees'
            options={employeesOptions}
            formId='organization-about'
            placeholder='Number of employees'
            leftElement={<Icons.Users2 color='gray.500' mr='3' />}
          />

          <OwnerInput
            id={id}
            owner={data?.organization?.owner}
            invalidateQuery={invalidateQuery}
          />

          <FormSocialInput
            bg='gray.25'
            name='socials'
            formId='organization-about'
            organizationId={id}
            addSocial={handleAddSocial}
            invalidateQuery={invalidateQuery}
            defaultValues={defaultValues.socials}
            placeholder='Social link'
            leftElement={<Icons.Share7 color='gray.500' />}
          />

          {!!data?.organization?.subsidiaries?.length &&
            showParentRelationshipSelector && (
              <Branches
                id={id}
                isReadOnly={parentRelationshipReadOnly}
                branches={
                  (data?.organization
                    ?.subsidiaries as Organization['subsidiaries']) ?? []
                }
              />
            )}
        </VStack>

        {data?.organization?.customerOsId && (
          <Tooltip label='Copy ID'>
            <Text
              py='3'
              w='fit-content'
              color='gray.400'
              cursor='pointer'
              onClick={() =>
                copyToClipboard(
                  data?.organization?.customerOsId ?? '',
                  'CustomerOS ID copied',
                )
              }
            >
              CustomerOS ID: {data?.organization?.customerOsId}
            </Text>
          </Tooltip>
        )}
      </Flex>
    </Flex>
  );
};
