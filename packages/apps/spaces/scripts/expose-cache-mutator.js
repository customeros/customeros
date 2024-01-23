const { parse } = require('graphql');
const camelCase = require('lodash/camelCase');
const upperFirst = require('lodash/upperFirst');
const { oldVisit } = require('@graphql-codegen/plugin-helpers');

const plugin = (_schema, documents, _config) => {
  let output;

  documents?.forEach((doc) => {
    const ast = parse(doc?.rawSDL);

    oldVisit(ast, {
      enter: {
        OperationDefinition(node) {
          if (node.operation === 'query') {
            const queryName = camelCase(node.name.value);

            output = '\n'.concat(
              template(queryName),
              '\n',
              template(queryName, { infinite: true }),
            );
          }
        },
      },
    });
  });

  if (output)
    return {
      content: output,
    };
};

function template(name, options = { infinite: false }) {
  const query = upperFirst(name); // GetOrganizations
  const queryName = `${query}Query`; // GetOrganizationsQuery
  const hookName = options.infinite
    ? `useInfinite${queryName}` // useInfiniteGetOrganizationsQuery
    : `use${queryName}`; // useGetOrganizationsQuery
  const variablesType = `${queryName}Variables`; // GetOrganizationsQueryVariables
  const queryType = options.infinite
    ? `InfiniteData<${queryName}>` // InfiniteData<GetOrganizationsQuery>
    : queryName; // GetOrganizationsQuery

  return `${hookName}.mutateCacheEntry =
  (queryClient: QueryClient, variables?: ${variablesType}) =>
  (mutator: (cacheEntry: ${queryType}) => ${queryType}) => {
    const cacheKey = ${hookName}.getKey(variables);
    const previousEntries =
      queryClient.getQueryData<${queryType}>(cacheKey);
    if (previousEntries) {
      queryClient.setQueryData<${queryType}>(cacheKey, mutator);
    }
    return { previousEntries };
  }`;
}

module.exports = {
  plugin,
};
