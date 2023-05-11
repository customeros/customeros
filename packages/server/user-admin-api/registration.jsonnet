function(ctx) {
properties: {
   email: ctx.identity.traits.email,
   firstname: ctx.identity.traits.name.first, 
   lastname: ctx.identity.traits.name.last, 
   [if 'workspace' in ctx.identity.traits then 'workspace' else null]: ctx.identity.traits.workspace,
   provider: ctx.identity.traits.provider
},
}
