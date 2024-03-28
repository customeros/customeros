import ReactSelect from 'react-select';

export const Select = () => {
  return (
    <>
      <ReactSelect
        unstyled
        classNames={{
          menu: (state) => 'bg-gray-100',
        }}
      />
    </>
  );
};
