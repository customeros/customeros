import { Tag, TagLabel } from '@ui/presentation/Tag';

export const PersonasFlowProfile = () => {
  return (
    <div className='flex'>
      <p className='font-semibold'> Tag contact with </p>
      <Tag variant='subtle'>
        <TagLabel>Solo RevOps</TagLabel>
      </Tag>
    </div>
  );
};
