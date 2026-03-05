<template>
  <span
    :class="[
      'inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium',
      colorClass
    ]"
  >
    {{ statusLabel }}
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{
  status: 'idle' | 'running' | 'passed' | 'failed';
}>();

const statusLabel = {
  idle: 'Idle',
  running: 'Running',
  passed: 'Passed',
  failed: 'Failed'
}[props.status];

const colorClass = {
  idle: 'bg-slate-700 text-slate-200',
  running: 'bg-amber-500 text-slate-900',
  passed: 'bg-emerald-500 text-slate-900',
  failed: 'bg-rose-500 text-slate-900'
}[props.status];
</script>
