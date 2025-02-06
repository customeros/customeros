import { useMemo } from 'react';

import { observer } from 'mobx-react-lite';
import { useFeatureIsOn } from '@growthbook/growthbook-react';
import { EditOrganizationTagUsecase } from '@domain/usecases/organization-details/edit-organization-tag.usecase';
import { SaveOrganizationRelationshipAndStageUsecase } from '@domain/usecases/organization-details/save-organization-relationship-and-stage.usecase';

import { cn } from '@ui/utils/cn';
import { Icon } from '@ui/media/Icon';
import { flags } from '@ui/media/flags';
import { Spinner } from '@ui/feedback/Spinner';
import { IconButton } from '@ui/form/IconButton';
import { useStore } from '@shared/hooks/useStore';
import { Tag, TagLabel } from '@ui/presentation/Tag';
import { Tooltip } from '@ui/overlay/Tooltip/Tooltip';
import { SocialMediaList } from '@organization/components/Tabs';
import { useCopyToClipboard } from '@shared/hooks/useCopyToClipboard';
import { TruncatedText } from '@ui/presentation/TruncatedText/TruncatedText';
import { Menu, MenuItem, MenuList, MenuButton } from '@ui/overlay/Menu/Menu';
import { IcpBadge } from '@shared/components/OrganizationDetails/components/icp';
import { Domains } from '@shared/components/OrganizationDetails/components/domains';
import { OwnerInput } from '@shared/components/OrganizationDetails/components/owner';
import { Branches } from '@shared/components/OrganizationDetails/components/branches';
import { AboutTabField } from '@shared/components/OrganizationDetails/components/AboutTabField';
import {
  FlagWrongFields,
  OrganizationRelationship,
} from '@shared/types/__generated__/graphql.types';
import {
  stageOptions,
  getStageOptions,
  relationshipOptions,
} from '@organization/components/Tabs/panels/AboutPanel/util';

import { Tags } from '../Tags';

const iconMap = {
  Customer: <Icon name='activity-heart' className='text-gray-500' />,
  Prospect: <Icon name='seeding' className='text-gray-500' />,
  'Not a fit': <Icon name='message-x-circle' className='text-gray-500' />,
  'Former Customer': <Icon name='broken-heart' className='text-gray-500' />,
  unknown: <Icon className='text-gray-500' name='align-horizontal-centre-02' />,
};

interface OrganizationDetailsProps {
  id: string;
}

export const OrganizationDetails = observer(
  ({ id }: OrganizationDetailsProps) => {
    const store = useStore();
    const organization = store.organizations.getById(id);

    const [_, copyToClipboard] = useCopyToClipboard();

    const showParentRelationshipSelector = useFeatureIsOn(
      'show-parent-relationship-selector',
    );
    const parentRelationshipReadOnly = useFeatureIsOn(
      'parent-relationship-selector-read-only',
    );

    const tagsUsecase = useMemo(() => new EditOrganizationTagUsecase(id), [id]);
    const saveRelationshipAndStageUsecase = useMemo(
      () => new SaveOrganizationRelationshipAndStageUsecase(id),
      [id],
    );

    const handleCreateOption = () => {
      tagsUsecase.create();
    };

    const isEnriching = organization?.isEnriching;

    if (!organization) return null;

    const applicableStageOptions = getStageOptions(
      organization?.value?.relationship,
    );

    const selectedRelationshipOption = relationshipOptions.find(
      (option) => option.value === organization?.value?.relationship,
    );

    const selectedStageOption = stageOptions.find(
      (option) => option.value === organization?.value?.stage,
    );

    return (
      <div className='flex pt-[6px] px-6 w-full h-full  flex-1 bg-gray-25 rounded-2xl'>
        <div className='flex h-full flex-col  overflow-visible w-full'>
          {isEnriching && (
            <div className='flex items-center justify-start gap-2 border-[1px] text-sm border-grayModern-100 bg-grayModern-50 rounded-[4px] py-1 px-2 '>
              <Spinner
                label='enriching company'
                className='text-grayModern-300 fill-grayModern-500 size-4'
              />
              <span className='font-medium'>
                We're enriching this company’s details.
              </span>
            </div>
          )}

          <div className='flex items-center justify-between'>
            <p className='font-semibold text-base mt-0.5 overflow-hidden overflow-ellipsis'>
              {organization?.value?.name ?? ''}
            </p>

            <div className='flex justify-between items-start h-full gap-x-2 mt-[11px]'>
              {organization.value?.referenceId && (
                <div className='ml-4'>
                  <Tooltip asChild={false} label={'Copy ID'}>
                    <Tag
                      variant='outline'
                      colorScheme='gray'
                      className='cursor-pointer w-full max-w-[100px]'
                      onClick={() => {
                        copyToClipboard(
                          organization.value?.referenceId ?? '',
                          'Reference ID copied ',
                        );
                      }}
                    >
                      <TagLabel className='truncate overflow-hidden whitespace-nowrap'>
                        {organization.value?.referenceId}
                      </TagLabel>
                    </Tag>
                  </Tooltip>
                </div>
              )}

              <IcpBadge id={store.ui.focusRow ?? id} />
              {store.ui.showPreviewCard && (
                <IconButton
                  size='xxs'
                  variant='ghost'
                  aria-label='close preview company'
                  icon={<Icon name='x-close' className='size-4' />}
                  onClick={() => store.ui.setShowPreviewCard(false)}
                />
              )}
            </div>
          </div>

          <Domains id={id} />

          <div className='flex flex-col w-full flex-1 items-start justify-start gap-3 mt-2'>
            {!!organization?.value?.description && (
              <TruncatedText
                maxLines={7}
                className='text-sm'
                data-test='org-about-description'
                text={organization.value.description}
              />
            )}
            <SocialMediaList dataTest='org-about-social-link' />
            <Tags
              dataTest='org-about-tags'
              placeholder='Company tags'
              inputPlaceholder='Search...'
              onCreate={handleCreateOption}
              options={tagsUsecase.tagList}
              value={tagsUsecase.selectedTags}
              inputValue={tagsUsecase.searchTerm}
              setInputValue={tagsUsecase.setSearchTerm}
              leftAccessory={
                <Icon name='tag-01' className='mr-3 text-gray-500' />
              }
              onChange={(selection) => {
                tagsUsecase.select(selection.map((o) => o.value));
              }}
            />

            <div className='flex items-center justify-center w-full '>
              <div
                data-test='org-about-relationship'
                className='flex-2 flex items-center'
              >
                <Menu>
                  <Tooltip align='start' label='Relationship'>
                    <MenuButton
                      data-test='org-about-relationship'
                      className='min-h-[20px] text-md outline-none focus:outline-none items-center'
                    >
                      {
                        iconMap[
                          (selectedRelationshipOption?.label ??
                            'unknown') as keyof typeof iconMap
                        ]
                      }
                      {''}
                      <span
                        className={cn(
                          'ml-3 text-sm',
                          !selectedRelationshipOption?.label && 'text-gray-400',
                        )}
                      >
                        {selectedRelationshipOption?.label ?? 'Relationship'}
                      </span>
                    </MenuButton>
                  </Tooltip>
                  <MenuList side='bottom' align='start'>
                    {relationshipOptions.map((option) => (
                      <MenuItem
                        key={option.value}
                        onClick={() => {
                          saveRelationshipAndStageUsecase.execute({
                            relationship: option.value,
                          });
                        }}
                      >
                        {iconMap[option.label as keyof typeof iconMap]}
                        {option.label}
                      </MenuItem>
                    ))}
                  </MenuList>
                </Menu>
              </div>
              {selectedRelationshipOption?.value !==
                OrganizationRelationship.Customer && (
                <div
                  data-test='org-about-stage'
                  className='flex-1 flex items-center'
                >
                  <Menu>
                    <Tooltip label='Stage' align='start'>
                      <MenuButton className='min-h-[20px] outline-none focus:outline-none'>
                        <Icon
                          name='target-05'
                          className='text-gray-500 mb-0.5'
                        />
                        <span className='ml-3 text-sm'>
                          {selectedStageOption?.label || 'Stage'}
                        </span>
                      </MenuButton>
                    </Tooltip>
                    <MenuList side='bottom' align='start'>
                      {applicableStageOptions.map((option) => (
                        <MenuItem
                          key={option.value}
                          onClick={() => {
                            saveRelationshipAndStageUsecase.execute({
                              stage: option.value,
                            });
                          }}
                        >
                          {iconMap[option.label as keyof typeof iconMap]}
                          {option.label}
                        </MenuItem>
                      ))}
                    </MenuList>
                  </Menu>
                </div>
              )}
            </div>
            <AboutTabField
              id={store.ui.focusRow ?? id}
              dataTest={'org-about-industry'}
              placeholder='Industry not found yet'
              value={organization?.value?.industryName}
              field={FlagWrongFields.OrganizationIndustry}
              flaggedAsIncorrect={organization?.value?.wrongIndustry ?? false}
              icon={<Icon name='building-07' className='text-gray-500 mr-3' />}
            />

            {organization.country && (
              <Tooltip align='start' label='Country'>
                <p className='text-sm flex items-center cursor-default'>
                  <span className='flex items-center mr-3'>
                    {organization.value.locations?.[0]?.countryCodeA2 &&
                      flags[organization.value.locations?.[0]?.countryCodeA2]}
                  </span>

                  {organization.country}
                </p>
              </Tooltip>
            )}

            {typeof organization.value!.employees === 'number' && (
              <Tooltip align='start' label='Number of employees'>
                <p className='text-sm flex items-center cursor-default '>
                  <Icon name='users-02' className='text-gray-500 mr-3' />
                  {organization.value!.employees} employees
                </p>
              </Tooltip>
            )}
            <OwnerInput
              id={id ?? store.ui.focusRow}
              dataTest='org-about-org-owner'
            />

            {showParentRelationshipSelector &&
              organization?.value?.subsidiaries?.length > 0 && (
                <Branches id={id} isReadOnly={parentRelationshipReadOnly} />
              )}
          </div>
          {organization?.value.customerOsId && (
            <Tooltip label='Copy ID'>
              <span
                className='py-3 w-fit text-gray-400 cursor-pointer text-sm'
                onClick={() =>
                  copyToClipboard(
                    organization.value?.customerOsId ?? '',
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
  },
);
