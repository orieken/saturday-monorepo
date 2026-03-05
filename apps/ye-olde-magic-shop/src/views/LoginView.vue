<template>
  <div class="min-h-[calc(100vh-200px)] flex items-center justify-center py-12 px-4">
    <div class="w-full max-w-md">
      <!-- Magical Header -->
      <div class="text-center mb-8">
        <h2 class="text-4xl font-magical text-dnd-gold mb-2">
          {{ isLogin ? 'Enter the Arcane Registry' : 'Join the Guild' }}
        </h2>
        <p class="text-gray-400">
          {{ isLogin ? 'Access your magical account' : 'Create your wizard account' }}
        </p>
      </div>

      <!-- Auth Form Card -->
      <div class="bg-dnd-dark border-2 border-dnd-slate rounded-lg p-8 shadow-2xl relative overflow-hidden">
        <!-- Magical glow effect -->
        <div class="absolute inset-0 bg-gradient-to-br from-purple-900/10 via-transparent to-blue-900/10 pointer-events-none"></div>
        
        <div class="relative z-10">
          <!-- Error Message -->
          <div v-if="authStore.error" class="mb-6 p-4 bg-red-900/30 border border-red-700 rounded text-red-300 text-sm" data-testid="login-error">
            <span class="font-magical">⚠️ {{ authStore.error }}</span>
          </div>

          <!-- Login Form -->
          <form v-if="isLogin" @submit.prevent="handleLogin" class="space-y-6">
            <div>
              <label for="login-email" class="block text-sm font-medium text-dnd-parchment mb-2">
                Email Address
              </label>
              <input
                id="login-email"
                v-model="loginForm.email"
                type="email"
                required
                class="w-full px-4 py-3 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold focus:ring-2 focus:ring-dnd-gold/50 outline-none transition-all text-white placeholder-gray-500"
                placeholder="wizard@example.com"
                data-testid="email-input"
              />
            </div>

            <div>
              <label for="login-password" class="block text-sm font-medium text-dnd-parchment mb-2">
                Password
              </label>
              <input
                id="login-password"
                v-model="loginForm.password"
                type="password"
                required
                class="w-full px-4 py-3 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold focus:ring-2 focus:ring-dnd-gold/50 outline-none transition-all text-white placeholder-gray-500"
                placeholder="Enter your secret incantation"
                data-testid="password-input"
              />
            </div>

            <button
              type="submit"
              :disabled="authStore.loading"
              class="w-full py-3 bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white font-magical text-lg rounded shadow-lg transition-all transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
              data-testid="login-button"
            >
              <span v-if="!authStore.loading">🔮 Cast Login Spell</span>
              <span v-else>✨ Channeling magic...</span>
            </button>
          </form>

          <!-- Register Form -->
          <form v-else @submit.prevent="handleRegister" class="space-y-6">
            <div>
              <label for="register-name" class="block text-sm font-medium text-dnd-parchment mb-2">
                Wizard Name
              </label>
              <input
                id="register-name"
                v-model="registerForm.name"
                type="text"
                required
                class="w-full px-4 py-3 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold focus:ring-2 focus:ring-dnd-gold/50 outline-none transition-all text-white placeholder-gray-500"
                placeholder="Gandalf the Grey"
                data-testid="register-name"
              />
            </div>

            <div>
              <label for="register-email" class="block text-sm font-medium text-dnd-parchment mb-2">
                Email Address
              </label>
              <input
                id="register-email"
                v-model="registerForm.email"
                type="email"
                required
                class="w-full px-4 py-3 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold focus:ring-2 focus:ring-dnd-gold/50 outline-none transition-all text-white placeholder-gray-500"
                placeholder="wizard@example.com"
                data-testid="register-email"
              />
            </div>

            <div>
              <label for="register-password" class="block text-sm font-medium text-dnd-parchment mb-2">
                Password
              </label>
              <input
                id="register-password"
                v-model="registerForm.password"
                type="password"
                required
                minlength="6"
                class="w-full px-4 py-3 bg-black/40 border border-dnd-slate rounded focus:border-dnd-gold focus:ring-2 focus:ring-dnd-gold/50 outline-none transition-all text-white placeholder-gray-500"
                placeholder="Create a secret incantation"
                data-testid="register-password"
              />
            </div>

            <button
              type="submit"
              :disabled="authStore.loading"
              class="w-full py-3 bg-gradient-to-r from-green-600 to-teal-600 hover:from-green-700 hover:to-teal-700 text-white font-magical text-lg rounded shadow-lg transition-all transform hover:scale-105 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none"
              data-testid="register-submit"
            >
              <span v-if="!authStore.loading">✨ Join the Guild</span>
              <span v-else>🌟 Inscribing your name...</span>
            </button>
          </form>

          <!-- Toggle between Login/Register -->
          <div class="mt-6 text-center">
            <button
              @click="toggleMode"
              class="text-dnd-gold hover:text-yellow-400 transition-colors text-sm"
              data-testid="toggle-auth-mode"
            >
              {{ isLogin ? "Don't have an account? Join the Guild" : 'Already a member? Enter the Registry' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Demo Credentials -->
      <div class="mt-6 p-4 bg-blue-900/20 border border-blue-700/50 rounded text-sm text-blue-300">
        <p class="font-magical mb-2">🧙 Demo Credentials:</p>
        <p><strong>Email:</strong> gandalf@middleearth.com</p>
        <p><strong>Password:</strong> password123</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/authStore';

const router = useRouter();
const authStore = useAuthStore();

const isLogin = ref(true);

const loginForm = ref({
  email: '',
  password: ''
});

const registerForm = ref({
  name: '',
  email: '',
  password: ''
});

function toggleMode() {
  isLogin.value = !isLogin.value;
  authStore.error = null;
}

async function handleLogin() {
  const success = await authStore.login(loginForm.value);
  if (success) {
    router.push('/account');
  }
}

async function handleRegister() {
  const success = await authStore.register(registerForm.value);
  if (success) {
    router.push('/account');
  }
}
</script>
