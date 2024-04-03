export const ContractCell = ({ value }: { value: string }) => {
  return (
    <span className='font-medium line-clamp-1'>{`${value}'s contract`}</span>
  );
};
