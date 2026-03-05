<template>
  <div class="grid grid-cols-12 gap-6">
    <div class="col-span-8">
      <CucumberDemoWidget />
    </div>
    <div class="col-span-4">
      <RunsPanel />
    </div>
  </div>
</template>

<script setup lang="ts">
import CucumberDemoWidget from '../components/CucumberDemoWidget.vue';
import RunsPanel from '../components/RunsPanel.vue';
</script>
