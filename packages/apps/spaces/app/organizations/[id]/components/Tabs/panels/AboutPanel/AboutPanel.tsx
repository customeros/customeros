'use client';
import { useParams } from 'next/navigation';
import { useForm } from 'react-inverted-form';
import { useQueryClient } from '@tanstack/react-query';

import { Flex } from '@ui/layout/Flex';
import { Icons } from '@ui/media/Icon';
import { FormInput } from '@ui/form/Input';
import { FormSelect } from '@ui/form/SyncSelect';
import { VStack, HStack } from '@ui/layout/Stack';
import { FormInputGroup } from '@ui/form/InputGroup';
import { OrganizationRelationship } from '@graphql/types';
import { FormAutoresizeTextarea } from '@ui/form/Textarea';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';

import { useOrganizationQuery } from '@organization/graphql/organization.generated';
import { useAddRelationshipMutation } from '@organization/graphql/addRelationship.generated';
import { useRemoveRelationshipMutation } from '@organization/graphql/removeRelationship.generated';
import { useUpdateOrganizationMutation } from '@organization/graphql/updateOrganization.generated';
import { useSetRelationshipStageMutation } from '@organization/graphql/setRelationshipStage.generated';
import { useRemoveRelationshipStageMutation } from '@organization/graphql/removeRelationshipStage.generated';

import {
  stageOptions,
  industryOptions,
  employeesOptions,
  relationshipOptions,
  businessTypeOptions,
  lastFundingRoundOptions,
} from './util';
import {
  OrganizationAboutForm,
  OrganizationAboutFormDto,
} from './OrganizationAbout.dto';
import { FormUrlInput } from './FormUrlInput';
import { FormSocialInput } from './FormSocialInput';

const placeholders = {
  valueProposition: `Value proposition (A company's value prop is its raison d'être, its sweet spot, its jam. It's the special sauce that makes customers come back for more. It's the secret behind "Shut up and take my money!")`,
};

export const AboutPanel = () => {
  const id = useParams()?.id as string;

  const client = getGraphQLClient();
  const queryClient = useQueryClient();
  const { data } = useOrganizationQuery(client, { id });


  const invalidateQuery = () =>
    queryClient.invalidateQueries(useOrganizationQuery.getKey({ id }));

  const addRelationship = useAddRelationshipMutation(client, {
    onSuccess: invalidateQuery,
  });
  const removeRelationship = useRemoveRelationshipMutation(client, {
    onSuccess: invalidateQuery,
  });
  const setRelationshipStage = useSetRelationshipStageMutation(client, {
    onSuccess: invalidateQuery,
  });
  const removeRelationshipStage = useRemoveRelationshipStageMutation(client, {
    onSuccess: invalidateQuery,
  });
  const updateOrganization = useUpdateOrganizationMutation(client, {
    onSuccess: invalidateQuery,
  });

  const defaultValues: OrganizationAboutForm = new OrganizationAboutFormDto(
    data?.organization,
  );

  const prevRelationship =
    data?.organization?.relationshipStages?.[0]?.relationship;

  const { state } = useForm<OrganizationAboutForm>({
    formId: 'organization-about',
    defaultValues,
    stateReducer: (state, action, next) => {
      const mutateOrganization = (
        variables: Partial<OrganizationAboutForm>,
      ) => {
        updateOrganization.mutate({
          input: OrganizationAboutFormDto.toPayload({
            ...state.values,
            ...variables,
          }),
        });
      };

      if (action.type === 'FIELD_CHANGE') {
        switch (action.payload.name) {
          case 'relationship': {
            const relationship = action.payload?.value?.value;

            const add = () => {
              addRelationship.mutate({
                organizationId: id,
                relationship,
              });
            };
            const remove = (onSuccess?: () => void) => {
              removeRelationship.mutate(
                {
                  organizationId: id,
                  relationship: prevRelationship as OrganizationRelationship,
                },
                { onSuccess },
              );
            };

            if (!relationship && !prevRelationship) break;
            if (!relationship && prevRelationship) remove();
            if (!prevRelationship && relationship) add();
            if (prevRelationship && relationship)
              remove(() => {
                if (!relationship) return;
                add();
              });

            return {
              ...next,
              values: {
                ...next.values,
                stage: null,
              },
            };
          }
          case 'stage': {
            const relationship = state?.values?.relationship
              ?.value as OrganizationRelationship;
            const stage = action?.payload?.value?.value;

            if (!relationship) break;
            if (!stage) {
              removeRelationshipStage.mutate({
                organizationId: id,
                relationship,
              });
              break;
            }

            setRelationshipStage.mutate({
              organizationId: id,
              relationship,
              stage,
            });
            break;
          }
          case 'industry':
          case 'employees':
          case 'businessType':
          case 'lastFundingRound': {
            mutateOrganization({
              [action.payload.name]: action.payload?.value,
            });
            break;
          }
          default:
            return next;
        }
      }
      if (action.type === 'FIELD_BLUR') {
        switch (action.payload.name) {
          case 'name':
          case 'website':
          case 'valueProposition':
          case 'targetAudience':
          case 'lastFundingAmount': {
            mutateOrganization({
              [action.payload.name]: action.payload?.value,
            });
            break;
          }
          default:
            return next;
        }
      }

      return next;
    },
  });

  return (
    <Flex
      pt='4'
      px='6'
      pb='0'
      w='full'
      h='calc(100% - 40px)'
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
        <FormInput
          name='name'
          fontSize='lg'
          autoComplete='off'
          fontWeight='semibold'
          variant='unstyled'
          borderRadius='unset'
          placeholder='Company name'
          formId='organization-about'
        />
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

        <VStack flex='1' align='flex-start' justify='flex-start' spacing='0'>
          <HStack w='full'>
            <FormSelect
              isClearable
              name='relationship'
              formId='organization-about'
              placeholder='Relationship'
              options={relationshipOptions}
              leftElement={<Icons.HeartHand color='gray.500' mr='3' />}
            />
            {!!state.values.relationship && (
              <FormSelect
                isClearable
                name='stage'
                placeholder='Stage'
                options={stageOptions}
                formId='organization-about'
                isDisabled={!state.values.relationship}
                leftElement={<Icons.ClockRefresh color='gray.500' mr='3' />}
              />
            )}
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
            <FormInputGroup
              name='lastFundingAmount'
              formId='organization-about'
              placeholder='Last funding amount'
              leftElement={<Icons.BankNote3 color='gray.500' />}
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

          <FormSocialInput
            name='socials'
            organizationId={id}
            placeholder='Social link'
            formId='organization-about'
            leftElement={<Icons.Share7 color='gray.500' />}
          />
        </VStack>
      </Flex>
    </Flex>
  );
};
