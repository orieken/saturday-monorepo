<template>
  <div class="app min-h-screen flex flex-col">
    <header class="bg-dnd-dark border-b border-dnd-slate shadow-lg sticky top-0 z-50">
      <div class="container mx-auto px-4 py-4 flex justify-between items-center">
        <router-link to="/" class="flex items-center gap-4 hover:opacity-90 transition-opacity">
          <img src="/logo.png" alt="Ye Olde Magic Shop Logo" class="w-12 h-12 object-contain drop-shadow-md" />
          <h1 class="text-2xl font-bold tracking-wider text-dnd-gold drop-shadow-md">Ye Olde Magic Shop</h1>
        </router-link>
        <nav class="flex gap-6 text-dnd-parchment font-magical items-center">
          <router-link to="/" class="hover:text-dnd-gold transition-colors text-lg" data-testid="shop-link">Shop</router-link>
          <router-link to="/cart" class="hover:text-dnd-gold transition-colors text-lg flex items-center gap-2" data-testid="cart-link">
            <span>Cart</span>
            <span v-if="totalItems > 0" class="bg-dnd-crimson text-white text-xs rounded-full px-2 py-0.5 font-sans" data-testid="cart-count">{{ totalItems }}</span>
          </router-link>
          <router-link 
            v-if="authStore.isAuthenticated" 
            to="/account" 
            class="hover:text-dnd-gold transition-colors text-lg" 
            data-testid="account-link"
          >
            <span data-testid="user-name">👤 {{ authStore.user?.name }}</span>
          </router-link>
          <router-link 
            v-else 
            to="/login" 
            class="hover:text-dnd-gold transition-colors text-lg" 
            data-testid="login-link"
          >
            🔮 Login
          </router-link>
        </nav>
      </div>
    </header>

    <main class="container mx-auto px-4 py-8 flex-1">
      <router-view />
    </main>

    <footer class="bg-dnd-dark border-t border-dnd-slate py-6 mt-auto">
      <div class="container mx-auto px-4 text-center text-gray-500 text-sm">
        &copy; 2025 Dungeon Master Supplies Co.
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { useCartStore } from './stores/cartStore';
import { useAuthStore } from './stores/authStore';

const cart = useCartStore();
const { totalItems } = storeToRefs(cart);

const authStore = useAuthStore();

cart.hydrate();
authStore.hydrate();
</script>

