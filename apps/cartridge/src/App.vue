<template>
  <main class="min-h-screen bg-slate-950 text-slate-100">
    <div class="max-w-6xl mx-auto px-4 py-6">
      <header class="mb-6 flex items-center justify-between gap-4 border-b border-slate-800 pb-4">
        <div class="flex items-center gap-4">
          <img src="/logo.png" alt="Cartridge Logo" class="w-10 h-10 object-contain" />
          <div>
            <h1 class="text-2xl font-bold tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-teal-400 font-mono">CARTRIDGE</h1>
            <p class="text-xs text-slate-400 flex items-center gap-1">
              <span>Saturday Framework Plugin</span>
            </p>
          </div>
        </div>
        <div class="flex items-center gap-6">
          <nav class="flex items-center gap-4 text-sm font-medium">
            <router-link to="/" class="text-slate-400 hover:text-white" active-class="text-emerald-400">Dashboard</router-link>
            <router-link to="/reports" class="text-slate-400 hover:text-white" active-class="text-emerald-400">Reports</router-link>
          </nav>
          <span class="text-xs text-slate-500 border-l border-slate-800 pl-6">
            Demo suite: <span class="font-mono text-emerald-400">demo-ecommerce</span>
          </span>
        </div>
      </header>

      <router-view />
    </div>
    <ToastContainer />
  </main>
</template>

<script setup lang="ts">
import ToastContainer from './components/ToastContainer.vue';
</script>
