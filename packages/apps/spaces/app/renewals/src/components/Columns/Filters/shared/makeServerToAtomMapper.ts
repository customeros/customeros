import { match } from 'ts-pattern';
import { Pattern } from 'ts-pattern/dist/types/Pattern';
import { PickReturnValue } from 'ts-pattern/dist/types/Match';

import { Filter } from '@graphql/types';

export const makeServerToAtomMapper =
  <AtomState>(
    matchPattern: Pattern<Filter>,
    mapper: (
      selections: Filter,
      value: Filter,
    ) => PickReturnValue<AtomState, AtomState>,
    defaultState: PickReturnValue<AtomState, AtomState>,
  ) =>
  (input: Filter) =>
    match(input)
      .returnType<AtomState>()
      .with(matchPattern, mapper)
      .otherwise(() => defaultState);
