import { useRef } from 'react';
import { useParams } from 'react-router-dom';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';

import { Input } from '@ui/form/Input';
import { Select } from '@ui/form/Select';
import { UrlInput } from '@ui/form/UrlInput';
import { Users03 } from '@ui/media/icons/Users03';
import { Share07 } from '@ui/media/icons/Share07';
import { useStore } from '@shared/hooks/useStore';
import { Seeding } from '@ui/media/icons/Seeding';
import { Target05 } from '@ui/media/icons/Target05';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { AutoresizeTextarea } from '@ui/form/Textarea';
import { Building07 } from '@ui/media/icons/Building07';
import { Tag, TagLabel } from '@ui/presentation/Tag/Tag';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { Briefcase02 } from '@ui/media/icons/Briefcase02';
import { ActivityHeart } from '@ui/media/icons/ActivityHeart';
import { MessageXCircle } from '@ui/media/icons/MessageXCircle';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { Organization, OrganizationRelationship } from '@graphql/types';
import { HorizontalBarChart03 } from '@ui/media/icons/HorizontalBarChart03';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { Branches } from '@organization/components/Tabs/panels/AboutPanel/branches/Branches';
import { OwnerInput } from '@organization/components/Tabs/panels/AboutPanel/owner/OwnerInput';
import { ParentOrgInput } from '@organization/components/Tabs/panels/AboutPanel/branches/ParentOrgInput';

import { SocialIconInput } from '../../shared/SocialIconInput';
import {
  stageOptions,
  industryOptions,
  employeesOptions,
  businessTypeOptions,
  relationshipOptions,
  lastFundingRoundOptions,
} from './util';

const placeholders = {
  valueProposition: `What is this organization about - what they do, who they serve, and what makes them unique?`,
};

const iconMap = {
  Customer: <ActivityHeart className='text-gray-500' />,
  Prospect: <Seeding className='text-gray-500' />,
  'Not a fit': <MessageXCircle className='text-gray-500' />,
  'Former Customer': <BrokenHeart className='text-gray-500' />,
};

export const AboutPanel = observer(() => {
  const store = useStore();
  const id = useParams()?.id as string;
  const [_, copyToClipboard] = useCopyToClipboard();
  const nameRef = useRef<HTMLInputElement | null>(null);

  const showParentRelationshipSelector = useFeatureIsOn(
    'show-parent-relationship-selector',
  );
  const parentRelationshipReadOnly = useFeatureIsOn(
    'parent-relationship-selector-read-only',
  );
  const orgNameReadOnly = useFeatureIsOn('org-name-readonly');

  const organization = store.organizations.value.get(id);
  if (!organization) return null;

  const selectedRelationshipOption = relationshipOptions.find(
    (option) => option.value === organization?.value.relationship,
  );

  const selectedStageOption = stageOptions.find(
    (option) => option.value === organization?.value.stage,
  );

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = e.target;

    organization?.update((org) => {
      // @ts-expect-error fixme
      org[name] = value;

      return org;
    });
  };

  const menuHandleChange = (name: string, value: string) => {
    organization?.update((org) => {
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (org as any)[name] = value;

      return org;
    });
  };

  return (
    <div className=' flex pt-4 px-6 w-full h-full overflow-y-auto flex-1 bg-gray-25 rounded-2xl'>
      <div className='flex h-full flex-col  overflow-visible w-full'>
        <div className='flex items-center justify-between'>
          <Input
            className='font-semibold text-lg border-none overflow-hidden overflow-ellipsis'
            name='name'
            ref={nameRef}
            autoComplete='off'
            variant='unstyled'
            placeholder='Company name'
            disabled={orgNameReadOnly}
            value={organization?.value.name || ''}
            onChange={handleChange}
            size='xs'
          />
          {organization?.value.referenceId && (
            <div className='h-full ml-4'>
              <Tooltip label={'Copy ID'} asChild={false}>
                <Tag
                  colorScheme='gray'
                  variant='outline'
                  className='rounded-full cursor-pointer'
                  onClick={() => {
                    copyToClipboard(
                      organization?.value.referenceId ?? '',
                      'Reference ID copied ',
                    );
                  }}
                >
                  <TagLabel>{organization?.value.referenceId}</TagLabel>
                </Tag>
              </Tooltip>
            </div>
          )}
        </div>
        <UrlInput
          name='website'
          autoComplete='off'
          placeholder='www.'
          value={organization?.value?.website || ''}
          onChange={handleChange}
        />
        <AutoresizeTextarea
          className='mb-6'
          spellCheck={false}
          size='xs'
          name='valueProposition'
          placeholder={placeholders.valueProposition}
          value={organization?.value?.valueProposition || ''}
          onChange={handleChange}
        />
        {showParentRelationshipSelector && (
          <ParentOrgInput id={id} isReadOnly={parentRelationshipReadOnly} />
        )}
        <div className='flex items-center justify-center w-full'>
          <div className='flex-2'>
            <Menu>
              <MenuButton className='min-h-[40px] '>
                {
                  iconMap[
                    selectedRelationshipOption?.label as keyof typeof iconMap
                  ]
                }{' '}
                <span className='ml-2'>
                  {selectedRelationshipOption?.label}
                </span>
              </MenuButton>
              <MenuList side='bottom' align='center' className='min-w-[280px]'>
                {relationshipOptions
                  .filter(
                    (option) =>
                      !(
                        selectedRelationshipOption?.label === 'Customer' &&
                        option.label === 'Prospect'
                      ) &&
                      !(
                        selectedRelationshipOption?.label === 'Not a fit' &&
                        option.label === 'Prospect'
                      ),
                  )
                  .map((option) => (
                    <MenuItem
                      key={option.value}
                      disabled={
                        (selectedRelationshipOption?.label === 'Customer' ||
                          selectedRelationshipOption?.label === 'Not a fit') &&
                        option.label === 'Prospect'
                      }
                      onClick={() => {
                        menuHandleChange('relationship', option.value);
                      }}
                    >
                      {iconMap[option.label as keyof typeof iconMap]}
                      {option.label}
                    </MenuItem>
                  ))}
              </MenuList>
            </Menu>
          </div>

          {organization?.value?.relationship ===
            OrganizationRelationship.Prospect && (
            <div className='flex-1'>
              <Menu>
                <MenuButton className='min-h-[40px]'>
                  <Target05 className='text-gray-500 mb-0.5' />
                  <span className='ml-2'>{selectedStageOption?.label}</span>
                </MenuButton>
                <MenuList
                  side='bottom'
                  align='center'
                  className='min-w-[280px]'
                >
                  {stageOptions.map((option) => (
                    <MenuItem
                      key={option.value}
                      onClick={() => {
                        menuHandleChange('stage', option.value);
                      }}
                    >
                      {iconMap[option.label as keyof typeof iconMap]}
                      {option.label} || 'Stage'
                    </MenuItem>
                  ))}
                </MenuList>
              </Menu>
            </div>
          )}
        </div>
        <div className='flex flex-col w-full flex-1 items-start justify-start gap-0'>
          <Select
            name='industry'
            isClearable
            placeholder='Industry'
            options={industryOptions}
            value={organization?.value.industry}
            onChange={(value) => {
              organization?.update((org) => {
                org.industry = value as string;

                return org;
              });
            }}
            leftElement={<Building07 className='text-gray-500 mr-3' />}
          />

          <Select
            isClearable
            name='businessType'
            placeholder='Business Type'
            value={businessTypeOptions.map((option) =>
              option.value === organization?.value.market ? option : null,
            )}
            onChange={(value) => {
              organization?.update((org) => {
                if (value === null) org.market = null;
                else org.market = value.value;

                return org;
              });
            }}
            options={businessTypeOptions}
            leftElement={<Briefcase02 className='text-gray-500 mr-3' />}
          />

          <div className='flex items-center justify-center w-full'>
            <div className='flex-1'>
              <Select
                isClearable
                name='lastFundingRound'
                value={lastFundingRoundOptions.map((option) =>
                  option.value === organization?.value.lastFundingRound
                    ? option
                    : null,
                )}
                onChange={(value) => {
                  organization?.update((org) => {
                    if (value === null) org.lastFundingRound = null;
                    else org.lastFundingRound = value.value;

                    return org;
                  });
                }}
                placeholder='Last funding round'
                options={lastFundingRoundOptions}
                leftElement={
                  <HorizontalBarChart03 className='text-gray-500 mr-3' />
                }
              />
            </div>
          </div>

          <Select
            isClearable
            name='employees'
            value={organization?.value.employees}
            onChange={(value) => {
              organization?.update((org) => {
                org.employees = value as string;

                return org;
              });
            }}
            options={employeesOptions}
            placeholder='Number of employees'
            leftElement={<Users03 className='text-gray-500 mr-3' />}
          />

          <OwnerInput id={id} owner={organization?.value.owner} />

          <SocialIconInput
            name='socials'
            organizationId={id}
            placeholder='Social link'
            leftElement={<Share07 className='text-gray-500' />}
          />

          {showParentRelationshipSelector &&
            organization?.value?.subsidiaries?.length > 0 && (
              <Branches
                id={id}
                isReadOnly={parentRelationshipReadOnly}
                branches={
                  (organization?.value
                    ?.subsidiaries as Organization['subsidiaries']) ?? []
                }
              />
            )}
        </div>
        {organization?.value.customerOsId && (
          <Tooltip label='Copy ID'>
            <span
              className='py-3 w-fit text-gray-400 cursor-pointer'
              onClick={() =>
                copyToClipboard(
                  organization?.value.customerOsId ?? '',
                  'CustomerOS ID copied',
                )
              }
            >
              CustomerOS ID: {organization?.value.customerOsId}
            </span>
          </Tooltip>
        )}
      </div>
    </div>
  );
});
