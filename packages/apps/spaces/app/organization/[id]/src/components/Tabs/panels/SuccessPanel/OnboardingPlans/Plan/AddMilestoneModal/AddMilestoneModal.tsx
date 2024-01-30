import addHours from 'date-fns/addHours';

import { Text } from '@ui/typography/Text';
import { Heading } from '@ui/typography/Heading';
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

interface AddMilestoneModalProps {
  isOpen: boolean;
  onClose: () => void;
  masterPlanId: string;
  existingMilestoneNames: string[];
  onAddMilestone: (name: string, dueDate: string) => void;
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

  const handleAdd = (name: string, durationHours: number) => () => {
    const dueDate = addHours(new Date(), durationHours).toISOString();
    onAddMilestone(name, dueDate);
  };

  return (
    <Modal closeOnEsc isOpen={isOpen} onClose={onClose}>
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
        <ModalBody>
          {filteredMilestones.length ? (
            filteredMilestones.map((milestone) => (
              <Card
                key={milestone.id}
                cursor='pointer'
                variant='outlinedElevated'
                onClick={handleAdd(milestone.name, milestone.durationHours)}
              >
                <CardBody pl='3' pb='4'>
                  <Text fontWeight='medium' w='full'>
                    {milestone.name}
                  </Text>
                  <Text color='gray.500'>
                    Max duration:{' '}
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
