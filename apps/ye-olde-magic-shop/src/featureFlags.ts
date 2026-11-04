// The build-time gate for /reagent-map — read by both the router and App.vue.
export function isReagentMapEnabled(value: string | undefined): boolean {
  return value === 'true';
}
