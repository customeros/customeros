import { useParams } from 'react-router-dom';
import { useForm } from 'react-inverted-form';

import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { useDebounce, useWillUnmount, useDeepCompareEffect } from 'rooks';

import { Organization } from '@graphql/types';
import { FormUrlInput } from '@ui/form/UrlInput';
import { Users03 } from '@ui/media/icons/Users03';
import { Share07 } from '@ui/media/icons/Share07';
import { Target05 } from '@ui/media/icons/Target05';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { FormSelect } from '@ui/form/Select/FormSelect';
import { Building07 } from '@ui/media/icons/Building07';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { HorizontalBarChart03 } from '@ui/media/icons/HorizontalBarChart03';
import { ArrowCircleBrokenUpLeft } from '@ui/media/icons/ArrowCircleBrokenUpLeft';
import { FormAutoresizeTextarea } from '@ui/form/Textarea/FormAutoresizeTextarea';
import { useOrganizationQuery } from '@organization/graphql/organization.generated';
import { Branches } from '@organization/components/Tabs/panels/AboutPanel/branches/Branches';
import { OwnerInput } from '@organization/components/Tabs/panels/AboutPanel/owner/OwnerInput';
import { ParentOrgInput } from '@organization/components/Tabs/panels/AboutPanel/branches/ParentOrgInput';
import { OrganizationNameInput } from '@organization/components/Tabs/panels/AboutPanel/OrganizationNameInput';

import { FormSocialInput } from '../../shared/FormSocialInput';
import { useAboutPanelMethods } from './hooks/useAboutPanelMethods';
import {
  OrganizationAboutForm,
  OrganizationAboutFormDto,
} from './OrganizationAbout.dto';
import {
  stageOptions,
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
  const { data, isLoading } = useOrganizationQuery(client, { id });

  const showParentRelationshipSelector = useFeatureIsOn(
    'show-parent-relationship-selector',
  );
  const parentRelationshipReadOnly = useFeatureIsOn(
    'parent-relationship-selector-read-only',
  );
  const orgNameReadOnly = useFeatureIsOn('org-name-readonly');
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
            debouncedMutateOrganization.cancel();
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
          } else {
            debouncedMutateOrganization.flush();
          }
        }
      }

      return next;
    },
  });

  useDeepCompareEffect(() => {
    setDefaultValues(defaultValues);
  }, [defaultValues]);

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
    <div className=' flex pt-4 px-6 w-full h-full overflow-y-auto flex-1 bg-gray-25 rounded-2xl'>
      <div className='flex h-full flex-col  overflow-visible w-full'>
        <div className='flex items-center justify-between'>
          <OrganizationNameInput
            orgNameReadOnly={orgNameReadOnly}
            isLoading={isLoading}
          />

          {data?.organization?.referenceId && (
            <div className='h-full ml-4'>
              <Tooltip label={'Copy ID'} asChild={false}>
                <Tag
                  colorScheme='gray'
                  variant='outline'
                  className='rounded-full cursor-pointer'
                  onClick={() => {
                    copyToClipboard(
                      data?.organization?.referenceId ?? '',
                      'Reference ID copied ',
                    );
                  }}
                >
                  <TagLabel>{data?.organization?.referenceId}</TagLabel>
                </Tag>
              </Tooltip>
            </div>
          )}
        </div>
        <FormUrlInput
          name='website'
          autoComplete='off'
          placeholder='www.'
          variant='unstyled'
          formId='organization-about'
        />

        <FormAutoresizeTextarea
          className='mb-6'
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

        <div className='flex flex-col w-full flex-1 items-start justify-start gap-0'>
          <div className='flex w-full'>
            <FormSelect
              isClearable
              name='isCustomer'
              formId='organization-about'
              placeholder='Relationship'
              options={relationshipOptions}
              leftElement={<HeartHand className='text-gray-500 mr-3' />}
            />
          </div>
          <div className='flex w-full'>
            <FormSelect
              isClearable
              name='stage'
              formId='organization-about'
              placeholder='Stage'
              options={stageOptions}
              leftElement={<HeartHand className='text-gray-500 mr-3' />}
            />
          </div>

          <FormSelect
            name='industry'
            isClearable
            placeholder='Industry'
            options={industryOptions}
            formId='organization-about'
            leftElement={<Building07 className='text-gray-500 mr-3' />}
          />

          <FormAutoresizeTextarea
            size='xs'
            className='items-start'
            name='targetAudience'
            formId='organization-about'
            placeholder='Target Audience'
            leftElement={<Target05 className='text-gray-500 mt-1 mr-1' />}
          />

          <FormSelect
            isClearable
            name='businessType'
            formId='organization-about'
            placeholder='Business Type'
            options={businessTypeOptions}
            leftElement={
              <ArrowCircleBrokenUpLeft className='text-gray-500 mr-3' />
            }
          />

          <div className='flex items-center justify-center w-full'>
            <div className='flex-1'>
              <FormSelect
                isClearable
                name='lastFundingRound'
                formId='organization-about'
                placeholder='Last funding round'
                options={lastFundingRoundOptions}
                leftElement={
                  <HorizontalBarChart03 className='text-gray-500 mr-3' />
                }
              />
            </div>

            {/* <FormNumberInputGroup
              name='lastFundingAmount'
              formId='organization-about'
              placeholder='Last funding amount'
              min={0}
              leftElement={<CurrencyDollar className='text-gray-500 size-4' />}
            /> */}
          </div>

          <FormSelect
            isClearable
            name='employees'
            options={employeesOptions}
            formId='organization-about'
            placeholder='Number of employees'
            leftElement={<Users03 className='text-gray-500 mr-3' />}
          />

          <OwnerInput
            id={id}
            owner={data?.organization?.owner}
            invalidateQuery={invalidateQuery}
          />

          <FormSocialInput
            name='socials'
            formId='organization-about'
            organizationId={id}
            addSocial={handleAddSocial}
            invalidateQuery={invalidateQuery}
            defaultValues={defaultValues.socials}
            placeholder='Social link'
            leftElement={<Share07 className='text-gray-500' />}
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
        </div>

        {data?.organization?.customerOsId && (
          <Tooltip label='Copy ID'>
            <span
              className='py-3 w-fit text-gray-400 cursor-pointer'
              onClick={() =>
                copyToClipboard(
                  data?.organization?.customerOsId ?? '',
                  'CustomerOS ID copied',
                )
              }
            >
              CustomerOS ID: {data?.organization?.customerOsId}
            </span>
          </Tooltip>
        )}
      </div>
    </div>
  );
};
