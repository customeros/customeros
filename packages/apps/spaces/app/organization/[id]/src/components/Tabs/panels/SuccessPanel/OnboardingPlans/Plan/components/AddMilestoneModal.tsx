import addHours from 'date-fns/addHours';

import { VStack } from '@ui/layout/Stack';
import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
import { MasterPlanMilestone } from '@graphql/types';
import { Card, CardBody } from '@ui/presentation/Card';
import { PlusSquare } from '@ui/media/icons/PlusSquare';
import { FeaturedIcon } from '@ui/media/Icon/FeaturedIcon';
import { getGraphQLClient } from '@shared/util/getGraphQLClient';
import { useMasterPlansQuery } from '@shared/graphql/masterPlans.generated';
import {
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalContent,
  ModalOverlay,
  ModalCloseButton,
} from '@ui/overlay/Modal';

import { NewMilestoneInput } from '../../types';

type MasterPlanMilestoneInput = Pick<
  MasterPlanMilestone,
  'name' | 'durationHours' | 'order' | 'items'
>;

interface AddMilestoneModalProps {
  isOpen: boolean;
  onClose: () => void;
  masterPlanId: string;
  existingMilestoneNames: string[];
  onAddMilestone: (input: NewMilestoneInput) => void;
}

export const AddMilestoneModal = ({
  isOpen,
  onClose,
  masterPlanId,
  onAddMilestone,
  existingMilestoneNames,
}: AddMilestoneModalProps) => {
  const client = getGraphQLClient();
  const { data } = useMasterPlansQuery(client);

  const masterPlans = data?.masterPlans ?? [];
  const foundMasterPlan = masterPlans.find((mp) => mp.id === masterPlanId);
  const milestones = foundMasterPlan?.milestones ?? [];
  const filteredMilestones = milestones.filter(
    (m) => !existingMilestoneNames.includes(m.name),
  );

  const handleAdd =
    ({ order, name, durationHours, items }: MasterPlanMilestoneInput) =>
    () => {
      const dueDate = addHours(new Date(), durationHours).toISOString();

      onAddMilestone({
        name,
        order,
        items,
        dueDate,
      });
    };

  return (
    <Modal closeOnEsc isOpen={isOpen} onClose={onClose} scrollBehavior='inside'>
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
        <ModalHeader>
          <FeaturedIcon size='lg' colorScheme='primary'>
            <PlusSquare />
          </FeaturedIcon>
          <Heading fontSize='lg' mt='4'>
            Add a milestone
          </Heading>
        </ModalHeader>

        <ModalCloseButton />
        <ModalBody as={VStack} spacing='2'>
          {filteredMilestones.length ? (
            filteredMilestones.map((milestone) => (
              <Card
                w='full'
                key={milestone.id}
                cursor='pointer'
                variant='outlinedElevated'
                onClick={handleAdd({
                  name: milestone.name,
                  order: milestone.order,
                  items: milestone.items,
                  durationHours: milestone.durationHours,
                })}
              >
                <CardBody pl='3' pb='4'>
                  <Text fontWeight='medium' w='full'>
                    {milestone.name}
                  </Text>
                  <Text color='gray.500'>
                    Time budget:{' '}
                    {`${milestone.durationHours / 24} ${
                      milestone.durationHours > 24 ? 'days' : 'day'
                    }`}
                  </Text>
                </CardBody>
              </Card>
            ))
          ) : (
            <Text>No milestones to add</Text>
          )}
        </ModalBody>

        <ModalFooter />
      </ModalContent>
    </Modal>
  );
};
